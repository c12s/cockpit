package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var ConfigsCmd = &cobra.Command{
	Use:   "configs",
	Short: "Get the configurations from region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}

var ConfigsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the configurations from region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Run: func(cmd *cobra.Command, args []string) {
		labels := cmd.Flag("labels").Value.String()
		compare := cmd.Flag("compare").Value.String()

		q := map[string]string{}
		if labels != "" {
			q["labels"] = labels

		}

		if labels != "" && compare == "" {
			q["compare"] = "any"
		} else if labels != "" && compare != "" {
			q["compare"] = compare
		}

		err, ctx := getContext()
		if err != nil {
			fmt.Println(err)
			return
		}

		callPath := formCall("configs", "list", ctx, q)
		// getCall(10*time.Second, callPath)

		fmt.Println(callPath)
	},
}

var ConfigsMutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Mutate state of the configurations for the region, cluster, node and/or jobs",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()
		if _, err := os.Stat(file); err == nil {
			f, err := mutateFile(file)
			if err != nil {
				fmt.Println(err)
			}

			err2, data := kind(f)
			if err2 != nil {
				fmt.Println(err2)
				return
			}

			err3, _ := getContext()
			if err != nil {
				fmt.Println(err3)
				return
			}
			// callPath := formCall("new", ctx)
			// postCall(10*time.Second, callPath, data)
			fmt.Println(data)
		} else {
			fmt.Println("File not exists")
		}
	},
}
