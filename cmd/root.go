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
}

var RootCmd = &cobra.Command{
	Use:   "cockpit",
	Short: "Get or update state of the regions, clusters, nodes and/or jobs",
	Long:  `This is simple longer desc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}
