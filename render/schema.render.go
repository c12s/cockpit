package render

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/c12s/cockpit/model"
)

func RenderSchemaVersionsTabWriter(versions []model.SchemaVersion) {
	if len(versions) == 0 {
		fmt.Println("No schema versions were found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "Organization\tNamespace\tSchema Name\tVersion\tCreation Time\t")

	for _, version := range versions {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\t\n",
			version.SchemaDetails.Organization,
			version.SchemaDetails.Namespace,
			version.SchemaDetails.SchemaName,
			version.SchemaDetails.Version,
			version.SchemaData.CreationTime)
	}
}

func RenderSchemaTabWriter(schema model.SchemaData) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "Schema\tCreation Time\t")

	fmt.Fprintf(w, "%s\t%s\t\t\n",
		schema.Schema,
		schema.CreationTime)
}
