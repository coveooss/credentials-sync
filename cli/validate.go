package cli

import (
	"github.com/coveooss/credentials-sync/logger"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Parses and validates the given configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := configuration.Sources.ValidateConfiguration(); err != nil {
			logger.Log.Errorf("The sources section of the config file is invalid: %v", err)
			return err
		}
		if err := configuration.Targets.ValidateConfiguration(); err != nil {
			logger.Log.Errorf("The targets section of the config file is invalid: %v", err)
			return err
		}
		logger.Log.Info("The config file is valid!")
		return nil
	},
}
