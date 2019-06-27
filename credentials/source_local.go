package credentials

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

// LocalSource represents local files containing credentials
type LocalSource struct {
	File string
}

// Credentials extracts credentials from the source
func (source *LocalSource) Credentials() ([]Credentials, error) {
	return getCredentialsFromFile(source.File)
}

// Type returns the type of the source
func (source *LocalSource) Type() string {
	return "Local file"
}

// ValidateConfiguration verifies that the source's attributes are valid
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
