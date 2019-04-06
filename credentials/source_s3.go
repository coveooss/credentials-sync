package credentials

type AWSS3Source struct{}

func (source *AWSS3Source) Credentials() ([]Credentials, error) {
	return []Credentials{}, nil
}

func (source *AWSS3Source) Type() string {
	return "Amazon S3"
}

func (source *AWSS3Source) ValidateConfiguration() bool {
	return true
}
