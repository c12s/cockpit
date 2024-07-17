package render

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/c12s/cockpit/model"
)

func RenderStandaloneConfigsTabWriter(configs []model.StandaloneConfig) {
	if len(configs) == 0 {
		fmt.Println("No standalone configurations were found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "Organization\tNamespace\tName\tVersion\tCreated At\tParams\t")

	for _, config := range configs {
		for _, param := range config.ParamSet {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s=%s\t\n", config.Organization, config.Namespace, config.Name, config.Version, config.CreatedAt, param.Key, param.Value)
		}
	}
}

func RenderStandaloneConfigTabWriter(config model.StandaloneConfig) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "Organization\tNamespace\tName\tVersion\tCreated At\tParams\t")

	for _, param := range config.ParamSet {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s=%s\t\n", config.Organization, config.Namespace, config.Name, config.Version, config.CreatedAt, param.Key, param.Value)
	}
}

func RenderStandaloneConfigDiffsTabWriter(diffResponse model.StandaloneConfigDiffResponse) {
	if len(diffResponse.Diffs) == 0 {
		fmt.Println("No diffs were found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "Key\tValue\tChange\t")

	for _, diff := range diffResponse.Diffs {
		switch diff.Type {
		case "deletion":
			fmt.Fprintf(w, "%s\t%s\t%s\t\n",
				diff.Diff["key"],
				diff.Diff["value"],
				"-",
			)
		case "addition":
			fmt.Fprintf(w, "%s\t%s\t%s\t\n",
				diff.Diff["key"],
				diff.Diff["value"],
				"+",
			)
		case "replacement":
			fmt.Fprintf(w, "%s\t%s -> %s\t%s\t\n",
				diff.Diff["key"],
				diff.Diff["old_value"],
				diff.Diff["new_value"],
				"->",
			)
		}
	}
}
