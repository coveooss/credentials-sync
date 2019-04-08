package cli

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	showSensitiveAttributes bool
)

var listCredentialsCmd = &cobra.Command{
	Use:   "list-credentials",
	Short: "Resolves and lists all configured credentials",
	Run: func(cmd *cobra.Command, args []string) {
		if !configuration.Sources.ValidateConfiguration() {
			log.Fatal("The sources section of the config file is invalid")
		}
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
