package targets

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/bndr/gojenkins"
	"github.com/coveooss/credentials-sync/credentials"
)

const credentialsDomain = "_"

// JenkinsTarget represents a Jenkins instance
type JenkinsTarget struct {
	Base `mapstructure:",squash"`

	CredentialsID      *string `mapstructure:"credentials_id"`
	URL                string
	InsecureConnection bool `mapstructure:"insecure_connection"`

	client              *gojenkins.Jenkins
	credentialsManager  *gojenkins.CredentialsManager
	existingCredentials []string
	loginCredentials    credentials.Credentials
}

// Initialize executes all necessary operations to prepare the Jenkins target for sync
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
	for _, creds := range allCredentials {
		if jenkins.CredentialsID != nil && creds.GetID() == *jenkins.CredentialsID {
			jenkins.loginCredentials = creds
		}
	}
	options := &gojenkins.JenkinsOptions{SslVerify: aws.Bool(!jenkins.InsecureConnection)}
	if jenkins.loginCredentials != nil {
		auth := jenkins.loginCredentials.(*credentials.UsernamePasswordCredentials)
		options.Username = &auth.Username
		options.Password = &auth.Password
	}
	jenkins.client = gojenkins.CreateJenkinsWithOptions(jenkins.URL, options)

	_, err = jenkins.client.Init()

	if err != nil {
		return err
	}

	jenkins.credentialsManager = &gojenkins.CredentialsManager{
		J: jenkins.client,
	}

	jenkins.existingCredentials, err = jenkins.credentialsManager.List(credentialsDomain)

	return err
}

// ToString prints out a description of the Jenkins instance
func (jenkins *JenkinsTarget) ToString() string {
	return fmt.Sprintf("%s (Jenkins) - %s", jenkins.BaseToString(), jenkins.URL)
}

// GetExistingCredentials returns a list of all credential IDs on the target
func (jenkins *JenkinsTarget) GetExistingCredentials() []string {
	return jenkins.existingCredentials
}

// DeleteCredentials deletes the credentials with the given ID on the target
func (jenkins *JenkinsTarget) DeleteCredentials(id string) error {
	return jenkins.credentialsManager.Delete(credentialsDomain, id)
}

// UpdateCredentials syncs the given credentials to the Jenkins instance
func (jenkins *JenkinsTarget) UpdateCredentials(cred credentials.Credentials) error {
	jenkinsCred := toJenkinsCredential(cred)
	if jenkinsCred == nil {
		return fmt.Errorf("unable to create jenkins credentials from %s", cred.GetID())
	}
	if HasCredential(jenkins, cred.GetTargetID()) {
		return jenkins.credentialsManager.Update(credentialsDomain, cred.GetTargetID(), jenkinsCred)
	}
	return jenkins.credentialsManager.Add(credentialsDomain, jenkinsCred)
}

// ValidateConfiguration verifies that Jenkins configuration is valid
func (jenkins *JenkinsTarget) ValidateConfiguration() error {
	if _, err := url.ParseRequestURI(jenkins.URL); err != nil {
		return fmt.Errorf("the Jenkins target `%s` has an invalid URL: %s", jenkins.Name, jenkins.URL)
	}
	return nil
}

// JenkinsGithubAppCredentials is the Jenkins Github plugin's credentials configuration.
/*
   It must be serializable to the following XML:
	<org.jenkinsci.plugins.github__branch__source.GitHubAppCredentials plugin="github-branch-source@2.8.2">
		<id>github-app-dev</id>
		<description>The GitHub app for Jenkins</description>
		<appID>73157</appID>
		<privateKey>{some_private_key}</privateKey>
		<apiUri>https://api.github.com</apiUri>
		<owner>coveo</owner>
	</org.jenkinsci.plugins.github__branch__source.GitHubAppCredentials>
*/
type JenkinsGithubAppCredentials struct {
	XMLName     xml.Name `xml:"org.jenkinsci.plugins.github__branch__source.GitHubAppCredentials"`
	ID          string   `xml:"id"`
	Description string   `xml:"description,omitempty"`
	AppID       int      `xml:"appID"`
	PrivateKey  string   `xml:"privateKey"`
	APIURI      string   `xml:"apiUri,omitempty"`
	Owner       string   `xml:"owner,omitempty"`
}

func toJenkinsCredential(creds credentials.Credentials) interface{} {
	switch castCreds := creds.(type) {
	case *credentials.AmazonWebServicesCredentials:
		return &gojenkins.AmazonWebServicesCredentials{
			ID:                 creds.GetTargetID(),
			Description:        castCreds.GetDescriptionOrID(),
			AccessKey:          castCreds.AccessKey,
			SecretKey:          castCreds.SecretKey,
			IAMRoleARN:         castCreds.RoleARN,
			IAMMFASerialNumber: castCreds.MFASerialNumber,
		}
	case *credentials.SecretTextCredentials:
		return &gojenkins.StringCredentials{
			ID:          creds.GetTargetID(),
			Description: castCreds.GetDescriptionOrID(),
			Secret:      castCreds.Secret,
		}
	case *credentials.UsernamePasswordCredentials:
		return &gojenkins.UsernameCredentials{
			ID:          castCreds.GetTargetID(),
			Description: castCreds.GetDescriptionOrID(),
			Username:    castCreds.Username,
			Password:    castCreds.Password,
		}
	case *credentials.SSHCredentials:
		return &gojenkins.SSHCredentials{
			ID:          castCreds.GetTargetID(),
			Description: castCreds.GetDescriptionOrID(),
			Username:    castCreds.Username,
			Passphrase:  castCreds.Passphrase,
			PrivateKeySource: &gojenkins.PrivateKey{
				Class: gojenkins.KeySourceDirectEntryType,
				Value: castCreds.PrivateKey,
			},
		}
	case *credentials.GithubAppCredentials:
		return &JenkinsGithubAppCredentials{
			ID:          castCreds.GetTargetID(),
			Description: castCreds.GetDescriptionOrID(),
			AppID:       castCreds.AppID,
			PrivateKey:  castCreds.PrivateKey,
			APIURI:      "https://api.github.com",
			Owner:       castCreds.Owner,
		}
	}
	return nil
}
