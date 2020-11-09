module github.com/coveooss/credentials-sync

go 1.12

require (
	github.com/aws/aws-sdk-go v1.35.23
	github.com/bndr/gojenkins v1.0.1
	github.com/golang/mock v1.4.4
	github.com/hashicorp/go-multierror v1.1.0
	github.com/mitchellh/mapstructure v1.3.3
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/bndr/gojenkins => github.com/julienduchesne/gojenkins v2.1.0+incompatible
