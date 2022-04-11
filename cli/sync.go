package cli

import (
	"github.com/coveooss/credentials-sync/logger"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetches credentials and syncs them to targets",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := configuration.Sources.ValidateConfiguration(); err != nil {
			logger.Log.Errorf("The sources section of the config file is invalid: %v", err)
			return err
		}
		if err := configuration.Targets.ValidateConfiguration(); err != nil {
			logger.Log.Errorf("The targets section of the config file is invalid: %v", err)
			return err
		}
		if err := configuration.Sync(); err != nil {
			logger.Log.Errorf("The synchronization process failed: %v", err)
			return err
		}
		return nil
	},
}
