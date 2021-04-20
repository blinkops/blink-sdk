package config

import (
	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type Config struct {
	Plugin struct {
		PluginDescriptionFilePath     string `yaml:"plugin_description_file_path"`
		ActionsFolderPath             string `yaml:"actions_folder_path"`
		ProviderConfigurationFilePath string `yaml:"provider_configuration_path"`
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
)

func loadConfigurationFromDisk(configFilePath string) (*Config, error) {

	rawYamlBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		logrus.Error("Failed to read action file ", err)
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(rawYamlBytes, &config)
	if err != nil {
		logrus.Error("Failed to unmarshal configuration ", err)
		return nil, err
	}

	return &config, nil
}

func GetConfig() *Config {

	configInitGuard.Do(func() {
		configFilePath := os.Getenv(ConfigurationPathEnvVar)
		if configFilePath == "" {
			panic("Plugin configuration path not supplied")
		}

		loadedConfiguration, err := loadConfigurationFromDisk(configFilePath)

		if err != nil {
			panic(err)
		}

		configInstance = loadedConfiguration
	})

	return configInstance
}
