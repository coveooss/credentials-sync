package credentials

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type AWSSecretsManagerSource struct {
	SecretPrefix string `mapstructure:"secret_prefix"`
	SecretID     string `mapstructure:"secret_id"`
}

func (source *AWSSecretsManagerSource) Credentials() ([]Credentials, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}))
	client := secretsmanager.New(sess)

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
	} else {
		return nil, fmt.Errorf("Either `secret_id` or `secret_prefix` must be defined")
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

func (source *AWSSecretsManagerSource) Type() string {
	return "Amazon SecretsManager"
}

func (source *AWSSecretsManagerSource) ValidateConfiguration() bool {
	return len(source.SecretID) > 0
}
