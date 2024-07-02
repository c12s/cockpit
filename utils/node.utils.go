package utils

import (
	"fmt"
	"github.com/c12s/cockpit/model"
	"strings"
)

func CreateNodeQuery(query string) ([]model.NodeQuery, error) {
	if query == "" {
		return nil, nil
	}

	parts := strings.Fields(query)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid query format. Please use 'key operation value'")
	}

	labelKey := parts[0]
	shouldBe := parts[1]
	value := parts[2]

	nodeQuery := model.NodeQuery{
		LabelKey: labelKey,
		ShouldBe: shouldBe,
		Value:    value,
	}

	return []model.NodeQuery{nodeQuery}, nil
}
