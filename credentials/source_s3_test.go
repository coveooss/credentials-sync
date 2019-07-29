package credentials

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
