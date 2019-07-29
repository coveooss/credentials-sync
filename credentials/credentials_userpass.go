package credentials

import (
	"fmt"
	"strings"
)

// UsernamePasswordCredentials represents credentials composed of a username and a password
type UsernamePasswordCredentials struct {
	Base     `mapstructure:",squash"`
	Username string
	Password string
}

// NewUsernamePassword instantiates a UsernamePasswordCredentials struct
func NewUsernamePassword() *UsernamePasswordCredentials {
	cred := &UsernamePasswordCredentials{}
	cred.CredType = "Username/Password"
	return cred
}

// ToString prints out the content of a UsernamePasswordCredentials struct.
// If showSensitive is true, the password will be shown
func (cred *UsernamePasswordCredentials) ToString(showSensitive bool) string {
	password := "********"
	username := cred.Username
	if showSensitive {
		password = cred.Password
	}
	if cred.Username == "" {
		username = "<empty>"
	}
	if password == "" {
		password = "<empty>"
	}
	return fmt.Sprintf("%s - %s:%s", cred.BaseToString(), username, password)
}

// Validate verifies that the credentials is valid.
// A UsernamePasswordCredentials is always considered valid, as empty values are accepted.
func (cred *UsernamePasswordCredentials) Validate() error {
	if cred.Username == "" && cred.Password == "" && cred.Value != "" {
		splitValue := strings.Split(cred.Value, ":")
		if len(splitValue) != 2 {
			return fmt.Errorf("The credentials with ID %s has an invalid username:password value: %s", cred.ID, cred.Value)
		}
		cred.Username = splitValue[0]
		cred.Password = splitValue[1]
	}
	return nil
}
