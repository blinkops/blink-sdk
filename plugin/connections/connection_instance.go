package connections

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	vaultapi "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	vaultClientTimeout = time.Second * 60
	vaultClientRetries = 3

	vaultSecretPathFormat = "secret/data/%s/%s"
	secretDataKey         = "data"
)

// ConnectionInstance represents the
type ConnectionInstance struct {
	VaultUrl string
	Name     string
	Id       string
	Token    string
}

func (c *ConnectionInstance) acquireVaultAPIConnection() (*vaultapi.Client, error) {
	vaultClient, err := vaultapi.NewClient(&vaultapi.Config{
		Address:    c.VaultUrl,
		HttpClient: cleanhttp.DefaultPooledClient(),
		Timeout:    vaultClientTimeout,
		MaxRetries: vaultClientRetries,
		Backoff:    retryablehttp.LinearJitterBackoff,
	})

	if err != nil {
		log.Error("Failed to create VaultAPI client, error: ", err)
		return nil, err
	}

	vaultClient.SetToken(c.Token)
	return vaultClient, nil
}

func (c *ConnectionInstance) ResolveCredentials() (map[string]interface{}, error) {

	client, err := c.acquireVaultAPIConnection()
	if err != nil {
		return nil, err
	}

	vaultSecretPath := fmt.Sprintf(vaultSecretPathFormat, c.Name, c.Id)
	secrets, err := client.Logical().Read(vaultSecretPath)
	if err != nil {
		log.Errorf("Failed to read secret credentials with path %s, error: %v", vaultSecretPath, err)
		return nil, err
	}

	internalSecretData, ok := secrets.Data[secretDataKey].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid secret structure, internal data missing")
	}

	return internalSecretData, nil
}
