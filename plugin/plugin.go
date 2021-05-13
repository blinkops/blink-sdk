package plugin

type ActionParameter struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
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
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	Provider    string   `yaml:"provider"`
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

type Implementation interface {
	Describe() Description

	GetActions() []Action
	ExecuteAction(context *ActionContext, request *ExecuteActionRequest) (*ExecuteActionResponse, error)
}
