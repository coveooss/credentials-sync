package credentials

import (
	"fmt"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"

	"github.com/stretchr/testify/assert"
)

const (
	prefix           = "the_prefix"
	firstSecretName  = "the_prefix/id1"
	firstSecretARN   = "arn:aws:secretsmanager:us-east-1:123456789012:secret:the_prefix/id1-123456"
	secondSecretName = "the_prefix/id2"
	secondSecretARN  = "arn:aws:secretsmanager:us-east-1:123456789012:secret:the_prefix/id2-123456"
	thirdSecretName  = "id3"
	thirdSecretARN   = "arn:aws:secretsmanager:us-east-1:123456789012:secret:id3-123456"
)

var expectedSecretsManagerCredentials = func() []Credentials {
	expectedCred := NewUsernamePassword()
	expectedCred.ID = "test3"
	expectedCred.Description = "a credential"
	expectedCred.Username = "user"
	expectedCred.Password = "pass"
	return append(testCredentials, expectedCred)
}()

func TestCreateSecretsManagerSource(t *testing.T) {
	t.Parallel()

	secretsManagerSource := &AWSSecretsManagerSource{}
	assert.NotNil(t, secretsManagerSource.getClient())
}

func TestSecretsManagerSourceValidate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		source        *AWSSecretsManagerSource
		expectedError error
	}{
		{
			name:          "No config",
			source:        &AWSSecretsManagerSource{},
			expectedError: fmt.Errorf("Either `secret_id` or `secret_prefix` must be defined on a secretsmanager source"),
		},
		{
			name:          "Both secret ID and prefix",
			source:        &AWSSecretsManagerSource{SecretID: "test", SecretPrefix: "test2"},
			expectedError: fmt.Errorf("Either `secret_id` or `secret_prefix` must be defined on a secretsmanager source"),
		},
		{
			name:          "Valid with only ID",
			source:        &AWSSecretsManagerSource{SecretID: "test"},
			expectedError: nil,
		},
		{
			name:          "Valid with only prefix",
			source:        &AWSSecretsManagerSource{SecretPrefix: "test"},
			expectedError: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, "Amazon SecretsManager", tt.source.Type())
			assert.Equal(t, tt.expectedError, tt.source.ValidateConfiguration())
		})
	}
}

type mockSecretsManagerClient struct {
	secretsmanageriface.SecretsManagerAPI
	t *testing.T
}

func (m *mockSecretsManagerClient) ListSecretsPages(input *secretsmanager.ListSecretsInput, theFunc func(*secretsmanager.ListSecretsOutput, bool) bool) error {
	theFunc(&secretsmanager.ListSecretsOutput{SecretList: []*secretsmanager.SecretListEntry{
		{
			ARN:  aws.String(firstSecretARN),
			Name: aws.String(firstSecretName),
		},
		{
			ARN:  aws.String(secondSecretARN),
			Name: aws.String(secondSecretName),
		},
		{
			ARN:  aws.String(thirdSecretARN),
			Name: aws.String(thirdSecretName),
		},
	}}, true)
	return nil
}

func (m *mockSecretsManagerClient) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	if *input.SecretId == firstSecretARN {
		return &secretsmanager.GetSecretValueOutput{SecretString: aws.String(testCredentialsAsMap)}, nil
	} else if *input.SecretId == secondSecretARN {
		return &secretsmanager.GetSecretValueOutput{SecretString: aws.String(`test3:
  type: usernamepassword
  description: a credential
  username: user
  password: pass`)}, nil
	}
	return nil, fmt.Errorf("Only first and second are valid")
}

func TestGetCredentialsFromSecretsManagerSourceWithPrefix(t *testing.T) {
	t.Parallel()

	secretsManagerSource := &AWSSecretsManagerSource{
		SecretPrefix: prefix,
		client:       &mockSecretsManagerClient{t: t},
	}

	credentials, err := secretsManagerSource.Credentials()
	sort.Slice(credentials, func(i, j int) bool { return credentials[i].GetID() < credentials[j].GetID() })

	assert.Nil(t, err)
	assert.Equal(t, expectedSecretsManagerCredentials, credentials)
}

func TestGetCredentialsFromSecretsManagerSourceWithID(t *testing.T) {
	t.Parallel()

	secretsManagerSource := &AWSSecretsManagerSource{
		SecretID: firstSecretARN,
		client:   &mockSecretsManagerClient{t: t},
	}

	credentials, err := secretsManagerSource.Credentials()
	sort.Slice(credentials, func(i, j int) bool { return credentials[i].GetID() < credentials[j].GetID() })

	assert.Nil(t, err)
	assert.Equal(t, testCredentials, credentials)
}

func TestGetCredentialsFromSecretsManagerSourceWithUnknownPrefix(t *testing.T) {
	t.Parallel()

	// Third credentials crashes the GetSecretValue call
	secretsManagerSource := &AWSSecretsManagerSource{
		SecretPrefix: "bad",
		client:       &mockSecretsManagerClient{t: t},
	}

	credentials, err := secretsManagerSource.Credentials()
	assert.EqualError(t, err, "No secrets found with the 'bad' prefix")
	assert.Nil(t, credentials)
}

func TestGetCredentialsFromSecretsManagerSourceWithBadPrefix(t *testing.T) {
	t.Parallel()

	// Third credentials crashes the GetSecretValue call
	secretsManagerSource := &AWSSecretsManagerSource{
		SecretPrefix: thirdSecretName,
		client:       &mockSecretsManagerClient{t: t},
	}

	credentials, err := secretsManagerSource.Credentials()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error while fetching secret "+thirdSecretARN)
	assert.Nil(t, credentials)
}
