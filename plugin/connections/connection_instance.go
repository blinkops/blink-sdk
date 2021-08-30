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

	secretRequest, err := http.NewRequest(http.MethodGet, c.VaultUrl, bytes.NewBuffer(marshalledData))

	if err != nil {
		return nil, err
	}

	secretResponse, err := client.Do(secretRequest)

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
