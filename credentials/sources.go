package credentials

import (
	"fmt"
	"sort"
)

type Source interface {
	Credentials() ([]Credentials, error)
	Type() string
	ValidateConfiguration() bool
}

type SourcesConfiguration struct {
	AWSS3Sources            []*AWSS3Source             `mapstructure:"aws_s3"`
	AWSSecretsManagerSource []*AWSSecretsManagerSource `mapstructure:"aws_secretsmanager"`
	AWSSSMSources           []*AWSSSMSource            `mapstructure:"aws_ssm"`
	LocalSources            []*LocalSource             `mapstructure:"local"`

	credentialsList []Credentials
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
	if sc.credentialsList != nil {
		return sc.credentialsList, nil
	}

	sc.credentialsList = []Credentials{}

	// Fetch all credentials
	for _, source := range sc.AllSources() {
		newCredentials, err := source.Credentials()
		if err != nil {
			return nil, err
		}
		sc.credentialsList = append(sc.credentialsList, newCredentials...)
	}

	// Sort credentials by ID
	sort.Slice(sc.credentialsList[:], func(i, j int) bool {
		return sc.credentialsList[i].GetID() < sc.credentialsList[j].GetID()
	})

	// Throw an error if IDs are not unique
	credentialIds := map[string]bool{}
	for _, cred := range sc.credentialsList {
		if _, ok := credentialIds[cred.GetID()]; ok {
			return nil, fmt.Errorf("There more than one credentials with this ID: %s", cred.GetID())
		}
		credentialIds[cred.GetID()] = true
	}

	return sc.credentialsList, nil
}
