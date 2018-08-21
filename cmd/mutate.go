package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/cmd/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func kind(file *model.MutateFile) {
	switch file.Content.Kind {
	case "NodeConfig": //add some configs to all present nodes based on labels in some region/cluster
		fmt.Println("NodeConfig file", file)
	case "NodeAction": // put some action to the all present nodes in some region/cluster like update,restart bash commands etc...
	case "NodeSecret": //add some secrets to all present nodes based on labels in some region/cluster

	}
}

var MutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Mutate state of the constallations",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()

		if _, err := os.Stat(file); err == nil {
			f, err := mutateFile(file)
			if err != nil {
				fmt.Println(err.Error())
			}

			kind(f)
		} else {
			fmt.Println("File not exists")
		}
	},
}

func mutateFile(n ...string) (*model.MutateFile, error) {
	path := ""
	if len(n) > 0 {
		path = n[0]
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var f model.MutateFile
	err = yaml.Unmarshal(yamlFile, &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
