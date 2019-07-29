package cli

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Parses and validates the given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := configuration.Sources.ValidateConfiguration(); err != nil {
			log.Fatalf("The sources section of the config file is invalid: %v", err)
		}
		if err := configuration.Targets.ValidateConfiguration(); err != nil {
			log.Fatalf("The targets section of the config file is invalid: %v", err)
		}
		log.Info("The config file is valid!")
	},
}
