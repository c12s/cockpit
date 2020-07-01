package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/cmd/helper"
	"github.com/spf13/cobra"
	"time"
)

var NodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Query the list of avalible nodes.",
	Long:  "Query avalible nodes, so that we can create new regions and clusters",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}

var NodesGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Query the list of avalible nodes.",
	Long:  "Query avalible nodes, so that we can create new regions and clusters. By adding k:v pairs we can filter nodes by some specific tags like memory, cpu, disk, storage etc",
	Run: func(cmd *cobra.Command, args []string) {
		labels := cmd.Flag("labels").Value.String()
		err, ctx := helper.GetContext()
		if err != nil {
			fmt.Println(err)
			return
		}
		pLabels, err := helper.PrepareLabels(labels)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		q := map[string]string{}
		q["user"] = ctx.Context.User
		q["namespace"] = ctx.Context.Namespace
		q["labels"] = pLabels

		h := map[string]string{
			"Content-Type":  "application/json; charset=UTF-8",
			"Authorization": ctx.Context.Token,
		}

		callPath := helper.FormCall("nodes", "list", ctx, q)
		err, resp := helper.Get(10*time.Second, callPath, h)
		if err != nil {
			fmt.Println(err)
			return
		}
		helper.Print("nodes", resp)
	},
}
