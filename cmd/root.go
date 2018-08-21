package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.PersistentFlags().StringP("list", "l", "nodes", "list all artifacts inside region, cluster")

	MutateCmd.Flags().StringP("file", "f", "", "mutate region or cluster with new configs or secrets provided in yml file")

	RootCmd.AddCommand(ClusterCmd)
	RootCmd.AddCommand(MutateCmd)
}

var RootCmd = &cobra.Command{
	Use:   "cockpit",
	Short: "This is simple short desc",
	Long:  `This is simple longer desc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Done!")
	},
}
