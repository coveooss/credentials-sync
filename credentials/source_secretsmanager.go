package credentials

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

// AWSSecretsManagerSource represents AWS SecretsManager secrets containing credentials
type AWSSecretsManagerSource struct {
	SecretPrefix string `mapstructure:"secret_prefix"`
	SecretID     string `mapstructure:"secret_id"`

	client secretsmanageriface.SecretsManagerAPI
}

func (source *AWSSecretsManagerSource) getClient() secretsmanageriface.SecretsManagerAPI {
	if source.client == nil {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState:       session.SharedConfigEnable,
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		}))
		source.client = secretsmanager.New(sess)
	}
	return source.client
}

// Credentials extracts credentials from the source
func (source *AWSSecretsManagerSource) Credentials() ([]Credentials, error) {
	client := source.getClient()

	secretIDs := []string{}

	if source.SecretPrefix != "" {
		if err := client.ListSecretsPages(&secretsmanager.ListSecretsInput{}, func(output *secretsmanager.ListSecretsOutput, lastPage bool) bool {
			for _, secret := range output.SecretList {
				if strings.HasPrefix(*secret.Name, source.SecretPrefix) {
					secretIDs = append(secretIDs, *secret.ARN)
				}
			}
			return !lastPage
		}); err != nil {
			return nil, fmt.Errorf("Error listing secrets: %v", err)
		}

		if len(secretIDs) == 0 {
			return nil, fmt.Errorf("No secrets found with the '%s' prefix", source.SecretPrefix)
		}
	} else if source.SecretID != "" {
		secretIDs = append(secretIDs, source.SecretID)
	}

	credentials := []Credentials{}
	for _, secretID := range secretIDs {
		value, err := client.GetSecretValue(&secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretID),
		})
		if err != nil {
			return nil, fmt.Errorf("Error while fetching secret %s: %v", secretID, err)
		}
		fetchedCredentials, err := getCredentialsFromBytes([]byte(*value.SecretString))
		if err != nil {
			return nil, fmt.Errorf("Error while parsing credentials from secret %s: %v", secretID, err)
		}
		credentials = append(credentials, fetchedCredentials...)
	}

	return credentials, nil
}

// Type returns the type of the source
func (source *AWSSecretsManagerSource) Type() string {
	return "Amazon SecretsManager"
}

// ValidateConfiguration verifies that the source's attributes are valid
func (source *AWSSecretsManagerSource) ValidateConfiguration() error {
	if (source.SecretID == "" && source.SecretPrefix == "") || (source.SecretID != "" && source.SecretPrefix != "") {
		return fmt.Errorf("Either `secret_id` or `secret_prefix` must be defined on a secretsmanager source")
	}
	return nil
}
