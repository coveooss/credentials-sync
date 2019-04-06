package credentials

type Source interface {
	Credentials() ([]Credentials, error)
	Type() string
	ValidateConfiguration() bool
}

type SourcesConfiguration struct {
	AWSS3Sources  []*AWSS3Source  `yaml:"aws_s3"`
	AWSSSMSources []*AWSSSMSource `yaml:"aws_ssm"`
	LocalSources  []*LocalSource  `yaml:"local"`
}

func (sc *SourcesConfiguration) AllSources() []Source {
	sources := []Source{}
	for _, source := range sc.LocalSources {
		sources = append(sources, source)
	}
	for _, source := range sc.AWSS3Sources {
		sources = append(sources, source)
	}
	for _, source := range sc.AWSSSMSources {
		sources = append(sources, source)
	}
	return sources
}

func (sc *SourcesConfiguration) ValidateConfiguration() bool {
	configIsOK := true
	for _, source := range sc.AllSources() {
		if !source.ValidateConfiguration() {
			configIsOK = false
		}
	}
	return configIsOK
}

func (sc *SourcesConfiguration) Credentials() ([]Credentials, error) {
	credentials := []Credentials{}
	for _, source := range sc.AllSources() {
		if newCredentials, err := source.Credentials(); err != nil {
			return nil, err
		} else {
			credentials = append(credentials, newCredentials...)
		}
	}
	return credentials, nil
}
