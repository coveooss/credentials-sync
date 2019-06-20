module github.com/coveo/credentials-sync

go 1.12

require (
	github.com/aws/aws-sdk-go v1.19.11
	github.com/bndr/gojenkins v0.0.0-00010101000000-000000000000
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/testify v1.2.2
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/bndr/gojenkins => github.com/julienduchesne/gojenkins v2.0.1+incompatible
