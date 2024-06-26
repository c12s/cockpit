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
	shortRegisterDescription = "Register a new user"
	longRegisterDescription  = `Register a new user by providing an email, name, organization, surname, and username. 
Once these details are entered, you will be prompted to input your password.

Example:
- cockpit register --email "example@gmail.com" --name "name" --org "org" --surname "surname" --username "username"`

	// Flag Constants
	emailFlag    = "email"
	nameFlag     = "name"
	orgFlag      = "org"
	surnameFlag  = "surname"
	usernameFlag = "username"

	// Flag Shorthand Constants
	emailShorthandFlag        = "e"
	nameShorthandFlag         = "n"
	organizationShorthandFlag = "r"
	surnameShorthandFlag      = "s"
	usernameShorthandFlag     = "u"

	// Flag Descriptions
	emailDescription    = "Email for registration"
	nameDescription     = "Name for registration"
	orgDescription      = "Organization for registration"
	surnameDescription  = "Surname for registration"
	usernameDescription = "Username for registration"
)

var (
	email    string
	name     string
	org      string
	surname  string
	username string
)

var RegisterCmd = &cobra.Command{
	Use:     "register",
	Short:   shortRegisterDescription,
	Aliases: aliases.RegisterAliases,
	Long:    longRegisterDescription,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{usernameFlag, emailFlag, nameFlag, orgFlag, surnameFlag})
	},
	Run: func(cmd *cobra.Command, args []string) {
		password, err := utils.PromptForPassword()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		err = register(email, name, org, password, surname, username)
		if err != nil {
			fmt.Println("Error:", err)
			println()
			os.Exit(1)
		}

		fmt.Println("Registration successful!\n")
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
	RegisterCmd.Flags().StringVarP(&email, emailFlag, emailShorthandFlag, "", emailDescription)
	RegisterCmd.Flags().StringVarP(&name, nameFlag, nameShorthandFlag, "", nameDescription)
	RegisterCmd.Flags().StringVarP(&org, orgFlag, organizationShorthandFlag, "", orgDescription)
	RegisterCmd.Flags().StringVarP(&surname, surnameFlag, surnameShorthandFlag, "", surnameDescription)
	RegisterCmd.Flags().StringVarP(&username, usernameFlag, usernameShorthandFlag, "", usernameDescription)

	RegisterCmd.MarkFlagRequired(emailFlag)
	RegisterCmd.MarkFlagRequired(nameFlag)
	RegisterCmd.MarkFlagRequired(orgFlag)
	RegisterCmd.MarkFlagRequired(surnameFlag)
	RegisterCmd.MarkFlagRequired(usernameFlag)
}
