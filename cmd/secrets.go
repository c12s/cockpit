package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var SecretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Get the secrets from region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}

var SecretsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the secrets from region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide region, cluster, node and/or job id")
			return
		}
		for _, a := range args {
			fmt.Println(a)
		}
	},
}

var SecretsMutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Mutate secrets of the region, cluster, node and/or jobs",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()

		if _, err := os.Stat(file); err == nil {
			f, err := mutateFile(file)
			if err != nil {
				fmt.Println(err.Error())
			}

			err2, data := kind(f)
			if err2 != nil {
				fmt.Println(data)
				return
			}
		} else {
			fmt.Println("File not exists")
		}
	},
}
