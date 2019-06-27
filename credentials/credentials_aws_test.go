package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAwsCredentialsFromValue(t *testing.T) {
	t.Parallel()

	credInterface, err := ParseSingleCredentials(map[string]interface{}{
		"id":    "test",
		"type":  "aws",
		"value": "key:secret",
	})

	assert.Nil(t, err)
	cred := credInterface.(*AmazonWebServicesCredentials)

	assert.Equal(t, "key", cred.AccessKey)
	assert.Equal(t, "secret", cred.SecretKey)
}

func TestAwsCredentialsValidationErrors(t *testing.T) {
	t.Parallel()

	credMap := map[string]interface{}{
		"id":         "test",
		"type":       "aws",
		"access_key": "key",
	}

	// No secret key
	_, err := ParseSingleCredentials(credMap)
	assert.Error(t, err)

	// All OK
	credMap["secret_key"] = "secret"
	_, err = ParseSingleCredentials(credMap)
	assert.Nil(t, err)

	// No access key
	delete(credMap, "access_key")
	_, err = ParseSingleCredentials(credMap)
	assert.Error(t, err)
}

func TestAwsCredentialsToString(t *testing.T) {
	cred := NewAmazonWebServicesCredentials()
	cred.ID = "test"
	cred.AccessKey = "key"
	cred.SecretKey = "secret"
	assert.Equal(t, "test -> Type: Amazon Web Services - key:********", cred.ToString(false))
	assert.Equal(t, "test -> Type: Amazon Web Services - key:secret", cred.ToString(true))
}
