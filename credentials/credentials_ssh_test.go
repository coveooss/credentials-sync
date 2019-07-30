package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSHCredentialsToString(t *testing.T) {
	cred := NewSSHCredentials()
	cred.ID = "test"
	cred.Username = "key"
	cred.Passphrase = "secret"
	assert.Equal(t, "test -> Type: SSH - key:********", cred.ToString(false))
	assert.Equal(t, "test -> Type: SSH - key:secret", cred.ToString(true))

	// Empty creds
	cred = NewSSHCredentials()
	cred.ID = "test"
	assert.Equal(t, "test -> Type: SSH - <empty>:********", cred.ToString(false))
	assert.Equal(t, "test -> Type: SSH - <empty>:<empty>", cred.ToString(true))
}

func TestSSHCredentialsValidation(t *testing.T) {
	cred := NewSSHCredentials()
	cred.ID = "test"
	cred.Username = "key"
	cred.Passphrase = "secret"
	assert.Error(t, cred.Validate())

	cred.PrivateKey = "private"
	assert.Nil(t, cred.Validate())
}
