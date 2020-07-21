package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubAppCredentials(t *testing.T) {
	tests := map[string]struct {
		givenID            string
		givenAppID         int
		givenOwner         string
		givenShowSensitive bool
		expectString       string
	}{
		"without owner": {givenID: "test", givenAppID: 1, givenOwner: "", givenShowSensitive: false, expectString: "test -> Type: Github App - 1"},
		"with owner":    {givenID: "test", givenAppID: 2, givenOwner: "owner", givenShowSensitive: false, expectString: "test -> Type: Github App - 2(owner)"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cred := NewGithubAppCredentials()
			cred.ID = test.givenID
			cred.AppID = test.givenAppID
			cred.Owner = test.givenOwner
			assert.Equal(t, test.expectString, cred.ToString(test.givenShowSensitive))
		})
	}
}

func TestGithubAppCredentialsValidation(t *testing.T) {
	tests := map[string]struct {
		givenCred   GithubAppCredentials
		expectError bool
	}{
		"valid": {givenCred: GithubAppCredentials{
			AppID:      12345,
			PrivateKey: "private",
			Owner:      "Me",
		}, expectError: false},
		"valid no owner": {givenCred: GithubAppCredentials{
			AppID:      12345,
			PrivateKey: "private",
		}, expectError: false},
		"missing app id": {givenCred: GithubAppCredentials{
			PrivateKey: "private",
		}, expectError: true},
		"missing private key": {givenCred: GithubAppCredentials{
			AppID: 12345,
		}, expectError: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cred := NewGithubAppCredentials()
			cred.AppID = test.givenCred.AppID
			cred.PrivateKey = test.givenCred.PrivateKey
			cred.Owner = test.givenCred.Owner
			err := cred.Validate()
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
