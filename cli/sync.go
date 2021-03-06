package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetches credentials and syncs them to targets",
	Run: func(cmd *cobra.Command, args []string) {
		if err := configuration.Sources.ValidateConfiguration(); err != nil {
			log.Fatalf("The sources section of the config file is invalid: %v", err)
		}
		if err := configuration.Targets.ValidateConfiguration(); err != nil {
			log.Fatalf("The targets section of the config file is invalid: %v", err)
		}
		if err := configuration.Sync(); err != nil {
			log.Fatalf("The synchronization process failed: %v", err)
		}
	},
}
