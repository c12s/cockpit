package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func checkAction(flag string) {
	switch flag {
	case "nodes":
		fmt.Println("List cluster nodes")
	case "configs":
		fmt.Println("List cluster configs [configs at all nodes]")
	case "secrets":
		fmt.Println("List cluster secrets [secrets at all nodes]")
	}
}

var ClusterCmd = &cobra.Command{
	Use: "cluster",
	Run: func(cmd *cobra.Command, args []string) {
		// if Verbose {
		// 	fmt.Println("About to greet friends from Mars...")
		// }
		// lang := cmd.Flag("lang").Value.String()
		// fmt.Printf("%s Mars :)\n", greeting(lang))

		flag := cmd.Flag("list").Value.String()
		checkAction(flag)

		if len(args) > 0 {
			for _, item := range args {
				fmt.Println(item)
			}
		}
	},
}
