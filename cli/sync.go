package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetches credentials and syncs them to targets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(verbose)
	},
}
