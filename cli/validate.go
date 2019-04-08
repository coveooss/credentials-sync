package cli

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Parses and validates the given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if !configuration.Sources.ValidateConfiguration() {
			log.Fatal("The sources section of the config file is invalid")
		}
		if !configuration.Targets.ValidateConfiguration() {
			log.Fatal("The targets section of the config file is invalid")
		}
		log.Info("The config file is valid!")
	},
}
