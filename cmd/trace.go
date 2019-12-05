package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/cmd/helper"
	"github.com/spf13/cobra"
	"time"
)

var TraceCmd = &cobra.Command{
	Use:   "trace",
	Short: "Get the state and trace execution of user defined job/s",
	Long:  "Show complate trace of user defined job/s",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}

var TraceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the job trace history by traceId",
	Long:  "Show complate trace of user defined job",
	Run: func(cmd *cobra.Command, args []string) {
		task := cmd.Flag("task").Value.String()
		err, ctx := helper.GetContext()
		if err != nil {
			fmt.Println(err)
			return
		}

		q := map[string]string{}
		q["user"] = ctx.Context.User
		q["traceId"] = task

		h := map[string]string{
			"Content-Type":  "application/json; charset=UTF-8",
			"Authorization": ctx.Context.Token,
		}

		callPath := helper.FormCall("trace", "get", ctx, q)
		err, resp := helper.Get(10*time.Second, callPath, h)
		if err != nil {
			fmt.Println(err)
			return
		}
		helper.Print("trace/get", resp)
	},
}

var TraceListCmd = &cobra.Command{
	Use:   "list",
	Short: "Query the job trace history by tags",
	Long:  "Show complate trace of user defined jobs by query",
	Run: func(cmd *cobra.Command, args []string) {
		tags := cmd.Flag("tags").Value.String()
		err, ctx := helper.GetContext()
		if err != nil {
			fmt.Println(err)
			return
		}

		q := map[string]string{}
		q["user"] = ctx.Context.User
		q["tags"] = tags

		h := map[string]string{
			"Content-Type":   "application/json; charset=UTF-8",
			"Authorization:": ctx.Context.Token,
		}

		callPath := helper.FormCall("trace", "list", ctx, q)
		err, resp := helper.Get(10*time.Second, callPath, h)
		if err != nil {
			fmt.Println(err)
			return
		}
		helper.Print("trace/list", resp)
	},
}
