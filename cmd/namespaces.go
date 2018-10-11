package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var NamespacesCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "Get the namespaces from the system",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}

var NamespacesGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the namespaces from the system",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
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

var NamespacesMutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Mutate state of the namespaces for the system",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()
		if _, err := os.Stat(file); err == nil {
			f, err := mutateNFile(file)
			if err != nil {
				fmt.Println(err.Error())
			}

			err2, data := namespaces(f)
			if err2 != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(data)
		} else {
			fmt.Println("File not exists")
		}
	},
}
