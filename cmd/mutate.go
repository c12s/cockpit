package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/cmd/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var MutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Mutate state of the constallations",
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()

		if _, err := os.Stat(file); err == nil {
			f, err := mutateFile(file)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println(f)
		} else {
			fmt.Println("File not exists")
		}
	},
}

func mutateFile(n ...string) (*model.MutateFile, error) {
	path := "config.yml"
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
