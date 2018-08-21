package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// cockpit region list regionId -> list all clusters inside region if user
// have right to see this data.
var RegionCmd = &cobra.Command{
	Use:   "region",
	Short: "Show clusters inside region",
	Long:  "Show all avalible clusters inside region provided by user. If user is logged in and have access to see that data, it will be presented to him.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flag := cmd.Flag("list").Value.String()

		regionId := args[0]
		fmt.Println(regionId, flag)
	},
}
