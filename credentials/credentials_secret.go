package credentials

import (
	"fmt"
)

type SecretTextCredentials struct {
	Base   `mapstructure:",squash"`
	Secret string
}

func NewSecretText() *SecretTextCredentials {
	cred := &SecretTextCredentials{}
	cred.CredType = "Secret text"
	return cred
}

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

func (cred *SecretTextCredentials) Validate() bool {
	return true
}
