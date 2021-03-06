package targets

import (
	"fmt"
	"strings"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/hashicorp/go-multierror"
)

// Target represents an endpoint where credentials can be synced
type Target interface {
	// Base
	BaseValidateConfiguration() error
	GetName() string
	GetTags() map[string]string
	ShouldDeleteUnsynced() bool
	ShouldTagUnsynced() bool // Not implemented

	// To implement
	GetExistingCredentials() []string
	Initialize([]credentials.Credentials) error
	ToString() string
	DeleteCredentials(id string) error
	UpdateCredentials(credentials.Credentials) error
	ValidateConfiguration() error
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
func (targetBase *Base) BaseValidateConfiguration() error {
	if targetBase.Name == "" {
		return fmt.Errorf("Every target need to define a name")
	}
	if targetBase.DeleteUnsynced && targetBase.TagUnsynced {
		return fmt.Errorf("Cannot set both `tag_unsynced` and `delete_unsynced` on %v", targetBase.Name)
	}
	return nil
}

// GetName returns the target's name
func (targetBase *Base) GetName() string {
	return targetBase.Name
}

// GetTags returns the target's tags
func (targetBase *Base) GetTags() map[string]string {
	return targetBase.Tags
}

// ShouldDeleteUnsynced returns true if the unsynced credentials should be deleted from the target
func (targetBase *Base) ShouldDeleteUnsynced() bool {
	return targetBase.DeleteUnsynced
}

// ShouldTagUnsynced returns true if the unsynced credentials should be tagged accordingly on the target
func (targetBase *Base) ShouldTagUnsynced() bool {
	return targetBase.TagUnsynced
}

// HasCredential returns true if the given ID is found on the target
func HasCredential(target Target, id string) bool {
	for _, existingID := range target.GetExistingCredentials() {
		if existingID == id {
			return true
		}
	}
	return false
}

// TargetCollection represents a collection of targets to which credentials can be synced
type TargetCollection interface {
	AllTargets() []Target
	ValidateConfiguration() error
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
func (config *Configuration) ValidateConfiguration() error {
	var validationErrors error
	for _, target := range config.AllTargets() {
		if err := target.BaseValidateConfiguration(); err != nil {
			validationErrors = multierror.Append(validationErrors, err)
		}
		if err := target.ValidateConfiguration(); err != nil {
			validationErrors = multierror.Append(validationErrors, err)
		}
	}
	return validationErrors
}
