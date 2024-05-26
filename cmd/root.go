package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	auth "github.com/c12s/cockpit/cmd/auth"
	claim "github.com/c12s/cockpit/cmd/claim"
	create "github.com/c12s/cockpit/cmd/create"
	deleteCmd "github.com/c12s/cockpit/cmd/delete"
	diff "github.com/c12s/cockpit/cmd/diff"
	get "github.com/c12s/cockpit/cmd/get"
	list "github.com/c12s/cockpit/cmd/list"
	place "github.com/c12s/cockpit/cmd/place"
	put "github.com/c12s/cockpit/cmd/put"
	validate "github.com/c12s/cockpit/cmd/validate"
)

const (
	apiVersionFlag = "api-version"
)

func init() {
	// Authentication Commands
	RootCmd.AddCommand(auth.LoginCmd)
	RootCmd.AddCommand(auth.RegisterCmd)

	// List Commands
	ListCmd.AddCommand(list.NodesCmd)
	ListCmd.AddCommand(ListConfigCmd)
	ListCmd.AddCommand(ListStandaloneConfigCmd)
	ListStandaloneConfigCmd.AddCommand(list.ListStandaloneConfigCmd)
	ListConfigCmd.AddCommand(list.ListConfigGroupCmd)
	list.ListStandaloneConfigCmd.AddCommand(list.ListStandaloneConfigPlacementsCmd)
	list.ListConfigGroupCmd.AddCommand(list.ListConfigGroupPlacementsCmd)
	list.NodesCmd.AddCommand(list.AllocatedNodesCmd)
	RootCmd.AddCommand(ListCmd)

	// Put Commands
	PutCmd.AddCommand(put.LabelsCmd)
	PutCmd.AddCommand(PutConfigGroupCmd)
	PutCmd.AddCommand(PutStandaloneConfigCmd)
	PutStandaloneConfigCmd.AddCommand(put.PutStandaloneConfigCmd)
	PutConfigGroupCmd.AddCommand(put.PutConfigGroupCmd)
	RootCmd.AddCommand(PutCmd)

	// Delete Commands
	DeleteCmd.AddCommand(deleteCmd.DeleteNodeLabelsCmd)
	DeleteCmd.AddCommand(deleteCmd.DeleteSchemaCmd)
	DeleteCmd.AddCommand(DeleteStandaloneConfigCmd)
	DeleteCmd.AddCommand(DeleteConfigCmd)
	DeleteStandaloneConfigCmd.AddCommand(deleteCmd.DeleteStandaloneConfigCmd)
	DeleteConfigCmd.AddCommand(deleteCmd.DeleteConfigGroupCmd)
	RootCmd.AddCommand(DeleteCmd)

	// Claim Commands
	ClaimCmd.AddCommand(claim.ClaimNodesCmd)
	RootCmd.AddCommand(ClaimCmd)

	// Get Commands
	GetCmd.AddCommand(get.GetSchemaCmd)
	GetCmd.AddCommand(GetConfigCmd)
	GetCmd.AddCommand(GetStandaloneConfigCmd)
	GetStandaloneConfigCmd.AddCommand(get.GetStandaloneConfigCmd)
	GetConfigCmd.AddCommand(get.GetSingleConfigGroupCmd)
	get.GetSchemaCmd.AddCommand(get.GetSchemaVersionCmd)
	RootCmd.AddCommand(GetCmd)

	// Validate Commands
	ValidateCmd.AddCommand(validate.ValidateSchemaVersionCmd)
	RootCmd.AddCommand(ValidateCmd)

	// Create Commands
	CreateCmd.AddCommand(create.CreateSchemaCmd)
	RootCmd.AddCommand(CreateCmd)

	// Diff Commands
	DiffCmd.AddCommand(DiffConfigCmd)
	DiffCmd.AddCommand(DiffStandaloneConfigCmd)
	DiffStandaloneConfigCmd.AddCommand(diff.DiffStandaloneConfigCmd)
	DiffConfigCmd.AddCommand(diff.DiffConfigGroupCmd)
	RootCmd.AddCommand(DiffCmd)

	// Place Commands
	PlaceCmd.AddCommand(PlaceConfigGroupCmd)
	PlaceCmd.AddCommand(PlaceStandaloneConfigGroupCmd)
	PlaceStandaloneConfigGroupCmd.AddCommand(place.PlaceStandaloneConfigPlacementsCmd)
	PlaceConfigGroupCmd.AddCommand(place.PlaceConfigGroupPlacementsCmd)
	RootCmd.AddCommand(PlaceCmd)

	// Persistent Flags
	RootCmd.PersistentFlags().String(apiVersionFlag, "1.0.0", "specify c12s API version")
}

var (
	ClaimCmd                      = &cobra.Command{Use: "claim", Short: "Claim resources"}
	DeleteCmd                     = &cobra.Command{Use: "delete", Short: "Delete resources"}
	DeleteStandaloneConfigCmd     = &cobra.Command{Use: "standalone", Short: "Delete resources"}
	PutCmd                        = &cobra.Command{Use: "put", Short: "Put resources"}
	PutStandaloneConfigCmd        = &cobra.Command{Use: "standalone", Short: "Put resources"}
	ListCmd                       = &cobra.Command{Use: "list", Short: "List resources"}
	CreateCmd                     = &cobra.Command{Use: "create", Short: "Create resources"}
	GetCmd                        = &cobra.Command{Use: "get", Short: "Get resources"}
	PlaceCmd                      = &cobra.Command{Use: "place", Short: "Place resources"}
	PutConfigGroupCmd             = &cobra.Command{Use: "config", Short: "Put resources"}
	PlaceConfigGroupCmd           = &cobra.Command{Use: "config", Short: "Place resources"}
	PlaceStandaloneConfigGroupCmd = &cobra.Command{Use: "standalone", Short: "Place resources"}
	DiffCmd                       = &cobra.Command{Use: "diff", Short: "Diff resources"}
	ListConfigCmd                 = &cobra.Command{Use: "config", Short: "Manipulate with config"}
	GetConfigCmd                  = &cobra.Command{Use: "config", Short: "Manipulate with config"}
	DiffConfigCmd                 = &cobra.Command{Use: "config", Short: "Manipulate with config"}
	DiffStandaloneConfigCmd       = &cobra.Command{Use: "standalone", Short: "Manipulate with config"}
	DeleteConfigCmd               = &cobra.Command{Use: "config", Short: "Manipulate with config"}
	GetStandaloneConfigCmd        = &cobra.Command{Use: "standalone", Short: "Manipulate with config"}
	ListStandaloneConfigCmd       = &cobra.Command{Use: "standalone", Short: "Manipulate with config"}
	ValidateCmd                   = &cobra.Command{Use: "validate", Short: "Get resources"}

	// Root Command
	RootCmd = &cobra.Command{
		Use:   "cockpit",
		Short: "Cockpit is a CLI tool for interacting with the c12s system",
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
