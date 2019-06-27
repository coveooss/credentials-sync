package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/mapstructure"

	"gopkg.in/yaml.v2"

	"github.com/coveooss/credentials-sync/sync"

	"github.com/spf13/cobra"
)

var (
	configurationFile string
	configuration     *sync.Configuration
	verbose           bool
)

var rootCmd = &cobra.Command{
	Use:   "credentials-sync",
	Short: "Fetches credentials and syncs them to targets",
	Long: `Grabs credentials from various sources and
	syncs them to the given targets. This CLI is useful for
	targets that do not support external credentials.
	Support Jenkins only for now.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		configuration = sync.NewConfiguration()
		configurationDict := map[string]interface{}{}
		var (
			err         error
			fileContent []byte
		)
		if fileContent, err = ioutil.ReadFile(configurationFile); err != nil {
			return err
		}
		if err = yaml.Unmarshal(fileContent, configurationDict); err != nil {
			return err
		}
		return mapstructure.Decode(configurationDict, configuration)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configurationFile, "config", "c", "", "configuration file")
	rootCmd.MarkPersistentFlagRequired("config")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
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
