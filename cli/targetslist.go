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
		if !configuration.Targets.ValidateConfiguration() {
			log.Fatal("The targets section of the config file is invalid")
		}
		for _, target := range configuration.Targets.AllTargets() {
			fmt.Println(target.ToString())
		}
	},
}
