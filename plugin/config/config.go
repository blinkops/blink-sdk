package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

type Config struct {
	Plugin struct {
		PluginDescriptionFilePath     string `yaml:"plugin_description_file_path"`
		ActionsFolderPath             string `yaml:"actions_folder_path"`
		ProviderConfigurationFilePath string `yaml:"provider_configuration_path"`
		Type                          string `yaml:"type"`
	}

	Server struct {
		Port string `yaml:"port"`
	}
}

var (
	configInstance  *Config
	configInitGuard sync.Once
)

const (
	ConfigurationPathEnvVar = "CONFIG_FILE_PATH"
	SharedPluginType        = "shared"
	PrivatePluginType       = "private"
)

func loadConfigurationFromDisk(configFilePath string) (*Config, error) {

	rawYamlBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Error("Failed to read action file ", err)
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(rawYamlBytes, &config)
	if err != nil {
		log.Error("Failed to unmarshal configuration ", err)
		return nil, err
	}

	return &config, nil
}

func validateConfiguration(config *Config) {

	pluginType := config.Plugin.Type
	if pluginType == "" {
		config.Plugin.Type = SharedPluginType
	}

	if pluginType != SharedPluginType && pluginType != PrivatePluginType {
		panic(fmt.Sprintf("Invalid plugin type: %s", config.Plugin.Type))
	}

}

func GetConfig() *Config {

	configInitGuard.Do(func() {

		configFilePath := os.Getenv(ConfigurationPathEnvVar)
		if configFilePath == "" {
			log.Warn("Plugin configuration path not supplied")
			return
		}

		loadedConfiguration, err := loadConfigurationFromDisk(configFilePath)
		if err != nil {
			panic(err)
		}

		validateConfiguration(loadedConfiguration)

		configInstance = loadedConfiguration
	})

	return configInstance
}

func GetServerPort() string {
	config := GetConfig()
	if config == nil {
		return "1337"
	}
	return config.Server.Port
}
