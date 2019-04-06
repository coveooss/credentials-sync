package credentials

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type Credentials interface {
	ToString() string
}

type CredentialsBase struct {
	ID          string
	Description string
	CredType    string
}

func (credBase *CredentialsBase) BaseToString() string {
	return fmt.Sprintf("Type: %s, ID: %s, Description: %s", credBase.CredType, credBase.ID, credBase.Description)
}

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
	default:
		return nil, errors.New("Unknown credentials type")
	}
	mapstructure.Decode(credentialsMap, credentials)
	return credentials, nil
}
