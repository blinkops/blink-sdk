package server

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/blinkops/blink-sdk/plugin"
	pb "github.com/blinkops/blink-sdk/plugin/proto"
	log "github.com/sirupsen/logrus"
)

type PluginGRPCService struct {
	pb.UnimplementedPluginServer

	plugin plugin.Implementation

	activeWorkers      int64
	activeWorkersMutex sync.Mutex

	lastUse int64
}

func (service *PluginGRPCService) Describe(ctx context.Context, empty *pb.Empty) (*pb.PluginDescription, error) {
	pluginDescription := service.plugin.Describe()
	actions, err := service.GetActions(ctx, empty)
	if err != nil {
		return nil, err
	}

	return &pb.PluginDescription{
		Name:                 pluginDescription.Name,
		IconUri:              pluginDescription.IconUri,
		Description:          pluginDescription.Description,
		Tags:                 pluginDescription.Tags,
		Provider:             pluginDescription.Provider,
		Actions:              actions.Actions,
		Connections:          translateToProtoConnections(pluginDescription.Connections),
		Version:              pluginDescription.Version,
		PluginType:           translatePluginType(),
		IsConnectionOptional: pluginDescription.IsConnectionOptional,
	}, nil
}

func (service *PluginGRPCService) GetActions(_ context.Context, _ *pb.Empty) (*pb.ActionList, error) {
	actions := service.plugin.GetActions()

	var protoActions []*pb.Action
	for _, action := range actions {
		protoAction := &pb.Action{
			Name:                 action.Name,
			IconUri:              action.IconUri,
			DisplayName:          action.DisplayName,
			CollectionName:       action.CollectionName,
			Description:          action.Description,
			Active:               action.Enabled,
			Connections:          translateToProtoConnections(action.Connections),
			IsConnectionOptional: action.IsConnectionOptional,
		}

		var protoParameters []*pb.ActionParameter
		for name, parameter := range action.Parameters {
			protoParameter := &pb.ActionParameter{
				Field: &pb.FormField{
					Name:        name,
					DisplayName: parameter.DisplayName,
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

func emplaceDefaultExecuteActionRequestValues(request *pb.ExecuteActionRequest) {
	if request.Parameters == nil {
		request.Parameters = map[string]string{}
	}

	if request.Connections == nil {
		request.Connections = map[string]*pb.ConnectionInstance{}
	}
}

func (service *PluginGRPCService) ExecuteAction(ctx context.Context, request *pb.ExecuteActionRequest) (*pb.ExecuteActionResponse, error) {
	defer func() {
		if recoveryInformation := recover(); recoveryInformation != nil {
			log.Errorf("[PANIC RECOVERY] Successfully recovered from panic inside ExecuteAction, Information: \"%v\", Trace: \"%s\"", recoveryInformation, string(debug.Stack()))
		}
	}()

	service.activateWorker()
	defer service.deactivateWorker()

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
	defer func() {
		if recoveryInformation := recover(); recoveryInformation != nil {
			log.Errorf("[PANIC RECOVERY] Successfully recovered from panic inside TestCredentials, Information: \"%v\", Trace: \"%s\"", recoveryInformation, string(debug.Stack()))
		}
	}()

	service.activateWorker()
	defer service.deactivateWorker()

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
	return &pb.HealthStatus{
		LastUse: service.getLastUse(),
	}, nil
}

func (service *PluginGRPCService) getLastUse() int64 {
	service.activeWorkersMutex.Lock()
	defer service.activeWorkersMutex.Unlock()

	if service.activeWorkers > 0 {
		return timeNowNano()
	}

	if pluginStatus, err := service.plugin.HealthStatus(); err == nil && pluginStatus.InUse {
		return timeNowNano()
	}

	return service.lastUse
}

func (service *PluginGRPCService) activateWorker() {
	service.activeWorkersMutex.Lock()
	defer service.activeWorkersMutex.Unlock()
	service.activeWorkers++

	service.lastUse = timeNowNano() // Update last use on worker activation
}

func (service *PluginGRPCService) deactivateWorker() {
	service.activeWorkersMutex.Lock()
	defer service.activeWorkersMutex.Unlock()
	service.activeWorkers--

	service.lastUse = timeNowNano() // Update last use on worker de-activation
}

func NewPluginServiceImplementation(plugin plugin.Implementation) *PluginGRPCService {
	return &PluginGRPCService{
		plugin:             plugin,
		lastUse:            timeNowNano(),
		activeWorkersMutex: sync.Mutex{},
	}
}

func timeNowNano() int64 {
	return time.Now().UnixNano()
}
