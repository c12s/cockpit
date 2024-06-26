package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const (
	shortLoginDescription = "Login into application"
	longLoginDescription  = `Input your username after that you will be prompted to input your password.
Your token will be saved in the token.txt file, which will be sent with all of your request headers.

Example:
- cockpit login --username "username"`
)

var (
	tokenResponse model.TokenResponse
)

var LoginCmd = &cobra.Command{
	Use:     "login",
	Aliases: aliases.LoginAliases,
	Short:   shortLoginDescription,
	Long:    longLoginDescription,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{usernameFlag})
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

		fmt.Println("Login successful!\n")
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
		Timeout:     10 * time.Second,
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
	LoginCmd.Flags().StringVarP(&username, usernameFlag, usernameShorthandFlag, "", "Username for login")
	LoginCmd.MarkFlagRequired(usernameFlag)
}
