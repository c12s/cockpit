package cmd

import (
	"errors"
	"fmt"
	"github.com/c12s/cockpit/cmd/helper"
	"github.com/c12s/cockpit/cmd/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func keyNotOK(key string) bool {
	for _, item := range helper.Configs_payloads {
		if key == item {
			return true
		}
	}
	return false
}

func in(payload map[string][]string) error {
	payloads := strings.Join(helper.Configs_payloads, ", ")
	var inside = false

	for _, item := range helper.Configs_payloads {
		if _, ok := payload[item]; ok {
			inside = true
		} else {
			inside = false
		}
	}

	for key, _ := range payload {
		if !keyNotOK(key) {
			return errors.New(fmt.Sprintf("Error: Allowed payloads for Configs are %s", payloads))
		}
	}

	if inside {
		return nil
	}
	return errors.New(fmt.Sprintf("Error: Allowed payloads for Configs %s not presented", payloads))
}

func splitter(entries []string) error {
	for _, entry := range entries {
		if len(strings.Split(entry, "=")) == 1 {
			return errors.New(fmt.Sprintf("Error: entry %s not in the right format KEY=VALUE", entry))
		}
	}
	return nil
}

func files(files []string) error {
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("Error: File %s does not exists", file))
		}
	}
	return nil
}

func validateConfigsPayload(payload map[string][]string) error {
	err := in(payload)
	if err != nil {
		return err
	}

	err = splitter(payload[helper.ENV])
	if err != nil {
		return err
	}

	err = files(payload[helper.FILES])
	if err != nil {
		return err
	}
	return nil
}

func mutateConfigs(file *model.MutateFile) {
	err = validateConfigsPayload(file.Content.Payload)
	if err != nil {
		fmt.Println(err)
	} else {
		data, err := helper.FileToJSON(&file.Content)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("JSON\n", data)
	}
}

func mutateActions(file *model.MutateFile) {

}

func kind(file *model.MutateFile) {
	switch file.Content.Kind {
	case helper.CONFIGS: //add some configs to all present nodes based on labels in some region/cluster
		mutateConfigs(file)

	case helper.ACTIONS: // put some action to the all present nodes in some region/cluster like update,restart bash commands etc...
		mutateActions(file)

	case helper.SECRETS: //add some secrets to all present nodes based on labels in some region/cluster
	case helper.NAMESPACES:
	}
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
