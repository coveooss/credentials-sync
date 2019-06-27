package credentials

import (
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

// AWSS3Source represents s3 objects containing credentials
type AWSS3Source struct {
	Bucket string
	Key    string
}

// Credentials extracts credentials from the source
func (source *AWSS3Source) Credentials() ([]Credentials, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client := s3.New(sess)

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
func (source *AWSS3Source) ValidateConfiguration() bool {
	if source.Bucket == "" {
		log.Errorf("S3 sources must define a bucket")
		return false
	}
	if source.Key == "" {
		log.Errorf("S3 sources must define a key")
		return false
	}
	return true
}
