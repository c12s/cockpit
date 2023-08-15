package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(putCmd)
}

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("put")
	},
}
