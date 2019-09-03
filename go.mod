module github.com/coveooss/credentials-sync

go 1.12

require (
	github.com/aws/aws-sdk-go v1.23.13
	github.com/bndr/gojenkins v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.3.1
	github.com/hashicorp/go-multierror v1.0.0
	github.com/jinzhu/copier v0.0.0-20190625015134-976e0346caa8
	github.com/mitchellh/mapstructure v1.1.2
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/bndr/gojenkins => github.com/julienduchesne/gojenkins v2.1.0+incompatible
