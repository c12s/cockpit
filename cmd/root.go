package cmd

import (
	"fmt"
	auth "github.com/c12s/cockpit/cmd/auth"
	claim "github.com/c12s/cockpit/cmd/claim"
	delete "github.com/c12s/cockpit/cmd/delete"
	list "github.com/c12s/cockpit/cmd/list"
	put "github.com/c12s/cockpit/cmd/put"
	"github.com/spf13/cobra"
	"os"
)

const (
	apiVersionFlag = "api-version"
)

func init() {
	RootCmd.AddCommand(auth.LoginCmd)
	RootCmd.AddCommand(auth.RegisterCmd)

	ListCmd.AddCommand(list.NodesCmd)
	list.NodesCmd.AddCommand(list.AllocatedNodesCmd)
	RootCmd.AddCommand(ListCmd)

	PutCmd.AddCommand(put.LabelsCmd)
	RootCmd.AddCommand(PutCmd)

	DeleteCmd.AddCommand(delete.DeleteNodeLabelsCmd)
	RootCmd.AddCommand(DeleteCmd)

	ClaimCmd.AddCommand(claim.ClaimNodesCmd)
	RootCmd.AddCommand(ClaimCmd)

	RootCmd.PersistentFlags().String(apiVersionFlag, "1.0.0", "specify c12s API version")
}

var ClaimCmd = &cobra.Command{
	Use:   "claim",
	Short: "Claim resources",
}

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
}

var PutCmd = &cobra.Command{
	Use:   "put",
	Short: "Put resources",
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List nodes",
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
