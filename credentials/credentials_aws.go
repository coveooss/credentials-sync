package credentials

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// SecretTextCredentials represents credentials composed of a single string value
type AmazonWebServicesCredentials struct {
	Base            `mapstructure:",squash"`
	AccessKey       string `mapstructure:"access_key"`
	SecretKey       string `mapstructure:"secret_key"`
	RoleARN         string `mapstructure:"role_arn"`
	MFASerialNumber string `mapstructure:"mfa_serial"`
}

// NewSecretText instantiates a SecretTextCredentials struct
func NewAmazonWebServicesCredentials() *AmazonWebServicesCredentials {
	cred := &AmazonWebServicesCredentials{}
	cred.CredType = "Amazon Web Services"
	return cred
}

// ToString prints out the content of a SecretTextCredentials struct.
// If showSensitive is true, the secret text will be shown
func (cred *AmazonWebServicesCredentials) ToString(showSensitive bool) string {
	secretKey := "********"
	if showSensitive {
		secretKey = cred.SecretKey
	}
	return fmt.Sprintf("%s - %s:%s", cred.BaseToString(), cred.AccessKey, secretKey)
}

// Validate verifies that the credentials is valid.
// A SecretTextCredentials is always considered valid, as empty values are accepted.
func (cred *AmazonWebServicesCredentials) Validate() bool {
	if cred.AccessKey == "" {
		log.Errorf("The credentials with ID %s does not define an access key", cred.ID)
		return false
	}
	if cred.SecretKey == "" {
		log.Errorf("The credentials with ID %s does not define an secret key", cred.ID)
		return false
	}
	return true
}
