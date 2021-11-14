package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blinkops/blink-sdk/plugin/connections"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"strings"
)

type ActionContext struct {
	// Context
	internalContext map[string]interface{} `json:"raw_context"`

	// Connections
	connections map[string]*connections.ConnectionInstance

	// Logging
	logger    *log.Logger
	logBuffer *bytes.Buffer
}

func NewActionContext(context map[string]interface{}, connections map[string]*connections.ConnectionInstance) *ActionContext {

	logBuffer := bytes.Buffer{}

	logger := log.New()
	logger.Out = &logBuffer

	return &ActionContext{
		internalContext: context,
		logger:          logger,
		logBuffer:       &logBuffer,
		connections:     connections,
	}
}

func resolveInnerItems(key string, createKeys bool, innerContext map[string]interface{}) (map[string]interface{}, error) {
	pathToInnerKey := strings.Split(key, ".")

	depth := len(pathToInnerKey) - 1
	if depth == 0 || pathToInnerKey[1] == "" {
		return nil, fmt.Errorf("provided key is not allowed, must be at-least 1 depth, %v", key)
	}

	innerContextIterator := interface{}(innerContext)
	for i := 1; i < depth; i++ {
		innerContextIteratorMap, ok := innerContextIterator.(map[string]interface{})
		if !ok {
			return nil, errors.New("failed to convert innerContextIterator to map[string] interface")
		}
		currentHead := pathToInnerKey[i]
		switch innerContextIteratorMap[currentHead].(type) {
		case map[string]interface{}:
			break
		default:
			if createKeys {
				innerContextIteratorMap[currentHead] = make(map[string]interface{})
			} else {
				return nil, errors.New("given path doesn't exists")
			}
		}
		innerContextIterator = innerContextIteratorMap[currentHead]
	}

	innerContextMap, ok := innerContextIterator.(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to convert last innerContextIterator to map[string] interface")
	}

	return innerContextMap, nil
}

func (ctx *ActionContext) getInnerContext(key string, createKeys bool) (map[string]interface{}, error) {
	// Usage: key will be path (json.dot.walking) that the start of the parameter
	// 		  strings.Split(parameter, ".")[0] -> Should be `inputs` or `variables`
	pathToInnerKey := strings.Split(key, ".")
	if len(pathToInnerKey) == 0 {
		return nil, fmt.Errorf("provided key is not allowed, must be at-least 1 depth, %v", key)
	}

	innerContextInterface, ok := ctx.internalContext[pathToInnerKey[0]]
	if !ok {
		return nil, fmt.Errorf("provided key is not allowed, tried to acces unknown path, %v", key)
	}

	innerContext, ok := innerContextInterface.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to access path, %v", key)
	}

	if innerContext == nil {
		return nil, fmt.Errorf("failed to get inner context with key: %v", key)
	}

	return resolveInnerItems(key, createKeys, innerContext)
}

func (ctx *ActionContext) GetValue(key string) (interface{}, error) {
	innerContext, err := ctx.getInnerContext(key, false)
	if err != nil || innerContext == nil {
		return nil, fmt.Errorf("no entry with name %s found, error: %v", key, err)
	}
	pathToInnerKey := strings.Split(key, ".")
	return innerContext[pathToInnerKey[len(pathToInnerKey)-1]], nil
}

func (ctx *ActionContext) SetValue(key string, value interface{}) {
	innerContext, err := ctx.getInnerContext(key, true)
	if err != nil || innerContext == nil {
		log.Errorf("failed to set entry with name %s, error: %v", key, err)
		return
	}
	pathToInnerKey := strings.Split(key, ".")
	innerContext[pathToInnerKey[len(pathToInnerKey)-1]] = value
}

func (ctx *ActionContext) DeleteEntry(key string) {
	innerContext, err := ctx.getInnerContext(key, true)
	if err != nil || innerContext == nil {
		log.Errorf("failed to delete entry with name %s, error: %v", key, err)
		return
	}
	pathToInnerKey := strings.Split(key, ".")
	delete(innerContext, pathToInnerKey[len(pathToInnerKey)-1])
}

func (ctx *ActionContext) GetMarshaledContext() ([]byte, error) {
	return json.Marshal(ctx.internalContext)
}

func (ctx *ActionContext) GetAllContextEntries() map[string]interface{} {
	return ctx.internalContext
}

func (ctx *ActionContext) GetRawLogBuffer() []byte {
	return ctx.logBuffer.Bytes()
}

func (ctx *ActionContext) GetLogger() *log.Logger {
	return ctx.logger
}

func (ctx *ActionContext) GetCredentials(name string) (map[string]interface{}, error) {
	connectionInstance, ok := ctx.connections[name]
	if !ok {
		// If no connection was provided, use the grpc-metadata.
		md, ok := ctx.internalContext[connections.MetadataHeader]
		if ok {
			conn := map[string]interface{}{}

			// Convert metadata to connection.
			for k, v := range md.(metadata.MD) {
				conn[k] = strings.Join(v, ", ")
			}

			return conn, nil
		}

		return nil, errors.New(fmt.Sprintf("connection with %s is missing in action context", name))
	}

	return connectionInstance.ResolveCredentials()
}

func (ctx *ActionContext) GetAllConnections() map[string]*connections.ConnectionInstance {
	return ctx.connections
}

func (ctx *ActionContext) ReplaceContext(context map[string]interface{}) {
	ctx.internalContext = context
}
