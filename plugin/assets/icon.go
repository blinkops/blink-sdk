package assets

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

const (
	PluginIconPath = "assets/icon.svg"
)

func ReadPluginIconBufferIntoMemory(iconUri string) ([]byte, error) {

	iconBuffer, err := ioutil.ReadFile(iconUri)
	if err != nil {
		log.Errorf("Failed to read icon into memory with error %v", err)
		return nil, err
	}

	return iconBuffer, nil
}
