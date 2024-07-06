package cmd

import (
	"fmt"
	"os"

	"github.com/c12s/cockpit/aliases"

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
	PutCmd.AddCommand(put.PutAppResourcesCmd)
	PutCmd.AddCommand(put.PutNamespaceResourcesCmd)
	PutStandaloneConfigCmd.AddCommand(put.PutStandaloneConfigCmd)
	PutConfigGroupCmd.AddCommand(put.PutConfigGroupCmd)
	RootCmd.AddCommand(PutCmd)

	// Delete Commands
	DeleteCmd.AddCommand(deleteCmd.DeleteNodeLabelsCmd)
	DeleteCmd.AddCommand(deleteCmd.DeleteSchemaCmd)
	DeleteCmd.AddCommand(DeleteStandaloneConfigCmd)
	DeleteCmd.AddCommand(DeleteConfigCmd)
	DeleteCmd.AddCommand(deleteCmd.DeleteNamespaceCmd)
	DeleteCmd.AddCommand(deleteCmd.DeleteAppCmd)
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
	GetCmd.AddCommand(NodesMetricsCmd)
	GetCmd.AddCommand(get.GetNamespaceHierarchyCmd)
	GetCmd.AddCommand(get.GetNamespaceCmd)
	NodesMetricsCmd.AddCommand(get.LatestMetricsCmd)
	GetStandaloneConfigCmd.AddCommand(get.GetStandaloneConfigCmd)
	GetConfigCmd.AddCommand(get.GetSingleConfigGroupCmd)
	get.GetSchemaCmd.AddCommand(get.GetSchemaVersionCmd)
	RootCmd.AddCommand(GetCmd)
	RootCmd.AddCommand(GetNodesMetricsCmd)

	// Validate Commands
	ValidateCmd.AddCommand(validate.ValidateSchemaVersionCmd)
	RootCmd.AddCommand(ValidateCmd)

	// Create Commands
	CreateCmd.AddCommand(create.CreateSchemaCmd)
	CreateCmd.AddCommand(create.CreateRelationsCmd)
	CreateCmd.AddCommand(create.CreatePoliciesCmd)
	CreateCmd.AddCommand(create.CreateAppCmd)
	CreateCmd.AddCommand(create.CreateNamespaceCmd)
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

	RootCmd.PersistentFlags().String(apiVersionFlag, "1.0.0", "specify c12s API version")
}

var (
	ClaimCmd                      = &cobra.Command{Use: "claim", Short: "Claim resources", Aliases: aliases.ClaimAliases}
	DeleteCmd                     = &cobra.Command{Use: "delete", Short: "Delete resources", Aliases: aliases.DeleteAliases}
	DeleteStandaloneConfigCmd     = &cobra.Command{Use: "standalone", Short: "Delete resources", Aliases: aliases.StandaloneAliases}
	PutCmd                        = &cobra.Command{Use: "put", Short: "Put resources"}
	PutStandaloneConfigCmd        = &cobra.Command{Use: "standalone", Short: "Put resources", Aliases: aliases.StandaloneAliases}
	ListCmd                       = &cobra.Command{Use: "list", Short: "List resources", Aliases: aliases.ListAliases}
	CreateCmd                     = &cobra.Command{Use: "create", Short: "Create resources", Aliases: aliases.CreateAliases}
	GetCmd                        = &cobra.Command{Use: "get", Short: "Get resources", Aliases: aliases.FetchAliases}
	PlaceCmd                      = &cobra.Command{Use: "place", Short: "Place resources", Aliases: aliases.PlaceAliases}
	PutConfigGroupCmd             = &cobra.Command{Use: "config", Short: "Put resources", Aliases: aliases.ConfigAliases}
	PlaceConfigGroupCmd           = &cobra.Command{Use: "config", Short: "Place resources", Aliases: aliases.ConfigAliases}
	PlaceStandaloneConfigGroupCmd = &cobra.Command{Use: "standalone", Short: "Place resources", Aliases: aliases.StandaloneAliases}
	DiffCmd                       = &cobra.Command{Use: "diff", Short: "Diff resources", Aliases: aliases.CompareAliases}
	ListConfigCmd                 = &cobra.Command{Use: "config", Short: "Manipulate with config", Aliases: aliases.ConfigAliases}
	GetConfigCmd                  = &cobra.Command{Use: "config", Short: "Manipulate with config", Aliases: aliases.ConfigAliases}
	DiffConfigCmd                 = &cobra.Command{Use: "config", Short: "Manipulate with config", Aliases: aliases.ConfigAliases}
	DiffStandaloneConfigCmd       = &cobra.Command{Use: "standalone", Short: "Manipulate with config", Aliases: aliases.StandaloneAliases}
	DeleteConfigCmd               = &cobra.Command{Use: "config", Short: "Manipulate with config", Aliases: aliases.ConfigAliases}
	GetStandaloneConfigCmd        = &cobra.Command{Use: "standalone", Short: "Manipulate with config", Aliases: aliases.StandaloneAliases}
	ListStandaloneConfigCmd       = &cobra.Command{Use: "standalone", Short: "Manipulate with config", Aliases: aliases.StandaloneAliases}
	ValidateCmd                   = &cobra.Command{Use: "validate", Short: "Validate resources", Aliases: aliases.ValidateAliases}
	GetNodesMetricsCmd            = &cobra.Command{Use: "get", Short: "Get resources", Aliases: aliases.FetchAliases}
	NodesMetricsCmd               = &cobra.Command{Use: "nodes", Short: "Nodes resources", Aliases: aliases.NodesAliases}

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
