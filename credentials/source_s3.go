package credentials

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// AWSS3Source represents s3 objects containing credentials
type AWSS3Source struct {
	Bucket string
	Key    string

	client s3iface.S3API
}

func (source *AWSS3Source) getClient() s3iface.S3API {
	if source.client == nil {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState:       session.SharedConfigEnable,
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		}))
		source.client = s3.New(sess)
	}
	return source.client
}

// Credentials extracts credentials from the source
func (source *AWSS3Source) Credentials() ([]Credentials, error) {
	client := source.getClient()

	response, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(source.Bucket),
		Key:    aws.String(source.Key),
	})
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return getCredentialsFromBytes(body)
}

// Type returns the type of the source
func (source *AWSS3Source) Type() string {
	return "Amazon S3"
}

// ValidateConfiguration verifies that the source's attributes are valid
func (source *AWSS3Source) ValidateConfiguration() error {
	if source.Bucket == "" {
		return fmt.Errorf("S3 sources must define a bucket")
	}
	if source.Key == "" {
		return fmt.Errorf("S3 sources must define a key")
	}
	return nil
}
