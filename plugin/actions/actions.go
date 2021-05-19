package actions

import (
	"github.com/blinkops/plugin-sdk/plugin"
	"github.com/go-yaml/yaml"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func loadActionFromDisk(actionPath string) (*plugin.Action, error) {

	rawYamlBytes, err := os.ReadFile(actionPath)
	if err != nil {
		log.Error("Failed to read action file ", err)
		return nil, err
	}

	action := plugin.Action{}
	err = yaml.Unmarshal(rawYamlBytes, &action)
	if err != nil {
		log.Error("Failed to unmarshal action ", err)
		return nil, err
	}

	return &action, nil
}

func LoadActionsFromDisk(actionFoldersPath string) ([]plugin.Action, error) {

	var actions []plugin.Action
	err := filepath.Walk(actionFoldersPath, func(filePath string, info os.FileInfo, err error) error {

		if !strings.HasSuffix(filePath, ".yaml") {
			return nil
		}

		loadedAction, err := loadActionFromDisk(filePath)
		if err != nil {
			log.Error("Failed to load action from disk ", err)
			return err
		}

		actions = append(actions, *loadedAction)
		return nil
	})

	if err != nil {
		log.Error("Failed to load actions from disk ", err)
		return make([]plugin.Action, 0), err
	}

	log.Infof("Loaded %d action from disk!", len(actions))
	return actions, nil
}
