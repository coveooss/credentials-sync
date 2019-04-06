package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Parses and validates the given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if !configuration.Sources.ValidateConfiguration() {
			os.Exit(1)
		}
	},
}
