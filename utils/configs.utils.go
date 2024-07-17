package utils

import (
	"fmt"
	"strings"

	"github.com/c12s/cockpit/model"
)

func PrepareConfigDiffRequest(namespace, names, versions, organization string) (interface{}, error) {
	namesList := strings.Split(names, "|")
	versionsList := strings.Split(versions, "|")

	if len(namesList) == 1 {
		namesList = append(namesList, namesList[0])
	}

	if len(versionsList) == 1 {
		versionsList = append(versionsList, versionsList[0])
	}

	if len(namesList) != 2 || len(versionsList) != 2 {
		return nil, fmt.Errorf("invalid names or versions format. Please use 'name1|name2' and 'version1|version2'")
	}

	requestBody := model.SingleConfigDiffRequest{
		Reference: model.SingleConfigReference{
			Name:         namesList[0],
			Organization: organization,
			Namespace:    namespace,
			Version:      versionsList[0],
		},
		Diff: model.SingleConfigReference{
			Name:         namesList[1],
			Organization: organization,
			Namespace:    namespace,
			Version:      versionsList[1],
		},
	}

	return requestBody, nil
}
