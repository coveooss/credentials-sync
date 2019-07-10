package targets

import (
	"fmt"
	"strings"

	"github.com/coveooss/credentials-sync/credentials"
	log "github.com/sirupsen/logrus"
)

// Target represents an endpoint where credentials can be synced
type Target interface {
	// Base
	BaseValidateConfiguration() bool
	GetName() string
	GetTags() map[string]string
	ShouldDeleteUnsynced() bool
	ShouldTagUnsynced() bool

	// To implement
	GetExistingCredentials() []string
	Initialize([]credentials.Credentials) error
	ToString() string
	DeleteCredentials(id string) error
	UpdateCredentials(credentials.Credentials) error
	ValidateConfiguration() bool
}

// UpdateListOfCredentials syncs the given list of credentials to the given target
func UpdateListOfCredentials(target Target, listOfCredentials []credentials.Credentials) error {
	isSynced := func(id string) bool {
		for _, credentials := range listOfCredentials {
			if credentials.GetID() == id {
				return true
			}
		}
		return false
	}

	for _, credentials := range listOfCredentials {
		log.Infof("[%s] Syncing %s", target.GetName(), credentials.GetID())
		if err := target.UpdateCredentials(credentials); err != nil {
			return err
		}
	}

	if target.ShouldDeleteUnsynced() {
		log.Debugf("Deleting unsynced credentials from %v", target.GetName())
		for _, existingID := range target.GetExistingCredentials() {
			if !isSynced(existingID) {
				log.Infof("[%s] Deleting %s", target.GetName(), existingID)
				if err := target.DeleteCredentials(existingID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Base contains attributes which are common to all targets
type Base struct {
	DeleteUnsynced bool              `mapstructure:"delete_unsynced"`
	TagUnsynced    bool              `mapstructure:"tag_unsynced"`
	Name           string            `mapstructure:"name"`
	Tags           map[string]string `mapstructure:"tags"`
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
	if targetBase.DeleteUnsynced && targetBase.TagUnsynced {
		log.Errorf("Cannot set both `tag_unsynced` and `delete_unsynced` on %v", targetBase.Name)
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

func (targetBase *Base) ShouldDeleteUnsynced() bool {
	return targetBase.DeleteUnsynced
}

func (targetBase *Base) ShouldTagUnsynced() bool {
	return targetBase.TagUnsynced
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
