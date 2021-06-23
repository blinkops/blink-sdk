package server

import (
	"context"
	"encoding/json"
	"github.com/blinkops/blink-sdk/plugin"
	"github.com/blinkops/blink-sdk/plugin/connections"
	pb "github.com/blinkops/blink-sdk/plugin/proto"
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
				Field: &pb.FormField{
					Name:        field.Name,
					Type:        field.FieldType,
					InputType:   field.InputType,
					Required:    field.Required,
					Description: field.Description,
					Placeholder: field.Placeholder,
					Default:     field.Default,
					Pattern:     field.Pattern,
					Options:     field.Options,
				},
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
		Version:     pluginDescription.Version,
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
				Field: &pb.FormField{
					Name:        name,
					Type:        parameter.Type,
					Description: parameter.Description,
					Placeholder: parameter.Placeholder,
					Required:    parameter.Required,
					Default:     parameter.Default,
					Pattern:     parameter.Pattern,
					Options:     parameter.Options,
				},
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

func translateConnectionInstances(protoConnections map[string]*pb.ConnectionInstance) (map[string]connections.ConnectionInstance, error) {

	concreteConnections := map[string]connections.ConnectionInstance{}
	for protoName, protoConnection := range protoConnections {
		concreteConnections[protoName] = connections.ConnectionInstance{
			VaultUrl: protoConnection.VaultUrl,
			Name:     protoConnection.Name,
			Id:       protoConnection.Id,
			Token:    protoConnection.Token,
		}
	}

	return concreteConnections, nil
}

func implaceDefaultExecuteActionRequestValues(request *pb.ExecuteActionRequest) {
	if request.Parameters == nil {
		request.Parameters = map[string]string{}
	}

	if request.Connections == nil {
		request.Connections = map[string]*pb.ConnectionInstance{}
	}
}

func (service *PluginGRPCService) ExecuteAction(_ context.Context, request *pb.ExecuteActionRequest) (*pb.ExecuteActionResponse, error) {
	implaceDefaultExecuteActionRequestValues(request)

	actionRequest := plugin.ExecuteActionRequest{
		Name:       request.Name,
		Parameters: request.Parameters,
	}

	rawContext, err := translateActionContext(request)
	if err != nil {
		return nil, err
	}

	connectionInstances, err := translateConnectionInstances(request.Connections)
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

func (service *PluginGRPCService) TestCredentials(_ context.Context, request *pb.TestCredentialsRequest) (*pb.TestCredentialsResponse, error) {

	connectionsToBeValidated, err := translateConnectionInstances(request.Connections)
	if err != nil {
		log.Error("Failed to translate connections, error: ", err)
		return nil, err
	}

	validationResponse, err := service.plugin.TestCredentials(connectionsToBeValidated)
	if err != nil {
		log.Error("Failed to validate connections, err: ", err)
		return nil, err
	}

	return &pb.TestCredentialsResponse{
		AreCredentialsValid:   validationResponse.AreCredentialsValid,
		RawValidationResponse: validationResponse.RawValidationResponse,
	}, nil
}

func NewPluginServiceImplementation(plugin plugin.Implementation) *PluginGRPCService {
	return &PluginGRPCService{plugin: plugin}
}
