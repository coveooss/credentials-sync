package cli

import (
	"fmt"

	"github.com/coveooss/credentials-sync/logger"

	"github.com/spf13/cobra"
)

var (
	showSensitiveAttributes bool
)

var listCredentialsCmd = &cobra.Command{
	Use:   "list-credentials",
	Short: "Resolves and lists all configured credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := configuration.Sources.ValidateConfiguration(); err != nil {
			logger.Log.Errorf("The sources section of the config file is invalid: %v", err)
			return err
		}
		allCredentials, err := configuration.Sources.Credentials()
		if err != nil {
			logger.Log.Errorf("The credential extraction for all configured sources failed: %v", err)
			return err
		}
		for _, credentials := range allCredentials {
			fmt.Println(credentials.ToString(showSensitiveAttributes))
		}
		return nil
	},
}

func initListCredentials() {
	listCredentialsCmd.Flags().BoolVarP(&showSensitiveAttributes, "show-sensitive", "s", false, "show sensitive credentials attributes, such as passwords")
	rootCmd.AddCommand(listCredentialsCmd)
}
