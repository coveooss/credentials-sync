package credentials

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/stretchr/testify/assert"
)

const (
	s3Bucket = "my-bucket"
	s3Key    = "a/key"
)


func TestCreateS3Source(t *testing.T) {
	s3Source := &AWSS3Source{}
	assert.NotNil(t, s3Source.getClient())
}

func TestS3SourceValidate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		source        *AWSS3Source
		expectedError error
	}{
		{
			name:          "No bucket",
			source:        &AWSS3Source{Key: "test"},
			expectedError: fmt.Errorf("S3 sources must define a bucket"),
		},
		{
			name:          "No key",
			source:        &AWSS3Source{Bucket: "bucket"},
			expectedError: fmt.Errorf("S3 sources must define a key"),
		},
		{
			name:          "Valid",
			source:        &AWSS3Source{Bucket: "bucket", Key: "test"},
			expectedError: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, "Amazon S3", tt.source.Type())
			assert.Equal(t, tt.expectedError, tt.source.ValidateConfiguration())
		})
	}
}

type mockS3Client struct {
	s3iface.S3API
	t *testing.T
}

func (m *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	assert.Equal(m.t, s3Bucket, *input.Bucket)
	assert.Equal(m.t, s3Key, *input.Key)

	return &s3.GetObjectOutput{Body: ioutil.NopCloser(strings.NewReader(`test_cred:
  type: usernamepassword
  description: a credential
  username: user
  password: pass`))}, nil
}

func TestGetCredentialsFromS3Source(t *testing.T) {
	s3Source := &AWSS3Source{
		Bucket: s3Bucket,
		Key:    s3Key,
		client: &mockS3Client{t: t},
	}

	credentials, err := s3Source.Credentials()

	expectedCred := NewUsernamePassword()
	expectedCred.ID = "test_cred"
	expectedCred.Description = "a credential"
	expectedCred.Username = "user"
	expectedCred.Password = "pass"
	assert.Nil(t, err)
	assert.Equal(t, []Credentials{expectedCred}, credentials)
}
