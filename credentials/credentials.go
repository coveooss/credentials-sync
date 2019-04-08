package credentials

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

// Credentials defines the methods that can be called by all types of credentials
type Credentials interface {
	BaseValidate() bool
	GetID() string
	ToString(bool) string
	Validate() bool
}

// Base defines that fields that are common to all types of credentials
type Base struct {
	ID          string
	Description string
	CredType    string

	// For multi-value fields. Such as SSM
	Value string
}

// BaseToString prints out the credentials fields common to all types of credentials
func (credBase *Base) BaseToString() string {
	return fmt.Sprintf("%s -> Type: %s, Description: %s", credBase.ID, credBase.CredType, credBase.Description)
}

// BaseValidate verifies that the credentials fields common to all types of credentials contain valid values
func (credBase *Base) BaseValidate() bool {
	if credBase.ID == "" {
		log.Errorf("Credentials (%s) has no defined ID", credBase.BaseToString())
	}
	if credBase.Description == "" {
		log.Errorf("Credentials (%s) has no defined description", credBase.ID)
	}
	if credBase.CredType == "" {
		log.Errorf("Credentials (%s) has no type. This is probably a bug in the software", credBase.ID)
	}
	return credBase.ID != "" && credBase.Description != "" && credBase.CredType != ""
}

// GetID returns a credentials' ID
func (credBase *Base) GetID() string {
	return credBase.ID
}

// ParseCredentials transforms a list of maps into a list of Credentials
// The credentials type is determined by the `type` attribute
func ParseCredentials(credentialsMaps []map[string]interface{}) ([]Credentials, error) {
	credentialsList := []Credentials{}
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
	if value, ok := credentialsMap["type"]; ok {
		if credentialsType, ok = value.(string); !ok {
			return nil, errors.New("Credentials type is not a string")
		}
	} else {
		return nil, errors.New("Unable to find the credentials type")
	}

	switch credentialsType {
	case "usernamepassword":
		credentials = NewUsernamePassword()
	case "secret":
		credentials = NewSecretText()
	default:
		return nil, errors.New("Unknown credentials type")
	}
	mapstructure.Decode(credentialsMap, credentials)
	if !credentials.BaseValidate() || !credentials.Validate() {
		return nil, errors.New("The following credentials failed to validate: \n	" + credentials.ToString(false))
	}
	return credentials, nil
}
