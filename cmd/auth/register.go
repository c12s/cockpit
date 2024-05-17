package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	shortRegisterDescription = "Register a new user"
	longRegisterDescription  = "Register a new user by providing an email, name, organization, surname, and username. \n" +
		"Once these details are entered, you will be prompted to input your password."
	flagEmail     = "email"
	flagName      = "name"
	flagOrg       = "org"
	flagSurname   = "surname"
	flagUsername  = "username"
	shortEmail    = "e"
	shortName     = "n"
	shortOrg      = "o"
	shortSurname  = "s"
	shortUsername = "u"
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
		password, err := PromptForPassword()
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
	},
}

func init() {
	RegisterCmd.Flags().StringVarP(&email, flagEmail, shortEmail, "", "Email for registration")
	RegisterCmd.Flags().StringVarP(&name, flagName, shortName, "", "Name for registration")
	RegisterCmd.Flags().StringVarP(&org, flagOrg, shortOrg, "", "Organization for registration")
	RegisterCmd.Flags().StringVarP(&surname, flagSurname, shortSurname, "", "Surname for registration")
	RegisterCmd.Flags().StringVarP(&username, flagUsername, shortUsername, "", "Username for registration")

	RegisterCmd.MarkFlagRequired(flagEmail)
	RegisterCmd.MarkFlagRequired(flagName)
	RegisterCmd.MarkFlagRequired(flagOrg)
	RegisterCmd.MarkFlagRequired(flagSurname)
	RegisterCmd.MarkFlagRequired(flagUsername)
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

	registrationJSON, err := json.Marshal(registrationDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal registration details: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := clients.Clients.Gateway + "/apis/core/v1/users"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(registrationJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	bodyString := string(bodyBytes)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registration failed: %s", bodyString)
	}

	return nil
}
