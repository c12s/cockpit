package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/cmd/get"
	"github.com/c12s/cockpit/cmd/list"
	"github.com/spf13/cobra"
	"os"
)

const (
	apiVersionFlag = "api-version"
)

func init() {
	ListCmd.AddCommand(list.NodesCmd)

	GetCmd.AddCommand(get.NodeCmd)

	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(GetCmd)
	RootCmd.AddCommand(PutCmd)
	RootCmd.AddCommand(DeleteCmd)
	RootCmd.PersistentFlags().String(apiVersionFlag, "1.0.0", "specify c12s API version")
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
}

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
}

var PutCmd = &cobra.Command{
	Use:   "put",
	Short: "Put resources",
}

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
}

var RootCmd = &cobra.Command{
	Use:   "cockpit",
	Short: "Cockpit is a CLI tool for interacting with the c12s system",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
