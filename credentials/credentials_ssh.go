package credentials

import (
	"fmt"
)

// SSHCredentials represents credentials composed of a private key, username and passphrase
type SSHCredentials struct {
	Base       `mapstructure:",squash"`
	Username   string
	Passphrase string
	PrivateKey string `mapstructure:"private_key"`
}

// NewSSHCredentials instantiates a SSHCredentials struct
func NewSSHCredentials() *SSHCredentials {
	cred := &SSHCredentials{}
	cred.CredType = "SSH"
	return cred
}

// ToString prints out the content of a UsernamePasswordCredentials struct.
// If showSensitive is true, the passphrase will be shown
func (cred *SSHCredentials) ToString(showSensitive bool) string {
	passphrase := "********"
	username := cred.Username
	if showSensitive {
		passphrase = cred.Passphrase
	}
	if cred.Username == "" {
		username = "<empty>"
	}
	if passphrase == "" {
		passphrase = "<empty>"
	}
	return fmt.Sprintf("%s - %s:%s", cred.BaseToString(), username, passphrase)
}

// Validate verifies that the credentials is valid.
// A SSHCredentials must have a private key, the username and passphrase are optional
func (cred *SSHCredentials) Validate() error {
	if cred.PrivateKey == "" {
		return fmt.Errorf("The credentials with ID %s does not define a private key", cred.ID)
	}
	return nil
}
