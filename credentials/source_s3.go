package credentials

import (
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

type AWSS3Source struct {
	Bucket string
	Key    string
}

func (source *AWSS3Source) Credentials() ([]Credentials, error) {
	downloader := s3manager.NewDownloader(session.New())

	file, err := ioutil.TempFile("", "credentials_sync_s3")
	defer os.Remove(file.Name())
	if err != nil {
		return nil, err
	}

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(source.Bucket),
			Key:    aws.String(source.Key),
		})
	if err != nil {
		return nil, err
	}
	log.Info("Downloaded", file.Name(), numBytes, "bytes")

	return getCredentialsFromFile(file.Name())
}

func (source *AWSS3Source) Type() string {
	return "Amazon S3"
}

func (source *AWSS3Source) ValidateConfiguration() bool {
	return true
}
