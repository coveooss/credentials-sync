package cli

import (
	"github.com/coveooss/credentials-sync/logger"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Parses and validates the given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := configuration.Sources.ValidateConfiguration(); err != nil {
			logger.Log.Fatalf("The sources section of the config file is invalid: %v", err)
		}
		if err := configuration.Targets.ValidateConfiguration(); err != nil {
			logger.Log.Fatalf("The targets section of the config file is invalid: %v", err)
		}
		logger.Log.Info("The config file is valid!")
	},
}
