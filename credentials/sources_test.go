package credentials

import (
	"sort"
	"testing"

	"github.com/mitchellh/mapstructure"

	"github.com/stretchr/testify/assert"
)

var testCredentials = []Credentials{
	func() (creds *SecretTextCredentials) {
		creds = NewSecretText()
		creds.ID = "test"
		creds.Secret = "my secret"
		creds.Description = "test-desc"
		return
	}(),
	func() (creds *UsernamePasswordCredentials) {
		creds = NewUsernamePassword()
		creds.ID = "test2"
		creds.Username = "my"
		creds.Password = "secret"
		creds.Description = "test2-desc"
		return
	}(),
}

func TestGetCredentialsFromBytes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		bytes   []byte
		result  []Credentials
		wantErr bool
	}{
		{
			name:    "Empty file",
			bytes:   []byte{},
			result:  []Credentials{},
			wantErr: false,
		},
		{
			name: "List",
			bytes: []byte(`[
				{
					"id": "test",
					"type": "secret",
					"description": "test-desc",
					"secret": "my secret"
				},
				{
					"id": "test2",
					"type": "usernamepassword",
					"description": "test2-desc",
					"username": "my",
					"password": "secret"
				}
			]`),
			result:  testCredentials,
			wantErr: false,
		},
		{
			name: "Map",
			bytes: []byte(`{
				"test": {
					"type": "secret",
					"description": "test-desc",
					"secret": "my secret"
				},
				"test2": {
					"type": "usernamepassword",
					"description": "test2-desc",
					"username": "my",
					"password": "secret"
				}
			}`),
			result:  testCredentials,
			wantErr: false,
		},
		{
			name: "Bad map",
			bytes: []byte(`{
				"test": "sg",
				"test2": "fsef"
			}`),
			result:  nil,
			wantErr: true,
		},
		{
			name: "Bad list",
			bytes: []byte(`[
				"test",
				"test2"
			]`),
			result:  nil,
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getCredentialsFromBytes(tt.bytes)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			var expectedAsMaps []map[string]interface{}
			var gottenAsMaps []map[string]interface{}
			mapstructure.Decode(tt.result, &expectedAsMaps)
			mapstructure.Decode(result, &gottenAsMaps)
			sort.Slice(gottenAsMaps, func(i int, j int) bool {
				return gottenAsMaps[i]["ID"].(string) < gottenAsMaps[j]["ID"].(string)
			})
			assert.Equal(t, expectedAsMaps, gottenAsMaps)
		})
	}
}
