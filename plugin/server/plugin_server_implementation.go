package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blinkops/blink-sdk/plugin"
	"github.com/blinkops/blink-sdk/plugin/assets"
	"github.com/blinkops/blink-sdk/plugin/config"
	"github.com/blinkops/blink-sdk/plugin/connections"
	pb "github.com/blinkops/blink-sdk/plugin/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
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
			Name:      connectionName,
			Fields:    protoConnectionFields,
			Reference: connection.Reference,
		}
	}

	return protoConnections
}

func translatePluginType() pb.PluginDescription_PluginType {
	if config.GetConfig() == nil {
		return pb.PluginDescription_SHARED
	}
	pluginType := config.GetConfig().Plugin.Type

	switch pluginType {
	case config.SharedPluginType:
		return pb.PluginDescription_SHARED
	case config.PrivatePluginType:
		return pb.PluginDescription_PRIVATE
	}

	panic(fmt.Sprintf("Invalid plugin type configured %s", pluginType))
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
		PluginType:  translatePluginType(),
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
					Index:       parameter.Index,
					Format:      parameter.Format,
					IsMulti:     parameter.IsMulti,
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

func translateActionContext(ctx context.Context, request *pb.ExecuteActionRequest) (map[string]interface{}, error) {

	rawContext := map[string]interface{}{}
	if len(request.Context) > 0 {
		err := json.Unmarshal(request.Context, &rawContext)
		if err != nil {
			log.Error("Failed to unmarshal action context with error: ", err)
			return nil, err
		}
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		// remove unnecessary headers from the metadata and store it on the raw context.
		delete(md, "user-agent")
		delete(md, "content-type")
		delete(md, ":authority")
		rawContext[connections.MetadataHeader] = md

	}

	return rawContext, nil
}

func translateConnectionInstances(protoConnections map[string]*pb.ConnectionInstance) (map[string]*connections.ConnectionInstance, error) {

	concreteConnections := map[string]*connections.ConnectionInstance{}
	for protoName, protoConnection := range protoConnections {
		concreteConnections[protoName] = &connections.ConnectionInstance{
			VaultUrl: protoConnection.VaultUrl,
			Name:     protoConnection.Name,
			Id:       protoConnection.Id,
			Token:    protoConnection.Token,
		}
	}

	return concreteConnections, nil
}

func emplaceDefaultExecuteActionRequestValues(request *pb.ExecuteActionRequest) {
	if request.Parameters == nil {
		request.Parameters = map[string]string{}
	}

	if request.Connections == nil {
		request.Connections = map[string]*pb.ConnectionInstance{}
	}
}

func (service *PluginGRPCService) ExecuteAction(ctx context.Context, request *pb.ExecuteActionRequest) (*pb.ExecuteActionResponse, error) {
	emplaceDefaultExecuteActionRequestValues(request)

	actionRequest := plugin.ExecuteActionRequest{
		Name:       request.Name,
		Parameters: request.Parameters,
		Timeout:    request.Timeout,
	}

	rawContext, err := translateActionContext(ctx, request)
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
		ErrorCode:    response.ErrorCode,
		Result:       response.Result,
		Context:      updatedActionContext,
		LogBuffer:    actionContext.GetRawLogBuffer(),
		ErrorMessage: response.ErrorMessage,
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

func (service *PluginGRPCService) HealthProbe(context.Context, *pb.Empty) (*pb.HealthStatus, error) {
	return &pb.HealthStatus{}, nil
}

func (service *PluginGRPCService) GetAssets(context.Context, *pb.Empty) (*pb.Assets, error) {

	pluginIconBuffer, err := assets.ReadPluginIconBufferIntoMemory()
	if err != nil {
		return nil, err
	}

	return &pb.Assets{Icon: &pb.PluginIcon{RawIconBuffer: pluginIconBuffer}}, nil
}

func NewPluginServiceImplementation(plugin plugin.Implementation) *PluginGRPCService {
	return &PluginGRPCService{plugin: plugin}
}
