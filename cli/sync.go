package cli

import (
	"github.com/coveooss/credentials-sync/logger"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetches credentials and syncs them to targets",
	Run: func(cmd *cobra.Command, args []string) {
		if err := configuration.Sources.ValidateConfiguration(); err != nil {
			logger.Log.Fatalf("The sources section of the config file is invalid: %v", err)
		}
		if err := configuration.Targets.ValidateConfiguration(); err != nil {
			logger.Log.Fatalf("The targets section of the config file is invalid: %v", err)
		}
		if err := configuration.Sync(); err != nil {
			logger.Log.Fatalf("The synchronization process failed: %v", err)
		}
	},
}
