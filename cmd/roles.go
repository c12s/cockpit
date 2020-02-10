package cmd

import (
	// "encoding/json"
	"fmt"
	"github.com/c12s/cockpit/cmd/helper"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var RolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Get the roles from the system",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}
var RolesGetCmd = &cobra.Command{
	Use:   "list",
	Short: "Get the roles from the system",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		users := cmd.Flag("users").Value.String()
		resources := cmd.Flag("resources").Value.String()
		namespaces := cmd.Flag("namespaces").Value.String()

		q := map[string]string{}
		err, ctx := helper.GetContext()
		if err != nil {
			fmt.Println(err)
			return
		}
		q["user"] = ctx.Context.User
		q["users"] = users
		q["resources"] = resources
		q["namespaces"] = namespaces

		h := map[string]string{
			"Content-Type":  "application/json; charset=UTF-8",
			"Authorization": ctx.Context.Token,
		}

		callPath := helper.FormCall("roles", "list", ctx, q)
		err1, resp := helper.Get(10*time.Second, callPath, h)
		if err1 != nil {
			fmt.Println("ERROR: ", err1)
			return
		}
		helper.Print("roles", resp)
	},
}
var RolesMutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Mutate state of the roles for the system",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()
		if _, err := os.Stat(file); err == nil {
			f, err := mutateRolesFile(file)
			if err != nil {
				fmt.Println(err.Error())
			}

			err2, data := roles(f)
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

			h := map[string]string{
				"Content-Type":  "application/json; charset=UTF-8",
				"Authorization": ctx.Context.Token,
			}

			callPath := helper.FormCall("roles", "mutate", ctx, q)
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
