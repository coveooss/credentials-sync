package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/coveo/credentials-sync/sync"

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
		configuration = &sync.Configuration{}
		if configurationFile != "" {
			var (
				err         error
				fileContent []byte
			)
			if fileContent, err = ioutil.ReadFile(configurationFile); err != nil {
				return err
			}
			return yaml.Unmarshal(fileContent, configuration)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configurationFile, "config", "c", "", "configuration file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	initListCredentials()
	rootCmd.AddCommand(listTargetsCmd, syncCmd, validateCmd)
}

func Execute(commit string, date string, version string) {
	rootCmd.Version = fmt.Sprintf("%s %s (%s)", version, commit, date)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
