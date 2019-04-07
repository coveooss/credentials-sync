package credentials

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

const ssmPathRegex = `^/(?:[A-Za-z0-9]+/?)+$`

type AWSSSMSource struct {
	Path string
}

func (source *AWSSSMSource) Credentials() ([]Credentials, error) {
	return []Credentials{}, nil
}

func (source *AWSSSMSource) Type() string {
	return "Amazon SSM"
}

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
