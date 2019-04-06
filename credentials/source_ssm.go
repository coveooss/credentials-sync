package credentials

type AWSSSMSource struct{}

func (source *AWSSSMSource) Credentials() ([]Credentials, error) {
	return []Credentials{}, nil
}

func (source *AWSSSMSource) Type() string {
	return "Amazon SSM"
}

func (source *AWSSSMSource) ValidateConfiguration() bool {
	return true
}
