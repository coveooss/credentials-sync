package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	showSensitiveAttributes bool
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
			fmt.Println(credentials.ToString(showSensitiveAttributes))
		}
	},
}

func initListCredentials() {
	listCredentialsCmd.Flags().BoolVarP(&showSensitiveAttributes, "show-sensitive", "s", false, "show sensitive credentials attributes, such as passwords")
	rootCmd.AddCommand(listCredentialsCmd)
}
