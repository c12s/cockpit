package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/cmd/helper"
	"github.com/spf13/cobra"
	"os"
	"time"
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
		err, ctx := helper.GetContext()
		if err != nil {
			fmt.Println(err)
			return
		}
		q["user"] = ctx.Context.User

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

		callPath := helper.FormCall("actions", "list", ctx, q)
		err1, resp := helper.GetCall(10*time.Second, callPath)
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		fmt.Println(resp)
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

			err3, ctx := helper.GetContext()
			if err != nil {
				fmt.Println(err3)
				return
			}

			q := map[string]string{}
			callPath := helper.FormCall("actions", "new", ctx, q)
			err4, resp := helper.PostCall(10*time.Second, callPath, data)
			if err4 != nil {
				fmt.Println(err4)
				return
			}
			fmt.Println(resp)

		} else {
			fmt.Println("File not exists")
		}
	},
}
