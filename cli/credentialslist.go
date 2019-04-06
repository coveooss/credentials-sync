package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCredentialsCmd = &cobra.Command{
	Use:   "list-credentials",
	Short: "Resolves and lists all configured credentials",
	Run: func(cmd *cobra.Command, args []string) {
		allCredentials, err := configuration.Sources.Credentials()
		if err != nil {
			panic(err)
		}
		for _, credentials := range allCredentials {
			fmt.Println(credentials.ToString())
		}
	},
}
