package credentials

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type LocalSource struct {
	File string
}

func (source *LocalSource) Credentials() ([]Credentials, error) {
	return getCredentialsFromFile(source.File)
}

func (source *LocalSource) Type() string {
	return "Local file"
}

func (source *LocalSource) ValidateConfiguration() bool {
	if _, err := os.Stat(source.File); os.IsNotExist(err) {
		log.Errorf("%s does not exist\n", source.File)
		return false
	}
	return true
}

func getCredentialsFromFile(fileName string) ([]Credentials, error) {
	var (
		err         error
		fileContent []byte
	)
	if fileContent, err = ioutil.ReadFile(fileName); err != nil {
		return nil, err
	}
	return getCredentialsFromBytes(fileContent)
}
