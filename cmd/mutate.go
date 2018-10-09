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

func keyNotOK(key string, test []string) bool {
	for _, item := range test {
		if key == item {
			return true
		}
	}
	return false
}

func in(payload map[string][]string, test []string) error {
	payloads := strings.Join(test, ", ")
	var inside = false

	for _, item := range test {
		if _, ok := payload[item]; ok {
			inside = true
		} else {
			inside = false
		}
	}

	for key, _ := range payload {
		if !keyNotOK(key, test) {
			return errors.New(fmt.Sprintf("Error1: Allowed payloads for this artifact are %s", payloads))
		}
	}

	if inside {
		return nil
	}
	return errors.New(fmt.Sprintf("Error2: Allowed payloads for this artifact %s not presented", payloads))
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
	err := in(payload, helper.Configs_payloads)
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

func validateActionsPayload(payload map[string][]string) error {
	err := in(payload, helper.Actions_payloads)
	if err != nil {
		return err
	}
	return nil
}

func mutateConfigs(file *model.MutateFile) (error, string) {
	err := validateConfigsPayload(file.Content.Payload)
	if err != nil {
		return err, ""
	} else {
		data, err := helper.FileToJSON(&file.Content)
		if err != nil {
			return err, ""
		}
		return nil, data
	}
}

func mutateActions(file *model.MutateFile) (error, string) {
	err := validateActionsPayload(file.Content.Payload)
	if err != nil {
		return err, ""
	} else {
		data, err := helper.FileToJSON(&file.Content)
		if err != nil {
			return err, ""
		}
		return nil, data
	}
}

func mutateSecrets(file *model.MutateFile) (error, string) {
	err := validateConfigsPayload(file.Content.Payload)
	if err != nil {
		return err, ""
	} else {
		data, err := helper.FileToJSON(&file.Content)
		if err != nil {
			return err, ""
		}
		return nil, data
	}
}

func kind(file *model.MutateFile) (error, string) {
	switch file.Content.Kind {
	case helper.CONFIGS: //add some configs to all present nodes based on labels in some region/cluster
		return mutateConfigs(file)

	case helper.ACTIONS: // put some action to the all present nodes in some region/cluster like update,restart bash commands etc...
		return mutateActions(file)

	case helper.SECRETS: //add some secrets to all present nodes based on labels in some region/cluster
		return mutateSecrets(file)

	default:
		return errors.New("Unspupported Kind"), ""
	}
}

func namespaces(file *model.NMutateFile) (error, string) {
	if file.Content.Name == "" {
		return errors.New("Error: Name must be provided, for Namespace artifact!"), ""
	} else {
		data, err := helper.FileToJSON(&file.Content)
		if err != nil {
			return err, ""
		}
		return nil, data
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

func mutateNFile(n ...string) (*model.NMutateFile, error) {
	path := ""
	if len(n) > 0 {
		path = n[0]
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var f model.NMutateFile
	err = yaml.Unmarshal(yamlFile, &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
