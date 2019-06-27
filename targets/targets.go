package targets

import (
	"fmt"
	"strings"

	"github.com/coveooss/credentials-sync/credentials"
	log "github.com/sirupsen/logrus"
)

// Target represents an endpoint where credentials can be synced
type Target interface {
	BaseValidateConfiguration() bool
	GetName() string
	GetTags() map[string]string
	Initialize([]credentials.Credentials) error
	ToString() string
	UpdateListOfCredentials([]credentials.Credentials) error
	UpdateCredentials(credentials.Credentials) error
	ValidateConfiguration() bool
}

// Base contains attributes which are common to all targets
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

// BaseValidateConfiguration validates the target's base attributes
func (targetBase *Base) BaseValidateConfiguration() bool {
	if targetBase.Name == "" {
		log.Errorf("Every target need to define a name")
		return false
	}
	return true
}

// GetName returns the target's name
func (targetBase *Base) GetName() string {
	return targetBase.Name
}

// GetTags returns the target's tags
func (targetBase *Base) GetTags() map[string]string {
	return targetBase.Tags
}

// Configuration contains all configured targets
type Configuration struct {
	JenkinsTargets []*JenkinsTarget `mapstructure:"jenkins"`
}

// AllTargets returns all configured targets
func (config *Configuration) AllTargets() []Target {
	targets := []Target{}
	for _, target := range config.JenkinsTargets {
		targets = append(targets, target)
	}
	return targets
}

// ValidateConfiguration verifies that all targets are correctly configured
func (config *Configuration) ValidateConfiguration() bool {
	configIsOK := true
	for _, target := range config.AllTargets() {
		if !target.ValidateConfiguration() || !target.BaseValidateConfiguration() {
			configIsOK = false
		}
	}
	return configIsOK
}
