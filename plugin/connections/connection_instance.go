package connections

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	connectionTypeKey = "type"
	connectionNameKey = "id"
	tokenKey          = "token"
)

// ConnectionInstance represents the
type ConnectionInstance struct {
	VaultUrl string
	Name     string
	Id       string
	Token    string
}

func (c *ConnectionInstance) ResolveCredentials() (map[string]interface{}, error) {
	secretResponse, err := c.requestSecret()

	if err != nil {
		return nil, err
	}

	defer func() { _ = secretResponse.Body.Close() }()

	body, err := ioutil.ReadAll(secretResponse.Body)

	if err != nil {
		return nil, err
	}

	internalSecretData := map[string]interface{}{}

	if err = json.Unmarshal(body, &internalSecretData); err != nil {
		return nil, err
	}

	return internalSecretData, nil
}

func (c *ConnectionInstance) requestSecret() (*http.Response, error) {
	client := &http.Client{}
	connectionData := map[string]string{
		connectionTypeKey: c.Name,
		connectionNameKey: c.Id,
		tokenKey:          c.Token,
	}

	marshalledData, err := json.Marshal(connectionData)

	if err != nil {
		return nil, err
	}

	secretRequest, err := http.NewRequest(http.MethodPost, c.VaultUrl, bytes.NewBuffer(marshalledData))

	if err != nil {
		return nil, err
	}

	secretResponse, err := client.Do(secretRequest)

	if err != nil {
		return nil, err
	}

	return secretResponse, nil
}
