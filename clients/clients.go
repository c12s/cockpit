package clients

import (
	"fmt"
	"github.com/c12s/cockpit/config"
	"github.com/c12s/cockpit/utils"
	"log"
	"os"
)

var cfg *config.Config

func Init() {
	err := utils.LoadEnvFile(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		fmt.Println("CONFIG_PATH is not set in the .env file")
		os.Exit(1)
	}

	cfg, err = config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Failed to load configuration:", err)
		os.Exit(1)
	}
}

func BuildURL(group, version, action string) string {
	gateway := fmt.Sprintf("http://%s:%s", "localhost", cfg.Gateway.Port)

	methodConfig, ok := cfg.Groups[group][version][action]
	if !ok {
		log.Fatalf("Configuration for %s/%s/%s not found", group, version, action)
	}

	fullMethodRoute := fmt.Sprintf("%s/%s/%s%s", cfg.Gateway.Route, group, version, methodConfig.MethodRoute)

	return fmt.Sprintf("%s%s", gateway, fullMethodRoute)
}
