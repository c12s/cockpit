package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
	"github.com/spf13/cobra"
	"os"
)

const (
	apiVersionFlag = "api-version"
)

var magnetar magnetarapi.MagnetarClient

func init() {
	rootCmd.PersistentFlags().String(apiVersionFlag, "1.0.0", "specify c12s API version")

	initClients()
}

func initClients() {
	magnetar = clients.NewMagnetar()
}

var rootCmd = &cobra.Command{
	Use:   "cockpit",
	Short: "Cockpit is a CLI tool for interacting with the c12s system",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
