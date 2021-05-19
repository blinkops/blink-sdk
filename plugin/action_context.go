package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blinkops/plugin-sdk/plugin/context"
	"github.com/sirupsen/logrus"
)

type ActionContext struct {
	// Context
	internalContext map[string]interface{} `json:"raw_context"`

	// TODO: Connections (Credentials)

	// Logging
	logger    *logrus.Logger
	logBuffer *bytes.Buffer
}

func NewActionContext(context map[string]interface{}) *ActionContext {

	logBuffer := bytes.Buffer{}

	logger := logrus.New()
	logger.Out = &logBuffer

	return &ActionContext{
		internalContext: context,
		logger:          logger,
		logBuffer:       &logBuffer,
	}
}

func (ctx *ActionContext) GetValue(key string) (interface{}, error) {
	value, err := context.Get(key, ctx.internalContext)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("no entry with name %s found, error: %v", key, err))
	}

	return value, nil
}

func (ctx *ActionContext) SetValue(key string, value interface{}) {
	err := context.Set(key, value, ctx.internalContext)
	if err != nil {
		return
	}
}

func (ctx *ActionContext) DeleteEntry(key string) {
	err := context.Delete(key, ctx.internalContext)
	if err != nil {
		return
	}
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
