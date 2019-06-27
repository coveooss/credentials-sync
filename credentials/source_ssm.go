package credentials

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	log "github.com/sirupsen/logrus"
)

const ssmPathRegex = `^/(?:[A-Za-z0-9]+/?)+$`

// AWSSSMSource represents AWS SSM Parameters containing credentials
type AWSSSMSource struct {
	Path string
}

// Credentials extracts credentials from the source
func (source *AWSSSMSource) Credentials() ([]Credentials, error) {
	svc := ssm.New(session.New())
	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(source.Path),
		WithDecryption: aws.Bool(true),
		Recursive:      aws.Bool(true),
	}
	credentialsMaps := []map[string]interface{}{}
	if err := svc.GetParametersByPathPages(input,
		func(page *ssm.GetParametersByPathOutput, lastPage bool) bool {
			for _, parameter := range page.Parameters {
				splitName := strings.Split(*parameter.Name, "/")
				credentialsMap := map[string]interface{}{
					"full_name":   *parameter.Name,
					"id":          splitName[len(splitName)-1],
					"description": splitName[len(splitName)-1],
					"value":       *parameter.Value,
				}
				credentialsMaps = append(credentialsMaps, credentialsMap)
			}
			return !lastPage
		}); err != nil {
		return nil, err
	}
	return ParseCredentials(credentialsMaps)
}

// Type returns the type of the source
func (source *AWSSSMSource) Type() string {
	return "Amazon SSM"
}

// ValidateConfiguration verifies that the source's attributes are valid
func (source *AWSSSMSource) ValidateConfiguration() bool {
	if strings.HasPrefix(source.Path, "/aws") {
		log.Errorf("%s should not start with /aws. This path is reserved to AWS", source.Path)
		return false
	}
	if matched, err := regexp.MatchString(ssmPathRegex, source.Path); !matched || err != nil {
		log.Errorf("%s does not match the allowed SSM path regex: %s", source.Path, ssmPathRegex)
		return false
	}
	return true
}
