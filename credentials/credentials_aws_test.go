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

func TestParseAwsCredentialsFromInvalidValue(t *testing.T) {
	t.Parallel()

	_, err := ParseSingleCredentials(map[string]interface{}{
		"id":    "test",
		"type":  "aws",
		"value": "key",
	})

	assert.Error(t, err)
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

	// Extra data
	credMap["extra"] = "data"
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

func TestCredentialWithTargetTags(t *testing.T) {
	credMap := map[string]interface{}{
		"id":          "test",
		"type":        "aws",
		"description": "test-desc",
		"access_key":  "key",
		"secret_key":  "secret_key",
		"target_tags": map[string]interface{}{
			"do_match": map[string]interface{}{
				"tag1": "value1",
			},
		},
	}

	cred, err := ParseSingleCredentials(credMap)
	assert.Nil(t, err)

	assert.Equal(t, &AmazonWebServicesCredentials{
		Base: Base{
			ID:          "test",
			CredType:    "Amazon Web Services",
			Description: "test-desc",
			TargetTags:  targetTagsMatcher{DoMatch: map[string]interface{}{"tag1": "value1"}},
		},
		AccessKey: "key",
		SecretKey: "secret_key",
	}, cred)
}

func TestCredentialWithTargetTagsMalformed(t *testing.T) {
	credMap := map[string]interface{}{
		"id":         "test",
		"type":       "aws",
		"access_key": "key",
		"secret_key": "secret_key",
		"target_tags": map[string]interface{}{
			"match": map[string]interface{}{
				"tag1": "value1",
			},
		},
	}

	_, err := ParseSingleCredentials(credMap)
	assert.Error(t, err)
}
