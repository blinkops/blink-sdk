package description

import (
	"github.com/blinkops/plugin-sdk/plugin"
	"github.com/go-yaml/yaml"
	log "github.com/sirupsen/logrus"
	"os"
)

func LoadPluginDescriptionFromDisk(descriptionFilePath string) (*plugin.Description, error) {

	rawYamlBytes, err := os.ReadFile(descriptionFilePath)
	if err != nil {
		log.Error("Failed to read description file: ", err)
		return nil, err
	}

	pluginDescription := &plugin.Description{}
	err = yaml.Unmarshal(rawYamlBytes, pluginDescription)
	if err != nil {
		log.Error("Failed to unmarshal raw yaml into description description: ", err)
		return nil, err
	}

	log.Infoln("Loaded description description into memory")
	return pluginDescription, nil
}
