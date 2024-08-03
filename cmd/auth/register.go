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
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
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
	Aliases: aliases.RegisterAliases,
	Short:   constants.ShortRegisterDesc,
	Long:    constants.LongRegisterDesc,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{constants.UsernameFlag, constants.EmailFlag, constants.NameFlag, constants.OrganizationFlag, constants.SurnameFlag})
	},
	Run: func(cmd *cobra.Command, args []string) {
		password, err := utils.PromptForPassword()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		bar := pb.StartNew(100)
		bar.SetWidth(50)

		err = register(email, name, org, password, surname, username, bar)
		bar.Finish()

		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println("Registration successful!")
	},
}

func register(email, name, org, password, surname, username string, bar *pb.ProgressBar) error {
	registrationDetails := model.RegistrationDetails{
		Email:    email,
		Name:     name,
		Org:      org,
		Password: password,
		Surname:  surname,
		Username: username,
	}

	url := clients.BuildURL("core", "v1", "RegisterUser")

	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			bar.Increment()
		}
	}()

	err := utils.SendHTTPRequestWithProgress(model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		RequestBody: registrationDetails,
		Timeout:     30 * time.Second,
	}, bar)

	if err != nil {
		bar.Finish()
		return err
	}

	bar.SetCurrent(100)
	return nil
}

func init() {
	RegisterCmd.Flags().StringVarP(&email, constants.EmailFlag, constants.EmailShorthandFlag, "", constants.EmailDescription)
	RegisterCmd.Flags().StringVarP(&name, constants.NameFlag, constants.NameShorthandFlag, "", constants.NameDescription)
	RegisterCmd.Flags().StringVarP(&org, constants.OrganizationFlag, constants.OrganizationShorthandFlag, "", constants.OrganizationDescription)
	RegisterCmd.Flags().StringVarP(&surname, constants.SurnameFlag, constants.SurnameShorthandFlag, "", constants.SurnameDescription)
	RegisterCmd.Flags().StringVarP(&username, constants.UsernameFlag, constants.UsernameShorthandFlag, "", constants.UsernameDescription)

	RegisterCmd.MarkFlagRequired(constants.EmailFlag)
	RegisterCmd.MarkFlagRequired(constants.NameFlag)
	RegisterCmd.MarkFlagRequired(constants.OrganizationFlag)
	RegisterCmd.MarkFlagRequired(constants.SurnameFlag)
	RegisterCmd.MarkFlagRequired(constants.UsernameFlag)
}
