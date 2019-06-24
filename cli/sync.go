package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetches credentials and syncs them to targets",
	Run: func(cmd *cobra.Command, args []string) {
		if !configuration.Sources.ValidateConfiguration() {
			log.Fatal("The sources section of the config file is invalid")
		}
		if !configuration.Targets.ValidateConfiguration() {
			log.Fatal("The targets section of the config file is invalid")
		}
		configuration.Sync()
	},
}
