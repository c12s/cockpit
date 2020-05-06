package cmd

import (
	// "encoding/json"
	"fmt"
	"github.com/c12s/cockpit/cmd/helper"
	"github.com/spf13/cobra"
	"os"
	"time"
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
	Use:   "list",
	Short: "Get the namespaces from the system",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		labels := cmd.Flag("labels").Value.String()
		compare := cmd.Flag("compare").Value.String()
		name := cmd.Flag("name").Value.String()

		q := map[string]string{}
		err, ctx := helper.GetContext()
		if err != nil {
			fmt.Println(err)
			return
		}
		q["user"] = ctx.Context.User
		q["namespace"] = ctx.Context.Namespace

		if labels != "" {
			q["labels"] = labels

		}

		if labels != "" && compare == "" {
			q["compare"] = "any"
		} else if labels != "" && compare != "" {
			q["compare"] = compare
		}

		if name != "" {
			q["name"] = name

		}

		h := map[string]string{
			"Content-Type":  "application/json; charset=UTF-8",
			"Authorization": ctx.Context.Token,
		}

		callPath := helper.FormCall("namespaces", "list", ctx, q)
		err1, resp := helper.Get(10*time.Second, callPath, h)
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		helper.Print("namespaces", resp)
	},
}

type Rez struct {
	Rez map[string]map[string]string `json:"rez"`
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
				return
			}

			err2, data := namespaces(f)
			if err2 != nil {
				fmt.Println(err)
				return
			}

			err3, ctx := helper.GetContext()
			if err != nil {
				fmt.Println(err3)
				return
			}

			q := map[string]string{}
			q["user"] = ctx.Context.User
			q["namespace"] = ctx.Context.Namespace

			h := map[string]string{
				"Content-Type":  "application/json; charset=UTF-8",
				"Authorization": ctx.Context.Token,
			}

			callPath := helper.FormCall("namespaces", "mutate", ctx, q)
			err4, resp := helper.Post(10*time.Second, callPath, data, h)
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
