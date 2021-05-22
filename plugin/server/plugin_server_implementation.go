package server

import (
	"context"
	"encoding/json"
	"github.com/blinkops/plugin-sdk/plugin"
	"github.com/blinkops/plugin-sdk/plugin/connections"
	pb "github.com/blinkops/plugin-sdk/plugin/proto"
	log "github.com/sirupsen/logrus"
)

type PluginGRPCService struct {
	pb.UnimplementedPluginServer

	plugin plugin.Implementation
}

func translateToProtoConnections(connections map[string]connections.Connection) map[string]*pb.Connection {

	protoConnections := map[string]*pb.Connection{}
	for connectionName, connection := range connections {

		protoConnectionFields := map[string]*pb.ConnectionField{}
		for fieldName, field := range connection.Fields {
			protoConnectionFields[fieldName] = &pb.ConnectionField{
				Name:        field.Name,
				Type:        field.FieldType,
				Required:    field.Required,
				PlaceHolder: field.Placeholder,
				InputType:   field.InputType,
				Patterns:    field.Pattern,
				Values:      field.Values,
			}
		}

		protoConnections[connectionName] = &pb.Connection{
			Name:   connectionName,
			Fields: protoConnectionFields,
		}
	}

	return protoConnections
}

func (service *PluginGRPCService) Describe(ctx context.Context, empty *pb.Empty) (*pb.PluginDescription, error) {
	pluginDescription := service.plugin.Describe()
	actions, err := service.GetActions(ctx, empty)
	if err != nil {
		return nil, err
	}

	return &pb.PluginDescription{
		Name:        pluginDescription.Name,
		Description: pluginDescription.Description,
		Tags:        pluginDescription.Tags, Provider: pluginDescription.Provider,
		Actions:     actions.Actions,
		Connections: translateToProtoConnections(pluginDescription.Connections),
	}, nil
}

func (service *PluginGRPCService) GetActions(ctx context.Context, empty *pb.Empty) (*pb.ActionList, error) {

	actions := service.plugin.GetActions()

	var protoActions []*pb.Action
	for _, action := range actions {

		protoAction := &pb.Action{
			Name:        action.Name,
			Description: action.Description,
			Active:      action.Enabled,
		}

		var protoParameters []*pb.ActionParameter
		for name, parameter := range action.Parameters {
			protoParameter := &pb.ActionParameter{
				Name:        name,
				Type:        parameter.Type,
				Description: parameter.Description,
				Required:    parameter.Required,
			}

			protoParameters = append(protoParameters, protoParameter)
		}
		protoAction.Parameters = protoParameters
		if action.Output != nil {
			protoAction.Output = &pb.Output{
				Table: action.Output.Name,
			}
			for _, field := range action.Output.Fields {
				protoField := &pb.Field{
					Name: field.Name,
					Type: field.Type,
				}
				protoAction.Output.Fields = append(protoAction.Output.Fields, protoField)
			}
		}

		protoActions = append(protoActions, protoAction)
	}

	return &pb.ActionList{Actions: protoActions}, nil
}

func translateActionContext(request *pb.ExecuteActionRequest) (map[string]interface{}, error) {

	rawContext := map[string]interface{}{}
	if len(request.Context) > 0 {
		err := json.Unmarshal(request.Context, &rawContext)
		if err != nil {
			log.Error("Failed to unmarshal action context with error: ", err)
			return nil, err
		}
	}

	return rawContext, nil
}

func translateConnectionInstances(request *pb.ExecuteActionRequest) (map[string]connections.ConnectionInstance, error) {

	concreteConnections := map[string]connections.ConnectionInstance{}
	for protoName, protoConnection := range request.Connections {
		concreteConnections[protoName] = connections.ConnectionInstance{
			VaultUrl: protoConnection.VaultUrl,
			Name:     protoConnection.Name,
			Id:       protoConnection.Id,
			Token:    protoConnection.Token,
		}
	}

	return concreteConnections, nil
}

func (service *PluginGRPCService) ExecuteAction(_ context.Context, request *pb.ExecuteActionRequest) (*pb.ExecuteActionResponse, error) {

	actionRequest := plugin.ExecuteActionRequest{
		Name:       request.Name,
		Parameters: request.Parameters,
	}

	rawContext, err := translateActionContext(request)
	if err != nil {
		return nil, err
	}

	connectionInstances, err := translateConnectionInstances(request)
	if err != nil {
		return nil, err
	}

	actionContext := plugin.NewActionContext(rawContext, connectionInstances)
	response, err := service.plugin.ExecuteAction(actionContext, &actionRequest)
	if err != nil {
		return nil, err
	}

	updatedActionContext, err := actionContext.GetMarshaledContext()
	if err != nil {
		log.Error("Failed to marshal context after action execution, error: ", err)
	}

	res := &pb.ExecuteActionResponse{
		ErrorCode: response.ErrorCode,
		Result:    response.Result,
		Context:   updatedActionContext,
		LogBuffer: actionContext.GetRawLogBuffer(),
	}

	for _, row := range response.Rows {
		pbRow := &pb.Row{Data: row}
		res.Rows = append(res.Rows, pbRow)
	}

	return res, nil
}

func NewPluginServiceImplementation(plugin plugin.Implementation) *PluginGRPCService {
	return &PluginGRPCService{plugin: plugin}
}
