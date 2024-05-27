package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const (
	shortRegisterDescription = "Register a new user"
	longRegisterDescription  = "Register a new user by providing an email, name, organization, surname, and username. \n" +
		"Once these details are entered, you will be prompted to input your password.\n\n" +
		"Example:\n" +
		"register --email \"example@gmail.com\" --name \"name\" --org \"org\" --surname \"surname\" --username \"username\""

	// Flag Constants
	flagEmail    = "email"
	flagName     = "name"
	flagOrg      = "org"
	flagSurname  = "surname"
	flagUsername = "username"

	// Flag Shorthand Constants
	shortEmail    = "e"
	shortName     = "n"
	shortOrg      = "r"
	shortSurname  = "s"
	shortUsername = "u"

	// Flag Descriptions
	emailDesc    = "Email for registration"
	nameDesc     = "Name for registration"
	orgDesc      = "Organization for registration"
	surnameDesc  = "Surname for registration"
	usernameDesc = "Username for registration"
)

var (
	email    string
	name     string
	org      string
	surname  string
	username string
)

var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: shortRegisterDescription,
	Long:  longRegisterDescription,
	Run: func(cmd *cobra.Command, args []string) {
		password, err := utils.PromptForPassword()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		err = register(email, name, org, password, surname, username)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Registration successful!")
		fmt.Println()
	},
}

func register(email, name, org, password, surname, username string) error {
	registrationDetails := model.RegistrationDetails{
		Email:    email,
		Name:     name,
		Org:      org,
		Password: password,
		Surname:  surname,
		Username: username,
	}

	url := clients.BuildURL("core", "v1", "RegisterUser")

	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		RequestBody: registrationDetails,
		Timeout:     10 * time.Second,
	})
}

func init() {
	RegisterCmd.Flags().StringVarP(&email, flagEmail, shortEmail, "", emailDesc)
	RegisterCmd.Flags().StringVarP(&name, flagName, shortName, "", nameDesc)
	RegisterCmd.Flags().StringVarP(&org, flagOrg, shortOrg, "", orgDesc)
	RegisterCmd.Flags().StringVarP(&surname, flagSurname, shortSurname, "", surnameDesc)
	RegisterCmd.Flags().StringVarP(&username, flagUsername, shortUsername, "", usernameDesc)

	RegisterCmd.MarkFlagRequired(flagEmail)
	RegisterCmd.MarkFlagRequired(flagName)
	RegisterCmd.MarkFlagRequired(flagOrg)
	RegisterCmd.MarkFlagRequired(flagSurname)
	RegisterCmd.MarkFlagRequired(flagUsername)
}
