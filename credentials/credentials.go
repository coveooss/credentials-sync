package credentials

import (
	"fmt"

	"github.com/coveooss/credentials-sync/logger"
	"github.com/hashicorp/go-multierror"
	"github.com/mitchellh/mapstructure"
)

// Credentials defines the methods that can be called by all types of credentials
type Credentials interface {
	BaseValidate() error
	GetID() string
	GetTargetID() string
	ShouldSync(targetName string, targetTags map[string]string) bool
	ToString(bool) string
	Validate() error
}

type targetTagsMatcher struct {
	DoMatch   map[string]interface{} `mapstructure:"do_match"`
	DontMatch map[string]interface{} `mapstructure:"dont_match"`
}

// Base defines that fields that are common to all types of credentials
type Base struct {
	ID          string
	Description string
	NoSync      bool              `mapstructure:"no_sync"`
	TargetName  string            `mapstructure:"target"`
	TargetTags  targetTagsMatcher `mapstructure:"target_tags"`
	TargetID    string            `mapstructure:"target_id"`

	// Field set by constructor
	CredType string

	// For multi-value fields. Such as SSM
	Value string
}

// BaseToString prints out the credentials fields common to all types of credentials
func (credBase *Base) BaseToString() string {
	value := fmt.Sprintf("%s -> Type: %s", credBase.ID, credBase.CredType)
	if credBase.Description != "" {
		value += ", Description: " + credBase.Description
	}
	return value
}

// BaseValidate verifies that the credentials fields common to all types of credentials contain valid values
func (credBase *Base) BaseValidate() error {
	if credBase.ID == "" {
		return fmt.Errorf("credentials (%s) has no defined ID", credBase.BaseToString())
	}
	if credBase.CredType == "" {
		return fmt.Errorf("credentials (%s) has no type. This is probably a bug in the software", credBase.ID)
	}
	return nil
}

// GetDescriptionOrID returns the description if it set, otherwise it returns the ID
func (credBase *Base) GetDescriptionOrID() string {
	if credBase.Description == "" {
		return credBase.ID
	}
	return credBase.Description
}

// GetID returns a credentials' ID
func (credBase *Base) GetID() string {
	return credBase.ID
}

// GetTargetID returns a credentials' Target ID (Essentially, the name that the credentials should have on a target)
// This is helpful to have different credentials with the same target ID (on different targets)
func (credBase *Base) GetTargetID() string {
	if credBase.TargetID != "" {
		return credBase.TargetID
	}
	return credBase.GetID()
}

// ShouldSync returns, given a target's name and tags, if a credentials should be synced to that target
// This is based on various credentials attributes such as the TargetTags DoMatch and DontMatch attributes
func (credBase *Base) ShouldSync(targetName string, targetTags map[string]string) bool {
	if credBase.NoSync {
		return false
	}
	if credBase.TargetName != "" && credBase.TargetName != targetName {
		return false
	}

	findMatch := func(match map[string]interface{}) bool {
		for key, value := range match {
			for tagKey, tag := range targetTags {
				if key != tagKey {
					continue
				}
				if valueAsString, ok := value.(string); ok {
					if valueAsString == tag {
						return true
					}
				} else if valueAsList, ok := value.([]string); ok {
					if listContainsElement(valueAsList, tag) {
						return true
					}
				} else {
					logger.Log.Warningf("%s ignored. Its value should either be a string or a list of string", key)
				}
			}
		}
		return false
	}

	return !findMatch(credBase.TargetTags.DontMatch) && (len(credBase.TargetTags.DoMatch) == 0 || findMatch(credBase.TargetTags.DoMatch))
}

// ParseCredentials transforms a list of maps into a list of Credentials
// The credentials type is determined by the `type` attribute
func ParseCredentials(credentialsMaps []map[string]interface{}) ([]Credentials, error) {
	credentialsList := make([]Credentials, 0)
	for _, credentialsMap := range credentialsMaps {
		newCredentials, err := ParseSingleCredentials(credentialsMap)
		if err != nil {
			return nil, err
		}
		credentialsList = append(credentialsList, newCredentials)
	}
	return credentialsList, nil
}

// ParseSingleCredentials transforms a map into a Credentials struct
// The credentials type is determined by the `type` attribute
func ParseSingleCredentials(credentialsMap map[string]interface{}) (Credentials, error) {
	var credentialsType string
	var credentials Credentials
	var id = credentialsMap["id"]
	if value, ok := credentialsMap["type"]; ok {
		if credentialsType, ok = value.(string); !ok {
			return nil, fmt.Errorf("entry %s: credentials type '%v' is not a string", id, credentialsType)
		}
	} else {
		return nil, fmt.Errorf("entry %s: unable to find the credentials type %s", id, credentialsType)
	}

	switch credentialsType {
	case "aws":
		credentials = NewAmazonWebServicesCredentials()
	case "usernamepassword":
		credentials = NewUsernamePassword()
	case "secret":
		credentials = NewSecretText()
	case "ssh":
		credentials = NewSSHCredentials()
	case "github_app":
		credentials = NewGithubAppCredentials()
	default:
		return nil, fmt.Errorf("entry %s: unknown credentials type: %s", id, credentialsType)
	}
	err := mapstructure.Decode(credentialsMap, credentials)
	if err != nil {
		return nil, fmt.Errorf("entry %s: invalid credentials data: %v", id, err)
	}
	var validationErrors error
	if err := credentials.BaseValidate(); err != nil {
		validationErrors = multierror.Append(validationErrors, err)
	}
	if err := credentials.Validate(); err != nil {
		validationErrors = multierror.Append(validationErrors, err)
	}
	if validationErrors != nil {
		return nil, fmt.Errorf("the following credentials failed to validate: %v -> %v", credentials.ToString(false), validationErrors)
	}
	return credentials, nil
}

func listContainsElement(list []string, element string) bool {
	for _, listElement := range list {
		if listElement == element {
			return true
		}
	}
	return false
}
