package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	ConfigsMutateCmd.Flags().StringP("file", "f", "", "mutate region, cluster, nod and/or job with new configurations provided in yml file")
	SecretsMutateCmd.Flags().StringP("file", "f", "", "mutate region, cluster, nod and/or job with new secrets provided in yml file")
	ActionsMutateCmd.Flags().StringP("file", "f", "", "mutate region, cluster, nod and/or job with new actions provided in yml file")
	NamespacesMutateCmd.Flags().StringP("file", "f", "", "mutate system with new namespace provided in yml file")
	RolesMutateCmd.Flags().StringP("file", "f", "", "mutate system with new roles provided in yml file")
	RegisterCmd.Flags().StringP("file", "f", "", "mutate system with new user provided in yml file")
	TopologyMutateCmd.Flags().StringP("file", "f", "", "mutate system with new user provided in yml file")

	LoginCmd.Flags().StringP("username", "u", "", "provide username to login to system")
	LoginCmd.Flags().StringP("password", "p", "", "provide password to login to system")
	InitCmd.Flags().StringP("address", "a", "", "provide service ip address, so CLI can comunicate with rest of the system")
	InitCmd.Flags().StringP("version", "v", "", "provide service version, so CLI can comunicate with rest of the system [default v1]")
	SwitchCmd.Flags().StringP("namespace", "n", "", "switch namespace")
	SwitchCmd.Flags().BoolP("context", "c", false, "switch context")

	ConfigsGetCmd.Flags().StringP("labels", "l", "", "list of key-value pairs for configs selection. [k1:v1,k2:v2,...]")
	ConfigsGetCmd.Flags().StringP("compare", "c", "", "compare rule, when selecting configs [any | all]")

	ActionsGetCmd.Flags().StringP("labels", "l", "", "list of key-value pairs for actions selection. [k1:v1,k2:v2,...]")
	ActionsGetCmd.Flags().StringP("compare", "c", "", "compare rule, when selecting actions [any | all]")
	ActionsGetCmd.Flags().StringP("from", "f", "", "timestamp filtering, from where to start lookup")
	ActionsGetCmd.Flags().StringP("to", "t", "", "timestamp filtering to where to end lookup")
	ActionsGetCmd.Flags().StringP("head", "e", "", "returning result contains top n elements")
	ActionsGetCmd.Flags().StringP("tail", "a", "", "returning result contains last n elements")

	SecretsGetCmd.Flags().StringP("labels", "l", "", "list of key-value pairs for secrets selection. [k1:v1,k2:v2,...]")
	SecretsGetCmd.Flags().StringP("compare", "c", "", "compare rule, when selecting secrets [any | all]")

	NamespacesGetCmd.Flags().StringP("labels", "l", "", "list of key-value pairs for namespaces selection. [k1:v1,k2:v2,...]")
	NamespacesGetCmd.Flags().StringP("compare", "c", "", "compare rule, when selecting namespaces [any | all]")
	NamespacesGetCmd.Flags().StringP("name", "n", "", "name, when selecting namespaces")
	RolesGetCmd.Flags().StringP("users", "u", "", "list of usernames separated with comma. [john,doe,...]")
	RolesGetCmd.Flags().StringP("resources", "r", "", "list of resources separated with comma. [configs,actions,...]")
	RolesGetCmd.Flags().StringP("namespaces", "n", "", "list of namespaces separated with comma. [configs,actions,...]")

	TraceListCmd.Flags().StringP("tags", "a", "", "list of key-value pairs for tags selection. [k1:v1,k2:v2,...]")
	TraceGetCmd.Flags().StringP("task", "t", "", "trace id to get complate trace")
	TopologyGetCmd.Flags().StringP("labels", "l", "", "list of key-value pairs for topology selection. [k1:v1,k2:v2,...]")
	TopologyGetCmd.Flags().StringP("compare", "c", "", "compare rule, when selecting topology [any | all]")
	TopologyGetCmd.Flags().StringP("name", "n", "", "name, when selecting topology")

	NodesGetCmd.Flags().StringP("labels", "l", "", "list of key-value pairs for topology selection. [k1:v1,k2:v2,...]")

	ConfigsCmd.AddCommand(ConfigsGetCmd)
	ConfigsCmd.AddCommand(ConfigsMutateCmd)
	RootCmd.AddCommand(ConfigsCmd)

	SecretsCmd.AddCommand(SecretsGetCmd)
	SecretsCmd.AddCommand(SecretsMutateCmd)
	RootCmd.AddCommand(SecretsCmd)

	ActionsCmd.AddCommand(ActionsGetCmd)
	ActionsCmd.AddCommand(ActionsMutateCmd)
	RootCmd.AddCommand(ActionsCmd)

	NamespacesCmd.AddCommand(NamespacesGetCmd)
	NamespacesCmd.AddCommand(NamespacesMutateCmd)
	RootCmd.AddCommand(NamespacesCmd)

	TraceCmd.AddCommand(TraceGetCmd)
	TraceCmd.AddCommand(TraceListCmd)
	RootCmd.AddCommand(TraceCmd)

	ContextCmd.AddCommand(InitCmd)
	ContextCmd.AddCommand(LoginCmd)
	ContextCmd.AddCommand(LogoutCmd)
	ContextCmd.AddCommand(DropCmd)
	ContextCmd.AddCommand(RegisterCmd)
	ContextCmd.AddCommand(SwitchCmd)
	RootCmd.AddCommand(ContextCmd)

	RolesCmd.AddCommand(RolesGetCmd)
	RolesCmd.AddCommand(RolesMutateCmd)
	RootCmd.AddCommand(RolesCmd)

	TopologyCmd.AddCommand(TopologyGetCmd)
	TopologyCmd.AddCommand(TopologyMutateCmd)
	RootCmd.AddCommand(TopologyCmd)

	NodesCmd.AddCommand(NodesGetCmd)
	RootCmd.AddCommand(NodesCmd)
}

var RootCmd = &cobra.Command{
	Use:   "cockpit",
	Short: "Get or update state of the regions, clusters, nodes and/or jobs",
	Long:  `This is simple longer desc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}
