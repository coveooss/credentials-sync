package credentials

import (
	"fmt"
	"sort"

	"github.com/coveooss/credentials-sync/logger"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

// Source represents a location to fetch credentials
type Source interface {
	Credentials() ([]Credentials, error)
	Type() string
	ValidateConfiguration() error
}

// SourcesConfiguration contains all configured sources
type SourcesConfiguration struct {
	AWSS3Sources            []*AWSS3Source             `mapstructure:"aws_s3"`
	AWSSecretsManagerSource []*AWSSecretsManagerSource `mapstructure:"aws_secretsmanager"`
	LocalSources            []*LocalSource             `mapstructure:"local"`

	credentialsList []Credentials
}

// SourceCollection represents a collection of sources from which credentials can be fetched
type SourceCollection interface {
	AllSources() []Source
	Credentials() ([]Credentials, error)
	ValidateConfiguration() error
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
func (sc *SourcesConfiguration) ValidateConfiguration() error {
	var validationErrors error
	for _, source := range sc.AllSources() {
		if err := source.ValidateConfiguration(); err != nil {
			validationErrors = multierror.Append(validationErrors, err)
		}
	}
	return validationErrors
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
		err             error
		credentialsList []map[string]interface{}
	)

	methods := []func(bytes []byte) ([]map[string]interface{}, error){tryReadingList, tryReadingMapOfCredentials, tryReadingSingleCredential}
	var success = false
	var errors []error
	for _, method := range methods {
		if credentialsList, err = method(byteArray); err == nil {
			success = true
			return ParseCredentials(credentialsList)
		}
		errors = append(errors, err)
	}

	if !success {
		logger.Log.Warning("Failed to get credential from data using all known formats (details below)")
		for _, err := range errors {
			if err != nil {
				logger.Log.Warning(err)
			}
		}
	}

	return nil, fmt.Errorf("Failed to parse %v. See debug for more info", string(byteArray))
}

// Accept list of credentials
func tryReadingList(bytes []byte) ([]map[string]interface{}, error) {
	var credentialsList []map[string]interface{}
	if err := yaml.Unmarshal(bytes, &credentialsList); err != nil {
		return nil, fmt.Errorf("Error reading as credentials list: %v", err)
	}

	return credentialsList, nil
}

// Accept map of credentials
func tryReadingMapOfCredentials(bytes []byte) ([]map[string]interface{}, error) {
	credentialsList := []map[string]interface{}{}

	var credentialsMap map[string]map[string]interface{}
	if err := yaml.Unmarshal(bytes, &credentialsMap); err != nil {
		return nil, fmt.Errorf("Error reading as credentials map: %v", err)
	}

	for id, value := range credentialsMap {
		value["id"] = id
		credentialsList = append(credentialsList, value)
	}

	return credentialsList, nil
}

// Accept a single credential
func tryReadingSingleCredential(bytes []byte) ([]map[string]interface{}, error) {
	var singleCredentials map[string]interface{}
	if err := yaml.Unmarshal(bytes, &singleCredentials); err != nil {
		return nil, fmt.Errorf("Error reading as a map: %v", err)
	}

	id, gotID := singleCredentials["id"]
	if !gotID {
		return nil, fmt.Errorf("The parsed credentials doesn't have an ID: %v", singleCredentials)
	}

	if _, idIsString := id.(string); !idIsString {
		return nil, fmt.Errorf("The given credentials' ID is not a string: %v", singleCredentials)
	}

	return []map[string]interface{}{singleCredentials}, nil
}
