package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var tokenResponse model.TokenResponse

var LoginCmd = &cobra.Command{
	Use:     "login",
	Aliases: aliases.LoginAliases,
	Short:   constants.ShortLoginDesc,
	Long:    constants.LongLoginDesc,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.UsernameFlag})
	},
	Run: func(cmd *cobra.Command, args []string) {
		password, err := utils.PromptForPassword()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		err = login(username, password)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println("Login successful!")
	},
}

func login(username, password string) error {
	credentials := model.Credentials{
		Username: username,
		Password: password,
	}

	url := clients.BuildURL("core", "v1", "LoginUser")

	err := utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		RequestBody: credentials,
		Response:    &tokenResponse,
		Timeout:     30 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("failed to send login request: %v", err)
	}

	if err := utils.SaveTokenToFile(tokenResponse.Token); err != nil {
		return fmt.Errorf("failed to save token: %v", err)
	}

	return nil
}

func init() {
	LoginCmd.Flags().StringVarP(&username, constants.UsernameFlag, constants.UsernameShorthandFlag, "", "Username for login")
	LoginCmd.MarkFlagRequired(constants.UsernameFlag)
}
