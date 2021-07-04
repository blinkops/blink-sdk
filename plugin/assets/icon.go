package assets

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

const (
	pluginIconPath = "assets/icon.svg"
)

func ReadPluginIconBufferIntoMemory() ([]byte, error) {

	iconBuffer, err := ioutil.ReadFile(pluginIconPath)
	if err != nil {
		log.Errorf("Failed to read icon into memory with error %v", err)
		return nil, err
	}

	return iconBuffer, nil
}
