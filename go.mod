module github.com/coveooss/credentials-sync

go 1.12

require (
	github.com/aws/aws-sdk-go v1.31.4
	github.com/bndr/gojenkins v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.4.3
	github.com/hashicorp/go-multierror v1.1.0
	github.com/mitchellh/mapstructure v1.2.2
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/bndr/gojenkins => github.com/julienduchesne/gojenkins v2.1.0+incompatible
