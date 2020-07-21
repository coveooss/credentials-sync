package credentials

import (
	"fmt"
)

// GithubAppCredentials represents credentials composed of an App ID, private key, and owner
type GithubAppCredentials struct {
	Base       `mapstructure:",squash"`
	AppID      int    `mapstructure:"app_id"`
	PrivateKey string `mapstructure:"private_key"`
	Owner      string `mapstructure:"owner"`
}

// NewGithubAppCredentials instantiates a GithubAppCredentials struct
func NewGithubAppCredentials() *GithubAppCredentials {
	cred := &GithubAppCredentials{}
	cred.CredType = "Github App"
	return cred
}

// ToString prints out the content of a GithubAppCredentials struct.
// showSensitive bool has not effect. It is provided to satisfy the interface.
func (cred *GithubAppCredentials) ToString(_ bool) string {
	appIDOwner := fmt.Sprintf("%d", cred.AppID)
	if len(cred.Owner) > 0 {
		appIDOwner = fmt.Sprintf("%s(%s)", appIDOwner, cred.Owner)
	}
	return fmt.Sprintf("%s - %s", cred.BaseToString(), appIDOwner)
}

// Validate verifies that the credentials is valid.
// A GithubAppCredentials must have an app id and a private key.  Owner is optional.
func (cred *GithubAppCredentials) Validate() error {
	switch {
	case cred.AppID == 0:
		return fmt.Errorf("the credentials with ID %s does not define an app ID", cred.ID)
	case len(cred.PrivateKey) == 0:
		return fmt.Errorf("the credentials with ID %s does not define a private key", cred.ID)
	default:
		return nil
	}
}
