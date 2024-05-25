package render

import (
	"fmt"
	"github.com/c12s/cockpit/model"
)

func HandleSchemaVersionResponse(response *model.SchemaVersionResponse) {
	println()
	fmt.Println("Message:", response.Message)
	for _, version := range response.SchemaVersions {
		fmt.Println("Schema Name:", version.SchemaDetails.SchemaName)
		fmt.Println("Version:", version.SchemaDetails.Version)
		fmt.Println("Organization:", version.SchemaDetails.Organization)
		fmt.Println("Schema Data:")
		fmt.Println(version.SchemaData.Schema)
		fmt.Println("Creation Time:", version.SchemaData.CreationTime)
		fmt.Println()
	}
}

func HandleSchemaResponse(response *model.SchemaResponse) {
	println()
	fmt.Println("Message:", response.Message)
	if response.SchemaData.Schema != "" {
		fmt.Println("Schema Data:")
		fmt.Println(response.SchemaData.Schema)
		fmt.Println("Creation Time:", response.SchemaData.CreationTime)
	}
	println()
}
