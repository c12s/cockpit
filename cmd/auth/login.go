package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"time"
)

const (
	shortLoginDescription = "Login into application"
	longLoginDescription  = "Input your username after that you will be prompted to input your password.\n" +
		"Your token will be saved in the token.txt file, which will be sent with all of your request headers.\n\n" +
		"Example:\n" +
		"login --username \"username\""
	tokenPath = "token.txt"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: shortLoginDescription,
	Long:  longLoginDescription,
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
		fmt.Println()
	},
}

func login(username, password string) error {
	credentials := model.Credentials{
		Username: username,
		Password: password,
	}

	loginURL := clients.BuildURL("core", "v1", "LoginUser")
	tokenResponse := model.TokenResponse{}

	err := utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         loginURL,
		Method:      "POST",
		RequestBody: credentials,
		Response:    &tokenResponse,
		Timeout:     10 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("failed to send login request: %v", err)
	}

	if err := saveTokenToFile(tokenResponse.Token); err != nil {
		return fmt.Errorf("failed to save token: %v", err)
	}

	return nil
}
func saveTokenToFile(token string) error {
	tokenFilePath := tokenPath
	return ioutil.WriteFile(tokenFilePath, []byte(token), 0600)
}

func init() {
	LoginCmd.Flags().StringVarP(&username, flagUsername, shortUsername, "", "Username for login")
	LoginCmd.MarkFlagRequired(flagUsername)
}
