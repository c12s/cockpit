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

func in(payload map[string]yaml.MapSlice, test []string) error {
	payloads := strings.Join(test, ", ")
	var inside = false

	for _, item := range test {
		if _, ok := payload[item]; ok {
			inside = true
		}
		// else {
		// 	inside = false
		// }
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

func files(files yaml.MapSlice) error {
	for _, sItem := range files {
		switch file := sItem.Value.(type) {
		case string:
			if _, err := os.Stat(file); os.IsNotExist(err) {
				return errors.New(fmt.Sprintf("Error: File %s does not exists", file))
			}
		}
	}
	return nil
}

func validateConfigsPayload(payload map[string]yaml.MapSlice) error {
	err := in(payload, helper.Configs_payloads)
	if err != nil {
		return err
	}

	err = files(payload[helper.FILES])
	if err != nil {
		return err
	}
	return nil
}

func validateActionsPayload(payload map[string]yaml.MapSlice) error {
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
		fmt.Println(err)
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
	if file.Content.MTData.Namespace == "" {
		return errors.New("Error: Name must be provided, for Namespace artifact!"), ""
	} else {
		data, err := helper.FileToJSON(&file.Content)
		if err != nil {
			return err, ""
		}
		return nil, data
	}
}

func users(file *model.NMutateFile) (error, string) {
	if len(file.Content.Payload) == 0 {
		return errors.New("Error: Username and password must be provided"), ""
	} else {
		data, err := helper.FileToJSON(&file.Content)
		if err != nil {
			return err, ""
		}
		return nil, data
	}
}

func roles(file *model.RolesFile) (error, string) {
	if file.Content.Payload.User == "" {
		return errors.New("Error: Username must be provided, for Role artifact!"), ""
	} else {
		data, err := helper.FileToJSON(&file.Content)
		if err != nil {
			return err, ""
		}
		return nil, data
	}
}

func topology(file *model.TopologyFile) (error, string) {
	if file.Content.Payload.Name == "" {
		return errors.New("Error: Name must be provided, for Topology artifact!"), ""
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

func mutateTopology(n ...string) (*model.TopologyFile, error) {
	path := ""
	if len(n) > 0 {
		path = n[0]
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var f model.TopologyFile
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

func mutateRolesFile(n ...string) (*model.RolesFile, error) {
	path := ""
	if len(n) > 0 {
		path = n[0]
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var f model.RolesFile
	err = yaml.Unmarshal(yamlFile, &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
