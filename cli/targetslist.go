package cli

import (
	"fmt"
	"github.com/coveooss/credentials-sync/logger"

	"github.com/spf13/cobra"
)

var listTargetsCmd = &cobra.Command{
	Use:   "list-targets",
	Short: "Resolves and lists all configured targets",
	Run: func(cmd *cobra.Command, args []string) {
		if err := configuration.Targets.ValidateConfiguration(); err != nil {
			logger.Log.Fatalf("The targets section of the config file is invalid: %v", err)
		}
		for _, target := range configuration.Targets.AllTargets() {
			fmt.Println(target.ToString())
		}
	},
}
