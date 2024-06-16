package render

import (
	"fmt"
	"github.com/c12s/cockpit/model"
	"os"
	"text/tabwriter"
)

func RenderTasksTabWriter(tasks []model.Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks were found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "ID\tNode\tNamespace\tStatus\tAccepted At\tResolved At\t")

	for _, task := range tasks {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n",
			task.ID,
			task.Node,
			task.Namespace,
			task.Status,
			task.AcceptedAt,
			task.ResolvedAt)
	}
}

func RenderConfigGroupsTabWriter(groups []model.ConfigGroup) {
	if len(groups) == 0 {
		fmt.Println("No configuration groups were found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "Organization\tName\tVersion\tCreated At\tParam Set Name\tParams\t")

	for _, group := range groups {
		for _, paramSet := range group.ParamSets {
			for _, param := range paramSet.ParamSet {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s=%s\t\n", group.Organization, group.Name, group.Version, group.CreatedAt, paramSet.Name, param.Key, param.Value)
			}
		}
	}
}

func RenderConfigGroupTabWriter(group model.ConfigGroup) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "Organization\tName\tVersion\tCreated At\tParam Set Name\tParams\t")

	for _, paramSet := range group.ParamSets {
		for _, param := range paramSet.ParamSet {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s=%s\t\n", group.Organization, group.Name, group.Version, group.CreatedAt, paramSet.Name, param.Key, param.Value)
		}
	}
}

func RenderConfigGroupDiffsTabWriter(diffResponse model.ConfigGroupDiffResponse) {
	if len(diffResponse.Diffs) == 0 {
		fmt.Println("No diffs were found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "Category\tKey\tValue\tChange\t")

	for category, diffSet := range diffResponse.Diffs {
		for _, diff := range diffSet.Diffs {
			switch diff.Type {
			case "deletion":
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n",
					category,
					diff.Diff.Key,
					diff.Diff.Value,
					"-",
				)
			case "replacement":
				fmt.Fprintf(w, "%s\t%s\t%s\t%s -> %s\t\n",
					category,
					diff.Diff.Key,
					diff.Diff.OldValue,
					diff.Diff.NewValue,
					"+",
				)
			}
		}
	}
}
