package connections

import (
	"github.com/go-yaml/yaml"
	log "github.com/sirupsen/logrus"
	"os"
)

type ConnectionField struct {
	Name        string   `yaml:"name"`
	FieldType   string   `yaml:"field_type"`
	InputType   string   `yaml:"input_type"`
	Required    bool     `yaml:"required"`
	Description string   `yaml:"description"`
	Placeholder string   `yaml:"placeholder"`
	Default     string   `yaml:"default"`
	Pattern     string   `yaml:"pattern"`
	Options     []string `yaml:"options"`
}

// Connection represents a connection type requested by the user.
// For each requested connection an instance of it will be sent
// for every action executed.
type Connection struct {
	Name      string
	Fields    map[string]ConnectionField
	Reference string
}

type RequestedConnections struct {
	ConnectionTypeReferences map[string]string                     `yaml:"connection_type_references"`
	ConnectionTypes          map[string]map[string]ConnectionField `yaml:"connection_types"`
}

func LoadConnectionsFromDisk(connectionsFilePath string) (map[string]Connection, error) {

	rawYamlBytes, err := os.ReadFile(connectionsFilePath)
	if err != nil {
		log.Error("Failed to read connections file: ", err)
		return nil, err
	}

	requestedConnections := &RequestedConnections{}
	err = yaml.Unmarshal(rawYamlBytes, requestedConnections)
	if err != nil {
		log.Error("Failed to unmarshal raw yaml into requested connections: ", err)
		return nil, err
	}

	connections := map[string]Connection{}
	for connectionName, connectionFields := range requestedConnections.ConnectionTypes {
		typeReference, _ := requestedConnections.ConnectionTypeReferences[connectionName]

		connections[connectionName] = Connection{
			Name:      connectionName,
			Fields:    connectionFields,
			Reference: typeReference,
		}
	}

	log.Infoln("Loaded requested connections into memory")
	return connections, nil
}
