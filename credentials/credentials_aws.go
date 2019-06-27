package credentials

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// AmazonWebServicesCredentials represents an access key id and a secret access key from AWS. A role to assume can also be defined
type AmazonWebServicesCredentials struct {
	Base            `mapstructure:",squash"`
	AccessKey       string `mapstructure:"access_key"`
	SecretKey       string `mapstructure:"secret_key"`
	RoleARN         string `mapstructure:"role_arn"`
	MFASerialNumber string `mapstructure:"mfa_serial"`
}

// NewAmazonWebServicesCredentials instantiates an AmazonWebServicesCredentials struct
func NewAmazonWebServicesCredentials() *AmazonWebServicesCredentials {
	cred := &AmazonWebServicesCredentials{}
	cred.CredType = "Amazon Web Services"
	return cred
}

// ToString prints out the content of a AmazonWebServicesCredentials struct.
// If showSensitive is true, the secret access key will be shown
func (cred *AmazonWebServicesCredentials) ToString(showSensitive bool) string {
	secretKey := "********"
	if showSensitive {
		secretKey = cred.SecretKey
	}
	return fmt.Sprintf("%s - %s:%s", cred.BaseToString(), cred.AccessKey, secretKey)
}

// Validate verifies that the credentials is valid.
// A AmazonWebServicesCredentials must define an access key and a secret access key
func (cred *AmazonWebServicesCredentials) Validate() bool {
	if cred.AccessKey == "" && cred.SecretKey == "" && cred.Value != "" {
		splitValue := strings.Split(cred.Value, ":")
		if len(splitValue) != 2 {
			log.Errorf("The credentials with ID %s has an invalid access_key:secret_key value: %s", cred.ID, cred.Value)
			return false
		}
		cred.AccessKey = splitValue[0]
		cred.SecretKey = splitValue[1]
	}

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
