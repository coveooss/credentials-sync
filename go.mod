module github.com/coveooss/credentials-sync

go 1.12

require (
	github.com/aws/aws-sdk-go v1.36.28
	github.com/bndr/gojenkins v1.0.1
	github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d // indirect
	github.com/evalphobia/logrus_sentry v0.8.2
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/getsentry/sentry-go v0.11.0
	github.com/golang/mock v1.4.4
	github.com/hashicorp/go-multierror v1.1.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/bndr/gojenkins => github.com/julienduchesne/gojenkins v2.1.0+incompatible
