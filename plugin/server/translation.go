package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blinkops/blink-sdk/plugin/config"
	"github.com/blinkops/blink-sdk/plugin/connections"
	pb "github.com/blinkops/blink-sdk/plugin/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

func translateActionContext(ctx context.Context, request *pb.ExecuteActionRequest) (map[string]interface{}, error) {
	rawContext := map[string]interface{}{}
	if len(request.Context) > 0 {
		if err := json.Unmarshal(request.Context, &rawContext); err != nil {
			log.Error("Failed to unmarshal action context with error: ", err)
			return nil, err
		}
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// remove unnecessary headers from the metadata and store it on the raw context.
		delete(md, "user-agent")
		delete(md, "content-type")
		delete(md, ":authority")
	}

	return rawContext, nil
}

func translateConnectionInstances(protoConnections map[string]*pb.ConnectionInstance) (map[string]*connections.ConnectionInstance, error) {
	concreteConnections := map[string]*connections.ConnectionInstance{}
	for protoName, protoConnection := range protoConnections {
		concreteConnections[protoName] = &connections.ConnectionInstance{
			Name: protoConnection.Name,
			Id:   protoConnection.Id,
			Data: protoConnection.Data,
		}
	}
	return concreteConnections, nil
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
