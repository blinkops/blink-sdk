package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blinkops/plugin-sdk/plugin/connections"
	"github.com/sirupsen/logrus"
)

type ActionContext struct {
	// Context
	internalContext map[string]interface{} `json:"raw_context"`

	connections map[string]connections.ConnectionInstance

	// Logging
	logger    *logrus.Logger
	logBuffer *bytes.Buffer
}

func NewActionContext(context map[string]interface{}, connections map[string]connections.ConnectionInstance) *ActionContext {

	logBuffer := bytes.Buffer{}

	logger := logrus.New()
	logger.Out = &logBuffer

	return &ActionContext{
		internalContext: context,
		logger:          logger,
		logBuffer:       &logBuffer,
		connections:     connections,
	}
}

func (ctx *ActionContext) GetValue(key string) (interface{}, error) {
	value, ok := ctx.internalContext[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no entry with name %s found", key))
	}

	return value, nil
}

func (ctx *ActionContext) SetValue(key string, value interface{}) {
	ctx.internalContext[key] = value
}

func (ctx *ActionContext) DeleteEntry(key string) {
	delete(ctx.internalContext, key)
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

func (ctx *ActionContext) GetLogger() *logrus.Logger {
	return ctx.logger
}

func (ctx *ActionContext) GetCredentials(name string) (map[string]interface{}, error) {
	connectionInstance, ok := ctx.connections[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("connection with %s is missing in action context", name))
	}

	return connectionInstance.ResolveCredentials()
}
