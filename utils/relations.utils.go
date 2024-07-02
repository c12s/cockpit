package utils

import (
	"fmt"
	"strings"
)

func ParseIDsAndKinds(ids, kinds string) ([]string, []string, error) {
	idsList := strings.Split(ids, "|")
	kindsList := strings.Split(kinds, "|")

	if len(idsList) != 2 || len(kindsList) != 2 {
		return nil, nil, fmt.Errorf("invalid ids or kinds format. Please provide exactly two values separated by '|'")
	}
	return idsList, kindsList, nil
}
