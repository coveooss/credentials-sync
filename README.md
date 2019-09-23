# credentials-sync
[![Build Status](https://travis-ci.org/coveooss/credentials-sync.svg?branch=master)](https://travis-ci.org/coveooss/credentials-sync)
[![codecov](https://codecov.io/gh/coveooss/credentials-sync/branch/master/graph/badge.svg)](https://codecov.io/gh/coveooss/credentials-sync)
[![Go Report Card](https://goreportcard.com/badge/github.com/coveooss/credentials-sync)](https://goreportcard.com/report/github.com/coveooss/credentials-sync)

Sync credentials from various sources to various targets. It currently only supports Jenkins, but LastPass is planned because that is what we use. However, we are open to supporting more targets.

What's the point?
1. Easier credentials rotations. Rotating credentials manually is simply not an option when credentials rotations are done too often
2. Uses a push-model instead of a pull-model which means that you can put your credentials in a secure environment to which targets don't have access, targets may have varying degrees of security (prod vs dev)
3. Decouples your credentials and the systems which use these credentials. Standardized credentials format for all targets

## Installation

 - Go to https://github.com/coveooss/credentials-sync/releases
 - Download the file appropriate for your system
 - Unzip it

## Usage

 - Write a config file, see [format here](#configuration-file)
 - Run the sync command

```bash
credentials-sync sync -c config.yml
```

## Logging

The log level can be set with either:
  - The `--log-level` option
  - The `SYNC_LOG_LEVEL` env variable

Valid levels are `debug`, `info`, `warning` and `error`

![example](https://raw.githubusercontent.com/coveooss/credentials-sync/master/example.png)

## Configuration file
A configuration file must be given to the application. Its path can either be a local path or a S3 path
The path can either be passed as a parameter (`-c/--config`) or as an environment variable (`SYNC_CONFIG`)  

A configuration file contains [sources](#supported-sources) which contain [credentials](#supported-types-of-credentials). It also defines targets to which these credentials will be synced

Here is the accepted format:
```yaml
sources:
  local:
    - path: /home/jdoe/path/to/file.yaml
  aws_s3:
    - bucket: name
      key: path/to/file.yaml
  aws_secretsmanager:
    - secret_prefix: credentials-sync/
    - secret_id: arn:aws:secretsmanager:us-west-2:123456789012:secret:production/MyAwesomeAppSecret-a1b2c3
stop_on_error: true   # If true, will completely stop the process if an operation fails. Otherwise, continues anyways
target_parallelism: 3 # Number of target on which to sync creds at the same time
targets:
  jenkins:
    - name: toolsjenkins
      url: https://toolsjenkins.my-domain.com
      credentials_id: toolsjenkins # Uses a set of username:password credentials
```

## Supported sources
Here are the supported sources:
 - Local (Single file)
 - AWS S3 (Single object)
 - AWS SecretsManager (Single secret or a secret prefix)

The source's value must either be a list or a map in the following formats (JSON or YAML):
```yaml
# list
- id: my_cred
  description: a description
  ...
- id: my_other_cred
  ...

# map
my_cred:
  description: a description
  ...
my_other_cred:
  ...
```

## Supported types of credentials
Credentials are defined as JSON, here are the supported types of credentials with definition examples:
 - Secret text
```yaml
secret_text:
  description: A secret text cred is only composed of a secret
  type: secret
  secret: xoxb-a-slack-token
```
 - Username password
```yaml
username_password:
  description: A username:password cred is composed of two values, a username and a password
  type: usernamepassword
  username: jdoe
  password: hunter42
```
 - AWS IAM credentials
```yaml
aws_iam:
  description: IAM creds are composed of an access key, a secret access key and optionally a role to assume
  type: aws
  access_key: AKIAMYFAKEKEY
  secret_key: fdjVEsefk4kgjVsdjfew54
  role_arn: arn:aws:iam::123456789012:role/S3Access
```
 - SSH credentials
```yaml
ssh_key:
  description: An SSH key is composed of a private key, a username and optionally, a passphrase
  type: ssh
  username: jdoe
  passphrase: hunter42
  private_key: |
    -----BEGIN RSA PRIVATE KEY-----
    MIIJKAIBAAKCAgEAvlXlTSaTs2VfvtYM+6UF9AQVhd6V7DeU1ViMQaLmEvWHkd/y
    vZEipSq+rI3vis0ObviouslyNotValidlollyiIkfd7bIoHRnQXCV8le/dzXBiAt
    Pa1fxCEkcsjJjwWBwDrGw3qGS6T+ElgngisI7YKBXrVqVKQaEeUEeqACwI5j9LoZ
    yuuoMQ9TLNjWFfcDK5Pl/0RhWrGEZDFAaSHm/lyLhxHvWR0GYSJ2V9XG2UTi6Zdq
    i11Ol956OqlmjyzpqmyYFCBzzhv7uLlI31/0MZfjQQJUa1JeQCL+Usjj+3GICDlu
    Yi7xX7n5GW4h7w43KXH1HHV+J1BE3w53uuzm+cATnEWc/raNopeDontEvenTryIt
    -----END RSA PRIVATE KEY-----
```

## Other features

### Unsynced credentials
Since credentials are also used for authentication, you may wish to not sync them:
```yaml
toolsjenkins:
  description: Login credentials for jenkins
  type: usernamepassword
  username: jdoe
  password: apikey
  no_sync: true # This will prevent this cred from being synced
```

### Target matching
Sometimes, certain credentials should only be synced to certain targets. There are two ways to make sure this happens:

1. Matching on target's name
```yaml
secret_text:
  description: A secret text cred is only composed of a secret
  type: secret
  secret: xoxb-a-slack-token
  target: toolsjenkins # This cred will only be synced to the toolsjenkins target
```
2. Matching on target tags
```yaml
# In config file
targets:
  jenkins:
    - name: toolsjenkins
      url: https://toolsjenkins.my-domain.com
      credentials_id: toolsjenkins
      tags:
        my_tag: my_value

# In credentials definition
secret_text:
  description: A secret text cred is only composed of a secret
  type: secret
  secret: xoxb-a-slack-token
  target_tags:
    do_match:
      my_tag: value # Will only sync to targets if my_tag == "value"
      my_other_tag: ["value1", "value2"] # OR if my_other_tag == "value1" OR my_other_tag == "value2"
    dont_match:
      my_tag: ["other_value", "some_value"] # Will not sync to targets if my_tag == "other_value" or if my_tag == "some_value", regardless of `do_match`
```

## Using the docker image

For every version, a docker image is published here: https://hub.docker.com/r/coveo/credentials-sync  
The only parameter needed for the credentials sync is the configuration file (You can set its location with `SYNC_CONFIG` env variable)  
This allows you to run this as a cron job in AWS Fargate or Kubernetes, for example

## Roadmap
- Incremental runs (keep a state file and only update credentials that have been modified at the source level. This would have to be optional because full runs will still be need to sync back credentials that have been modified at the target level)
- LastPass target
- Terraform state file source
- SSM Parameter store source (not in the regular JSON format)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
