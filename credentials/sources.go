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
	credentialsList := []Credentials{}

	// Fetch all credentials
	for _, source := range sc.AllSources() {
		newCredentials, err := source.Credentials()
		if err != nil {
			return nil, err
		}
		credentialsList = append(credentialsList, newCredentials...)
	}

	// Sort credentials by ID
	sort.Slice(credentialsList[:], func(i, j int) bool {
		return credentialsList[i].GetID() < credentialsList[j].GetID()
	})

	// Throw an error if IDs are not unique
	credentialIds := map[string]bool{}
	for _, cred := range credentialsList {
		if _, ok := credentialIds[cred.GetID()]; ok {
			return nil, fmt.Errorf("There more than one credentials with this ID: %s", cred.GetID())
		}
		credentialIds[cred.GetID()] = true
	}

	return credentialsList, nil
}
