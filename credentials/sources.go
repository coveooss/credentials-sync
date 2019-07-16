package credentials

import (
	"fmt"
	"sort"

	"gopkg.in/yaml.v2"
)

// Source represents a location to fetch credentials
type Source interface {
	Credentials() ([]Credentials, error)
	Type() string
	ValidateConfiguration() bool
}

// SourcesConfiguration contains all configured sources
type SourcesConfiguration struct {
	AWSS3Sources            []*AWSS3Source             `mapstructure:"aws_s3"`
	AWSSecretsManagerSource []*AWSSecretsManagerSource `mapstructure:"aws_secretsmanager"`
	LocalSources            []*LocalSource             `mapstructure:"local"`

	credentialsList []Credentials
}

type SourceCollection interface {
	AllSources() []Source
	Credentials() ([]Credentials, error)
	ValidateConfiguration() bool
}

// AllSources returns all configured sources in a single list
func (sc *SourcesConfiguration) AllSources() []Source {
	sources := []Source{}
	for _, source := range sc.LocalSources {
		sources = append(sources, source)
	}
	for _, source := range sc.AWSS3Sources {
		sources = append(sources, source)
	}
	for _, source := range sc.AWSSecretsManagerSource {
		sources = append(sources, source)
	}
	return sources
}

// ValidateConfiguration verifies that all configured sources are correctly configured
func (sc *SourcesConfiguration) ValidateConfiguration() bool {
	configIsOK := true
	for _, source := range sc.AllSources() {
		if !source.ValidateConfiguration() {
			configIsOK = false
		}
	}
	return configIsOK
}

// Credentials extracts credentials from all configured sources
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

func getCredentialsFromBytes(byteArray []byte) ([]Credentials, error) {
	var (
		err               error
		yamlContentAsList []map[string]interface{}
		yamlContentAsMap  map[string]map[string]interface{}
	)
	if err = yaml.Unmarshal(byteArray, &yamlContentAsList); err != nil {
		if err = yaml.Unmarshal(byteArray, &yamlContentAsMap); err != nil {
			return nil, err
		}
		yamlContentAsList = []map[string]interface{}{}
		for id, value := range yamlContentAsMap {
			value["id"] = id
			yamlContentAsList = append(yamlContentAsList, value)
		}
	}
	return ParseCredentials(yamlContentAsList)
}
