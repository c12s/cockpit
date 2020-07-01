package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/cmd/helper"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var TopologyCmd = &cobra.Command{
	Use:   "topology",
	Short: "Get the topology configurations",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}

var TopologyGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the configurations from region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Run: func(cmd *cobra.Command, args []string) {
		labels := cmd.Flag("labels").Value.String()
		compare := cmd.Flag("compare").Value.String()
		name := cmd.Flag("name").Value.String()

		fmt.Println(labels, compare, name)
	},
}

var TopologyMutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Mutate state of the configurations for the region, cluster, node and/or jobs",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()
		if _, err := os.Stat(file); err == nil {
			f, err := mutateTopology(file)
			if err != nil {
				fmt.Println(err)
			}

			err, data := topology(f)
			if err != nil {
				fmt.Println(err)
				return
			}

			err, ctx := helper.GetContext()
			if err != nil {
				fmt.Println(err)
				return
			}

			q := map[string]string{}
			q["user"] = ctx.Context.User
			q["namespace"] = ctx.Context.Namespace

			h := map[string]string{
				"Content-Type":  "application/json; charset=UTF-8",
				"Authorization": ctx.Context.Token,
			}

			callPath := helper.FormCall("topology", "mutate", ctx, q)
			err, resp := helper.Post(10*time.Second, callPath, data, h)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(resp)
		} else {
			fmt.Println("File not exists")
		}
	},
}
