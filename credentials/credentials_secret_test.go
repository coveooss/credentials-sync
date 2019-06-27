package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsSecretCredentialsFromValue(t *testing.T) {
	t.Parallel()

	credInterface, err := ParseSingleCredentials(map[string]interface{}{
		"id":    "test",
		"type":  "secret",
		"value": "my secret",
	})

	assert.Nil(t, err)
	cred := credInterface.(*SecretTextCredentials)

	assert.Equal(t, "my secret", cred.Secret)
}

func TestSecretCredentialsToString(t *testing.T) {
	cred := NewSecretText()
	cred.ID = "test"
	cred.Secret = "secret"
	assert.Equal(t, "test -> Type: Secret text - ********", cred.ToString(false))
	assert.Equal(t, "test -> Type: Secret text - secret", cred.ToString(true))

	// Empty creds
	cred = NewSecretText()
	cred.ID = "test"
	assert.Equal(t, "test -> Type: Secret text - ********", cred.ToString(false))
	assert.Equal(t, "test -> Type: Secret text - <empty>", cred.ToString(true))
}
