package credentials

import (
	"os"
	"path"
	"sort"
	"testing"

	"github.com/mitchellh/mapstructure"

	"github.com/stretchr/testify/assert"
)

func TestSourcesConfigWithLocalSource(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)
	filePath := path.Join(tempDir, "local_file.json")
	os.WriteFile(filePath, []byte(`[{"id": "test", "type": "secret", "description": "test-desc", "secret": "my secret"}]`), 0777)
	localSource := &LocalSource{
		File: filePath,
	}

	sourcesConfig := SourcesConfiguration{LocalSources: []*LocalSource{localSource}}

	credentials, err := sourcesConfig.Credentials()
	assert.Nil(t, err)
	assert.Equal(t, []Credentials{testCredentials[0]}, credentials)
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
			name:    "List",
			bytes:   []byte(testCredentialsAsList),
			result:  testCredentials,
			wantErr: false,
		},
		{
			name:    "Map",
			bytes:   []byte(testCredentialsAsMap),
			result:  testCredentials,
			wantErr: false,
		},
		{
			name: "Single cred",
			bytes: []byte(`{
				"id": "test",
				"type": "secret",
				"description": "test-desc",
				"secret": "my secret"
			}`),
			result:  testCredentials[0:1],
			wantErr: false,
		},
		{
			name:    "SSH Creds",
			bytes:   []byte(testSSHCredentialsString),
			result:  []Credentials{testSSHCredentials},
			wantErr: false,
		},
		{
			name: "Cred without ID",
			bytes: []byte(`{
				"type": "secret",
				"description": "test-desc",
				"secret": "my secret"
			}`),
			result:  nil,
			wantErr: true,
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
