package credentials

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCredentialsFromLocalSource(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	filePath := path.Join(tempDir, "local_file.json")
	localSource := &LocalSource{
		File: filePath,
	}
	assert.Equal(t, "Local file", localSource.Type())
	assert.Error(t, localSource.ValidateConfiguration())

	os.WriteFile(filePath, []byte(`test_cred:
  type: usernamepassword
  description: a credential
  username: user
  password: pass`), 0777)

	assert.Nil(t, localSource.ValidateConfiguration())
	credentials, err := localSource.Credentials()
	expectedCred := NewUsernamePassword()
	expectedCred.ID = "test_cred"
	expectedCred.Description = "a credential"
	expectedCred.Username = "user"
	expectedCred.Password = "pass"
	assert.Nil(t, err)
	assert.Equal(t, []Credentials{expectedCred}, credentials)
}
