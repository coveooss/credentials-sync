package targets

import (
	"testing"

	"github.com/bndr/gojenkins"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/stretchr/testify/assert"
)

func TestJenkinsToString(t *testing.T) {
	jenkins := &JenkinsTarget{
		Base: Base{Name: "targetName", Tags: map[string]string{"my_tag": "tag_value"}},
		URL:  "test.com",
	}

	assert.Equal(t, "targetName [Tags: my_tag=tag_value] (Jenkins) - test.com", jenkins.ToString())
}

func TestAwsCredentialsToJenkinsCred(t *testing.T) {
	secret := credentials.NewAmazonWebServicesCredentials()
	secret.ID = "test-id"
	secret.Description = "a test description"
	secret.AccessKey = "AKIASOMETHING"
	secret.SecretKey = "test-secret"
	secret.RoleARN = "my-role"

	jenkinsSecretInterface := toJenkinsCredential(secret)

	jenkinsSecret := jenkinsSecretInterface.(*gojenkins.AmazonWebServicesCredentials)
	assert.Equal(t, secret.ID, jenkinsSecret.ID)
	assert.Equal(t, secret.Description, jenkinsSecret.Description)
	assert.Equal(t, secret.AccessKey, jenkinsSecret.AccessKey)
	assert.Equal(t, secret.SecretKey, jenkinsSecret.SecretKey)
	assert.Equal(t, secret.RoleARN, jenkinsSecret.IAMRoleARN)
}

func TestSecretToJenkinsCred(t *testing.T) {
	secret := credentials.NewSecretText()
	secret.ID = "test-id"
	secret.Description = "a test description"
	secret.Secret = "a-secret"

	jenkinsSecretInterface := toJenkinsCredential(secret)

	jenkinsSecret := jenkinsSecretInterface.(*gojenkins.StringCredentials)
	assert.Equal(t, secret.ID, jenkinsSecret.ID)
	assert.Equal(t, secret.Description, jenkinsSecret.Description)
	assert.Equal(t, secret.Secret, jenkinsSecret.Secret)
}

func TestSSHKeyToJenkinsCred(t *testing.T) {
	secret := credentials.NewSSHCredentials()
	secret.ID = "test-id"
	secret.Username = "a-user"
	secret.Passphrase = "a-password"
	secret.PrivateKey = "a-key"

	jenkinsSecretInterface := toJenkinsCredential(secret)

	jenkinsSecret := jenkinsSecretInterface.(*gojenkins.SSHCredentials)
	assert.Equal(t, secret.ID, jenkinsSecret.ID)
	assert.Equal(t, secret.ID, jenkinsSecret.Description) // No description, then use ID as description
	assert.Equal(t, secret.Username, jenkinsSecret.Username)
	assert.Equal(t, secret.Passphrase, jenkinsSecret.Passphrase)
	privateKeySource := jenkinsSecret.PrivateKeySource.(*gojenkins.PrivateKey)
	assert.Equal(t, secret.PrivateKey, privateKeySource.Value)
}

func TestUsernamePasswordToJenkinsCred(t *testing.T) {
	secret := credentials.NewUsernamePassword()
	secret.ID = "test-id"
	secret.Description = "a test description"
	secret.Username = "a-user"
	secret.Password = "a-password"

	jenkinsSecretInterface := toJenkinsCredential(secret)

	jenkinsSecret := jenkinsSecretInterface.(*gojenkins.UsernameCredentials)
	assert.Equal(t, secret.ID, jenkinsSecret.ID)
	assert.Equal(t, secret.Description, jenkinsSecret.Description)
	assert.Equal(t, secret.Username, jenkinsSecret.Username)
	assert.Equal(t, secret.Password, jenkinsSecret.Password)
}
