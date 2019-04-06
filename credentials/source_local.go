package credentials

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type LocalSource struct {
	File *string `yaml:"file"`
}

func (source *LocalSource) Credentials() ([]Credentials, error) {
	var (
		err         error
		fileContent []byte
		yamlContent []map[string]interface{}
	)
	if fileContent, err = ioutil.ReadFile(*source.File); err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(fileContent, &yamlContent); err != nil {
		return nil, err
	}
	return ParseCredentials(yamlContent)
}

func (source *LocalSource) Type() string {
	return "Local file"
}

func (source *LocalSource) ValidateConfiguration() bool {
	if _, err := os.Stat(*source.File); os.IsNotExist(err) {
		log.Errorf("%s does not exist\n", *source.File)
		return false
	}
	return true
}
