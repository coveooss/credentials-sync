package targets

import (
	"errors"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/bndr/gojenkins"
	"github.com/coveo/credentials-sync/credentials"
)

type JenkinsTarget struct {
	Base `mapstructure:",squash"`

	LoginCredentials *credentials.Credentials
	URL              string
	SecureConnection bool

	client              *gojenkins.Jenkins
	existingCredentials []credentials.Credentials
}

func (jenkins *JenkinsTarget) Initialize() error {
	auth := (*jenkins.LoginCredentials).(*credentials.UsernamePasswordCredentials)
	jenkins.client = gojenkins.CreateJenkins(jenkins.URL, auth.Username, auth.Password)
	jenkins.client.Requester.SslVerify = jenkins.SecureConnection
	jenkins.client.Init()
	var err error
	defer func() {
		if r := recover(); r != nil {
			log.Error("Recovered the following error while initializing Jenkins: ", r)
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
	return err
}

func (jenkins *JenkinsTarget) ToString() string {
	return fmt.Sprintf("%s (Jenkins) - %s", jenkins.BaseToString(), jenkins.URL)
}

func (jenkins *JenkinsTarget) UpdateListOfCredentials(listOfCredentials []*credentials.Credentials) error {
	for _, credentials := range listOfCredentials {
		if err := jenkins.UpdateCredentials(credentials); err != nil {
			return err
		}
	}
	return nil
}

func (jenkins *JenkinsTarget) UpdateCredentials(credentials *credentials.Credentials) error {
	return nil
}

func (jenkins *JenkinsTarget) ValidateConfiguration() bool {
	if _, err := url.ParseRequestURI(jenkins.URL); err != nil {
		log.Errorf("The Jenkins target `%s` has an invalid URL: %s", jenkins.Name, jenkins.URL)
		return false
	}
	return true
}
