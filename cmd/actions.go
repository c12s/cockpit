package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var ActionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Get the actions history from region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}

var ActionsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the actinos history from region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Run: func(cmd *cobra.Command, args []string) {
		labels := cmd.Flag("labels").Value.String()
		compare := cmd.Flag("compare").Value.String()
		from := cmd.Flag("from").Value.String()
		to := cmd.Flag("to").Value.String()
		top := cmd.Flag("top").Value.String()

		q := map[string]string{}
		if labels != "" {
			q["labels"] = labels

		}

		if labels != "" && compare == "" {
			q["compare"] = "any"
		} else if labels != "" && compare != "" {
			q["compare"] = compare
		}

		if from != "" {
			q["from"] = from
		}

		if to != "" {
			q["to"] = to
		}

		if top != "" {
			q["top"] = top
		}

		err, ctx := getContext()
		if err != nil {
			fmt.Println(err)
			return
		}

		callPath := formCall("actions", "list", ctx, q)
		// getCall(10*time.Second, callPath)

		fmt.Println(callPath)
	},
}

var ActionsMutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Mutate state of the actios for the region, cluster, node and/or jobs",
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
				fmt.Println(err2)
				return
			}
			fmt.Println(data)
		} else {
			fmt.Println("File not exists")
		}
	},
}
