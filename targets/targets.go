package targets

import (
	"fmt"
	"strings"

	"github.com/coveo/credentials-sync/credentials"
	log "github.com/sirupsen/logrus"
)

type Target interface {
	BaseValidateConfiguration() bool
	GetName() string
	Initialize([]credentials.Credentials) error
	ToString() string
	UpdateListOfCredentials([]credentials.Credentials) error
	UpdateCredentials(credentials.Credentials) error
	ValidateConfiguration() bool
}

type Base struct {
	Name string
	Tags map[string]string
}

// BaseToString prints out the target fields common to all types of targets
func (targetBase *Base) BaseToString() string {
	tagString := ""
	for tagKey, tagValue := range targetBase.Tags {
		tagString = fmt.Sprintf("%s %s=%s", tagString, tagKey, tagValue)
	}
	return fmt.Sprintf("%s [Tags: %s]", targetBase.Name, strings.TrimSpace(tagString))
}

func (targetBase *Base) BaseValidateConfiguration() bool {
	if targetBase.Name == "" {
		log.Errorf("Every target need to define a name")
		return false
	}
	return true
}

type Configuration struct {
	JenkinsTargets []*JenkinsTarget `mapstructure:"jenkins"`
}

func (config *Configuration) AllTargets() []Target {
	targets := []Target{}
	for _, target := range config.JenkinsTargets {
		targets = append(targets, target)
	}
	return targets
}

func (config *Configuration) ValidateConfiguration() bool {
	configIsOK := true
	for _, target := range config.AllTargets() {
		if !target.ValidateConfiguration() || !target.BaseValidateConfiguration() {
			configIsOK = false
		}
	}
	return configIsOK
}
