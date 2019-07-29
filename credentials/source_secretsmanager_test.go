package credentials

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
