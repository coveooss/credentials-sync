package targets

import (
	"errors"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/bndr/gojenkins"
	"github.com/coveo/credentials-sync/credentials"
)

const credentialsDomain = "_"

type JenkinsTarget struct {
	Base `mapstructure:",squash"`

	CredentialsID    *string `mapstructure:"credentials_id"`
	URL              string
	SecureConnection bool `mapstructure:"secure_connections"`

	client              *gojenkins.Jenkins
	credentialsManager  *gojenkins.CredentialsManager
	existingCredentials []string
	loginCredentials    credentials.Credentials
}

func (jenkins *JenkinsTarget) GetName() string {
	return jenkins.Name
}

func (jenkins *JenkinsTarget) Initialize(allCredentials []credentials.Credentials) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()
	for _, credentials := range allCredentials {
		if jenkins.CredentialsID != nil && credentials.GetID() == *jenkins.CredentialsID {
			jenkins.loginCredentials = credentials
		}
	}
	if jenkins.loginCredentials != nil {
		auth := jenkins.loginCredentials.(*credentials.UsernamePasswordCredentials)
		jenkins.client = gojenkins.CreateJenkins(nil, jenkins.URL, auth.Username, auth.Password)
	} else {
		jenkins.client = gojenkins.CreateJenkins(nil, jenkins.URL)
	}

	jenkins.client.Requester.SslVerify = jenkins.SecureConnection
	jenkins.client.Init()
	jenkins.credentialsManager = &gojenkins.CredentialsManager{
		J: jenkins.client,
	}

	jenkins.existingCredentials, err = jenkins.credentialsManager.List(credentialsDomain)

	return err
}

func (jenkins *JenkinsTarget) HasCredentials(cred credentials.Credentials) bool {
	for _, id := range jenkins.existingCredentials {
		if cred.GetID() == id {
			return true
		}
	}
	return false
}

func (jenkins *JenkinsTarget) ToString() string {
	return fmt.Sprintf("%s (Jenkins) - %s", jenkins.BaseToString(), jenkins.URL)
}

func (jenkins *JenkinsTarget) UpdateListOfCredentials(listOfCredentials []credentials.Credentials) error {
	for _, credentials := range listOfCredentials {
		if err := jenkins.UpdateCredentials(credentials); err != nil {
			return err
		}
	}
	return nil
}

func (jenkins *JenkinsTarget) UpdateCredentials(cred credentials.Credentials) error {
	jenkinsCred := toJenkinsCredential(cred)
	if jenkinsCred == nil {
		return fmt.Errorf("Unable to create jenkins credentials from %s", cred.GetID())
	}
	if jenkins.HasCredentials(cred) {
		return jenkins.credentialsManager.Update(credentialsDomain, cred.GetID(), jenkinsCred)
	}
	return jenkins.credentialsManager.Add(credentialsDomain, jenkinsCred)
}

func (jenkins *JenkinsTarget) ValidateConfiguration() bool {
	if _, err := url.ParseRequestURI(jenkins.URL); err != nil {
		log.Errorf("The Jenkins target `%s` has an invalid URL: %s", jenkins.Name, jenkins.URL)
		return false
	}
	return true
}

func toJenkinsCredential(creds credentials.Credentials) interface{} {
	switch creds.(type) {
	case *credentials.SecretTextCredentials:
		castCreds := creds.(*credentials.SecretTextCredentials)
		return &gojenkins.StringCredentials{
			ID:          creds.GetID(),
			Description: castCreds.Description,
			Secret:      castCreds.Secret,
		}
	case *credentials.UsernamePasswordCredentials:
		castCreds := creds.(*credentials.UsernamePasswordCredentials)
		return &gojenkins.UsernameCredentials{
			ID:          castCreds.GetID(),
			Description: castCreds.Description,
			Username:    castCreds.Username,
			Password:    castCreds.Password,
		}
	}
	return nil
}
