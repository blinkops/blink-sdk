package plugin

import (
	"errors"
	"github.com/blinkops/blink-sdk/plugin/connections"
	"time"

	"encoding/json"
	log "github.com/sirupsen/logrus"
)

const DefaultTimeout = time.Second

type ActionParameter struct {
	Type        string   `yaml:"type"`
	Description string   `yaml:"description"`
	Placeholder string   `yaml:"placeholder"`
	Required    bool     `yaml:"required"`
	Default     string   `yaml:"default"`
	Pattern     string   `yaml:"pattern"` // optional: regex to validate in case of input component
	Options     []string `yaml:"options"` // optional: the option list in case of dropdown\checkbox
}

type Action struct {
	Name        string                     `yaml:"name"`
	Description string                     `yaml:"description"`
	Enabled     bool                       `yaml:"enabled"`
	EntryPoint  string                     `yaml:"entry_point"`
	Parameters  map[string]ActionParameter `yaml:"parameters"`
	Output      *Output
}

type Field struct {
	Name string
	Type string
}

type Output struct {
	Name   string
	Fields []Field
}

type Description struct {
	Name        string                            `yaml:"name"`
	Description string                            `yaml:"description"`
	Tags        []string                          `yaml:"tags"`
	Provider    string                            `yaml:"provider"`
	Connections map[string]connections.Connection `yaml:"connection_types"`
	Version     string                            `yaml:"version"`
}

type ExecuteActionRequest struct {
	Name       string            `yaml:"name"`
	Parameters map[string]string `yaml:"parameters"`
}

type ExecuteActionResponse struct {
	ErrorCode int64  `yaml:"error_code"`
	Result    []byte `yaml:"result"`
	Rows      []map[string]string
}

type ProviderConfiguration struct {
	ConfigurationMap map[interface{}]interface{}
}

type CredentialsValidationResponse struct {
	AreCredentialsValid   bool
	RawValidationResponse []byte
}

type Implementation interface {
	Describe() Description

	GetActions() []Action
	ExecuteAction(context *ActionContext, request *ExecuteActionRequest) (*ExecuteActionResponse, error)

	TestCredentials(map[string]connections.ConnectionInstance) (*CredentialsValidationResponse, error)
}

func (req *ExecuteActionRequest) GetParameters() (map[string]string, error) {
	_, ok := req.Parameters["parameters_as_json"]
	if ok {
		return nil, errors.New(ErrParametersAsJsonProvided)
	}
	return req.Parameters, nil
}

func (req *ExecuteActionRequest) GetUnmarshalledParameters() (map[string]interface{}, error) {
	parametersAsJson, ok := req.Parameters["parameters_as_json"]
	if !ok {
		return nil, errors.New(ErrNoParametersAsJsonProvided)
	}

	actionParameters := make(map[string]interface{})
	err := json.Unmarshal([]byte(parametersAsJson), &actionParameters)
	if err != nil {
		log.Error("Failed to unmarshal action parameters, err: ", err)
		return nil, err
	}

	return actionParameters, nil
}
