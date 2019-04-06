package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listTargetsCmd = &cobra.Command{
	Use:   "list-targets",
	Short: "Resolves and lists all configured targets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}
