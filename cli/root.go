package cli

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mitchellh/mapstructure"

	"gopkg.in/yaml.v2"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/coveooss/credentials-sync/sync"
	"github.com/coveooss/credentials-sync/targets"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configuration *sync.Configuration

var rootCmd = &cobra.Command{
	Use:   "credentials-sync",
	Short: "Fetches credentials and syncs them to targets",
	Long: `Grabs credentials from various sources and
	syncs them to the given targets. This CLI is useful for
	targets that do not support external credentials.
	Support Jenkins only for now.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var (
			configurationDict = map[string]interface{}{}
			configurationFile = viper.GetString("config")
			err               error
			fileContent       []byte
		)

		level, err := log.ParseLevel(viper.GetString("log-level"))
		if err != nil {
			return fmt.Errorf("Invalid log level: %s", err)
		}
		log.SetLevel(level)

		if configurationFile == "" {
			return fmt.Errorf("A configuration file must be defined")
		}

		configuration = sync.NewConfiguration()
		sourcesConfiguration := &credentials.SourcesConfiguration{}
		targetsConfiguration := &targets.Configuration{}

		if strings.HasPrefix(configurationFile, "s3://") {
			sess := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))
			s3Client := s3.New(sess)
			splitS3Path, err := url.Parse(configurationFile)
			if err != nil {
				return fmt.Errorf("Failed to parse the given S3 config path: %v", err)
			}

			resp, err := s3Client.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(splitS3Path.Host),
				Key:    aws.String(splitS3Path.Path),
			})
			if err != nil {
				return fmt.Errorf("Failed to download the config file from S3, %v", err)
			}

			if fileContent, err = ioutil.ReadAll(resp.Body); err != nil {
				return fmt.Errorf("Failed to read the config file from S3, %v", err)
			}
		} else {
			if fileContent, err = ioutil.ReadFile(configurationFile); err != nil {
				return err
			}
		}

		if err = yaml.Unmarshal(fileContent, configurationDict); err != nil {
			return err
		}

		// Get the config
		if err = mapstructure.Decode(configurationDict, configuration); err != nil {
			return err
		}

		// Get sources from config
		if err = mapstructure.Decode(configurationDict["sources"], sourcesConfiguration); err != nil {
			return err
		}
		configuration.SetSources(sourcesConfiguration)

		// Get targets from config
		if err = mapstructure.Decode(configurationDict["targets"], targetsConfiguration); err != nil {
			return err
		}
		configuration.SetTargets(targetsConfiguration)

		return nil
	},
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	viper.AutomaticEnv()
	viper.SetEnvPrefix("sync")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	rootCmd.PersistentFlags().StringP("config", "c", "", "configuration file")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().StringP("log-level", "l", log.InfoLevel.String(), `"debug", "info", "warning" or "error"`)
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))

	initListCredentials()
	rootCmd.AddCommand(listTargetsCmd, syncCmd, validateCmd)
}

// Execute runs the CLI
func Execute(commit string, date string, version string) {
	rootCmd.Version = fmt.Sprintf("%s %s (%s)", version, commit, date)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
