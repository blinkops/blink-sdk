package server

import (
	"context"
	"github.com/blinkops/plugin-sdk/plugin"
	pb "github.com/blinkops/plugin-sdk/plugin/proto"
)

type PluginGRPCService struct {
	pb.UnimplementedPluginServer

	plugin plugin.Implementation
}

func (service *PluginGRPCService) Describe(ctx context.Context, empty *pb.Empty) (*pb.PluginDescription, error) {
	pluginDescription := service.plugin.Describe()
	actions, _ := service.GetActions(ctx, empty)

	return &pb.PluginDescription{
		Name:        pluginDescription.Name,
		Description: pluginDescription.Description,
		Tags:        pluginDescription.Tags, Provider: pluginDescription.Provider,
		Actions:  	 actions.Actions,
	}, nil
}

func (service *PluginGRPCService) GetActions(ctx context.Context, empty *pb.Empty) (*pb.ActionList, error) {

	actions := service.plugin.GetActions()

	var protoActions []*pb.Action
	for _, action := range actions {

		protoAction := &pb.Action{
			Name:        action.Name,
			Description: action.Description,
			Active:     action.Enabled,
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

func (service *PluginGRPCService) ExecuteAction(_ context.Context, request *pb.ExecuteActionRequest) (*pb.ExecuteActionResponse, error) {

	actionRequest := plugin.ExecuteActionRequest{
		Name:       request.Name,
		Parameters: request.Parameters,
	}

	response, err := service.plugin.ExecuteAction(&actionRequest)
	if err != nil {
		return nil, err
	}

	res := &pb.ExecuteActionResponse{
		ErrorCode: response.ErrorCode,
		Result:    response.Result,
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
