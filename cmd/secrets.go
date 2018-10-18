package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/cmd/helper"
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
		labels := cmd.Flag("labels").Value.String()
		compare := cmd.Flag("compare").Value.String()

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

		callPath := helper.FormCall("secrets", "list", ctx, q)
		// getCall(10*time.Second, callPath)

		fmt.Println(callPath)
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
