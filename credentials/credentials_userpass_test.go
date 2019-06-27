package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUserPassCredentialsFromValue(t *testing.T) {
	t.Parallel()

	credInterface, err := ParseSingleCredentials(map[string]interface{}{
		"id":    "test",
		"type":  "usernamepassword",
		"value": "user:pass",
	})

	assert.Nil(t, err)
	cred := credInterface.(*UsernamePasswordCredentials)

	assert.Equal(t, "user", cred.Username)
	assert.Equal(t, "pass", cred.Password)
}

func TestUserPassCredentialsToString(t *testing.T) {
	cred := NewUsernamePassword()
	cred.ID = "test"
	cred.Username = "key"
	cred.Password = "secret"
	assert.Equal(t, "test -> Type: Username/Password - key:********", cred.ToString(false))
	assert.Equal(t, "test -> Type: Username/Password - key:secret", cred.ToString(true))

	// Empty creds
	cred = NewUsernamePassword()
	cred.ID = "test"
	assert.Equal(t, "test -> Type: Username/Password - <empty>:********", cred.ToString(false))
	assert.Equal(t, "test -> Type: Username/Password - <empty>:<empty>", cred.ToString(true))
}
