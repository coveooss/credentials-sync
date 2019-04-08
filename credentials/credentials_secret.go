package credentials

import (
	"fmt"
)

// SecretTextCredentials represents credentials composed of a single string value
type SecretTextCredentials struct {
	Base   `mapstructure:",squash"`
	Secret string
}

// NewSecretText instantiates a SecretTextCredentials struct
func NewSecretText() *SecretTextCredentials {
	cred := &SecretTextCredentials{}
	cred.CredType = "Secret text"
	return cred
}

// ToString prints out the content of a SecretTextCredentials struct.
// If showSensitive is true, the secret text will be shown
func (cred *SecretTextCredentials) ToString(showSensitive bool) string {
	secretText := "*******"
	if showSensitive {
		secretText = cred.Secret
	}
	if cred.Secret == "" {
		secretText = "<empty>"
	}
	return fmt.Sprintf("%s - %s", cred.BaseToString(), secretText)
}

// Validate verifies that the credentials is valid.
// A SecretTextCredentials is always considered valid, as empty values are accepted.
func (cred *SecretTextCredentials) Validate() bool {
	if cred.Secret == "" && cred.Value != "" {
		cred.Secret = cred.Value
	}
	return true
}
