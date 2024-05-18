package utils

import (
	"fmt"
	"golang.org/x/term"
	"io/ioutil"
	"syscall"
)

const (
	tokenFilePath = "token.txt"
)

func PromptForPassword() (string, error) {
	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %v", err)
	}
	fmt.Println()
	return string(passwordBytes), nil
}

func ReadTokenFromFile() (string, error) {
	token, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read token file: %w", err)
	}
	return string(token), nil
}
