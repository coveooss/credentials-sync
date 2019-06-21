package credentials

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type AWSSecretsManagerSource struct {
	SecretID string `mapstructure:"secret_id"`
}

func (source *AWSSecretsManagerSource) Credentials() ([]Credentials, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}))
	client := secretsmanager.New(sess)
	value, err := client.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(source.SecretID),
	})
	if err != nil {
		return nil, err
	}
	return getCredentialsFromBytes([]byte(*value.SecretString))
}

func (source *AWSSecretsManagerSource) Type() string {
	return "Amazon SecretsManager"
}

func (source *AWSSecretsManagerSource) ValidateConfiguration() bool {
	return len(source.SecretID) > 0
}
