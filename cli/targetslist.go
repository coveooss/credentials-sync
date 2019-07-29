package cli

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listTargetsCmd = &cobra.Command{
	Use:   "list-targets",
	Short: "Resolves and lists all configured targets",
	Run: func(cmd *cobra.Command, args []string) {
		if err := configuration.Targets.ValidateConfiguration(); err != nil {
			log.Fatalf("The targets section of the config file is invalid: %v", err)
		}
		for _, target := range configuration.Targets.AllTargets() {
			fmt.Println(target.ToString())
		}
	},
}
