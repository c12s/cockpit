package cmd

import (
	"fmt"
	"strings"

	"github.com/c12s/cockpit/model"
	"gopkg.in/yaml.v3"
)

// specOrg reads only metadata.org, leniently: the org selects the schema in
// Quasar, which then reports every other problem in the document.
func specOrg(doc []byte) string {
	var head struct {
		Metadata struct {
			Org string `yaml:"org"`
		} `yaml:"metadata"`
	}
	_ = yaml.Unmarshal(doc, &head)
	return head.Metadata.Org
}

// validateConfigurationResponse is Quasar's ValidateConfiguration reply as
// the gateway returns it (proto JSON: false and empty fields are omitted).
type validateConfigurationResponse struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
	IsValid bool   `json:"isValid"`
}

// validateSpecWithQuasar sends the spec, as written, to Quasar, which
// validates it against the org's host-image-spec schema. The shape of the
// file (fields, types, formats) is defined there, once, for every client.
func validateSpecWithQuasar(org string, doc []byte) error {
	if org == "" {
		return fmt.Errorf("validating spec: metadata.org is required")
	}
	body := struct {
		SchemaDetails model.SchemaDetails `json:"schema_details"`
		Configuration string              `json:"configuration"`
	}{
		SchemaDetails: model.SchemaDetails{
			Organization: org,
			Namespace:    applySchemaNamespace,
			SchemaName:   hostImageSchemaName,
			Version:      applySchemaVersion,
		},
		Configuration: string(doc),
	}
	var resp validateConfigurationResponse
	if err := gatewaySend("ValidateConfiguration", "GET", body, &resp); err != nil {
		// Quasar checks access to the schema before looking it up, so a
		// schema that was never registered answers "permission denied"
		if strings.Contains(err.Error(), "permission denied") {
			return missingSchemaError(org, err)
		}
		return fmt.Errorf("validating spec in Quasar: %v", err)
	}
	return quasarVerdict(org, resp)
}

func missingSchemaError(org string, cause error) error {
	return fmt.Errorf("schema %s/%s/%s@%s is not registered in Quasar (or you may not read it: %v); register it with:\n"+
		"  cockpit create schema --org %s --namespace %s --schema-name %s --version %s \\\n"+
		"    --path tools/bootc/specs/host-image-spec.schema.yaml",
		org, applySchemaNamespace, hostImageSchemaName, applySchemaVersion, strings.TrimSpace(cause.Error()),
		org, applySchemaNamespace, hostImageSchemaName, applySchemaVersion)
}

func quasarVerdict(org string, resp validateConfigurationResponse) error {
	if resp.IsValid {
		return nil
	}
	if strings.Contains(resp.Message, "No schema with key") {
		return missingSchemaError(org, fmt.Errorf("%s", resp.Message))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "spec does not match schema %s@%s:", hostImageSchemaName, applySchemaVersion)
	for _, m := range strings.Split(resp.Message, "; ") {
		fmt.Fprintf(&b, "\n  - %s", m)
	}
	return fmt.Errorf("%s", b.String())
}
