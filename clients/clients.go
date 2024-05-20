package clients

import (
	"fmt"
	"github.com/c12s/cockpit/config"
	"log"
	"os"
)

var cfg *config.Config

func Init() {
	var err error
	cfg, err = config.LoadConfig("../lunar-gateway/config/config.yml")
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
