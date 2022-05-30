package plugin

import (
	"time"

	"github.com/blinkops/blink-sdk/plugin/connections"
)

type MockPlugin struct {
}

func (MockPlugin) Describe() Description {
	return Description{Name: "Mock", Description: "Im mock"}
}

func (MockPlugin) GetActions() []Action {
	return []Action{}
}

func (MockPlugin) ExecuteAction(context *ActionContext, request *ExecuteActionRequest) (*ExecuteActionResponse, error) {
	time.Sleep(time.Duration(request.Timeout) * time.Second)
	return &ExecuteActionResponse{ErrorCode: 0, Result: []byte("Executed Action")}, nil
}

func (MockPlugin) TestCredentials(map[string]*connections.ConnectionInstance) (*CredentialsValidationResponse, error) {
	return &CredentialsValidationResponse{AreCredentialsValid: false, RawValidationResponse: []byte("Not Supported")}, nil
}
