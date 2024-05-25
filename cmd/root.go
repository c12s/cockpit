package cmd

import (
	"fmt"
	auth "github.com/c12s/cockpit/cmd/auth"
	claim "github.com/c12s/cockpit/cmd/claim"
	create "github.com/c12s/cockpit/cmd/create"
	delete "github.com/c12s/cockpit/cmd/delete"
	diff "github.com/c12s/cockpit/cmd/diff"
	get "github.com/c12s/cockpit/cmd/get"
	list "github.com/c12s/cockpit/cmd/list"
	place "github.com/c12s/cockpit/cmd/place"
	put "github.com/c12s/cockpit/cmd/put"
	validate "github.com/c12s/cockpit/cmd/validate"
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
	ListCmd.AddCommand(ListConfigCmd)
	ListConfigCmd.AddCommand(list.ListConfigGroupCmd)
	list.ListConfigGroupCmd.AddCommand(list.ListConfigGroupPlacementsCmd)
	list.NodesCmd.AddCommand(list.AllocatedNodesCmd)
	RootCmd.AddCommand(ListCmd)

	PutCmd.AddCommand(put.LabelsCmd)
	PutCmd.AddCommand(PutConfigGroupCmd)
	PutConfigGroupCmd.AddCommand(put.PutConfigGroupCmd)
	RootCmd.AddCommand(PutCmd)

	DeleteCmd.AddCommand(delete.DeleteNodeLabelsCmd)
	DeleteCmd.AddCommand(delete.DeleteSchemaCmd)
	DeleteCmd.AddCommand(DeleteConfigCmd)
	DeleteConfigCmd.AddCommand(delete.DeleteConfigGroupCmd)
	RootCmd.AddCommand(DeleteCmd)

	ClaimCmd.AddCommand(claim.ClaimNodesCmd)
	RootCmd.AddCommand(ClaimCmd)

	GetCmd.AddCommand(get.GetSchemaCmd)
	GetCmd.AddCommand(GetConfigCmd)
	GetCmd.AddCommand(GetStandaloneConfigCmd)
	GetStandaloneConfigCmd.AddCommand(get.GetStandaloneConfigCmd)
	GetConfigCmd.AddCommand(get.GetSingleConfigGroupCmd)
	get.GetSchemaCmd.AddCommand(get.GetSchemaVersionCmd)
	RootCmd.AddCommand(GetCmd)

	ValidateCmd.AddCommand(validate.ValidateSchemaVersionCmd)
	RootCmd.AddCommand(ValidateCmd)

	CreateCmd.AddCommand(create.CreateSchemaCmd)
	RootCmd.AddCommand(CreateCmd)

	DiffCmd.AddCommand(DiffConfigCmd)
	DiffConfigCmd.AddCommand(diff.DiffConfigGroupCmd)
	RootCmd.AddCommand(DiffCmd)

	PlaceCmd.AddCommand(PlaceConfigGroupCmd)
	PlaceConfigGroupCmd.AddCommand(place.PlaceConfigGroupPlacementsCmd)
	RootCmd.AddCommand(PlaceCmd)

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
	Short: "List resources",
}

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources",
}

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
}

var PlaceCmd = &cobra.Command{
	Use:   "place",
	Short: "Place resources",
}

var PutConfigGroupCmd = &cobra.Command{
	Use:   "config",
	Short: "Put resources",
}

var PlaceConfigGroupCmd = &cobra.Command{
	Use:   "config",
	Short: "Place resources",
}

var DiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Diff resources",
}

var ListConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manipulate with config",
}

var GetConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manipulate with config",
}

var DiffConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manipulate with config",
}

var DeleteConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manipulate with config",
}

var GetStandaloneConfigCmd = &cobra.Command{
	Use:   "standalone",
	Short: "Manipulate with config",
}

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Get resources",
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
