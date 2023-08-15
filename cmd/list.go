package cmd

import (
	"context"
	"fmt"
	"github.com/c12s/magnetar/pkg/api"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"os"
	"strings"
)

const (
	labelFlag      = "label"
	labelFlagShort = "l"
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringArrayP(labelFlag, labelFlagShort, []string{}, "node query selector <key>[=|!=|>|<]<value>")
}

var listHandlers = map[string]func(cmd *cobra.Command, args []string){
	"nodes": listNodes,
	"node":  listNodes,
}

var listCmd = &cobra.Command{
	Use:       "list",
	Short:     "",
	ValidArgs: maps.Keys(listHandlers),
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.OnlyValidArgs(cmd, args)
		if err != nil {
			return err
		}
		return cobra.ExactArgs(1)(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list")
		listHandlers[args[0]](cmd, args[1:])
	},
}

func listNodes(cmd *cobra.Command, args []string) {
	selector, err := cmd.Flags().GetStringArray(labelFlag)
	if err != nil || len(selector) == 0 {
		fmt.Println("listing nodes")
		resp, err := magnetar.ListNodes(context.Background(), &api.ListNodesReq{})
		if err != nil {
			fmt.Println(err)
			return
		}
		renderNodes(resp.Nodes)
		return
	}
	fmt.Println("querying nodes")
	req := &api.QueryNodesReq{
		Queries: make([]*api.Query, 0),
	}
	for _, q := range selector {
		op := ">"
		split := strings.Split(q, op)
		if len(split) == 2 {
			query := &api.Query{
				LabelKey: split[0],
				ShouldBe: op,
				Value:    split[1],
			}
			req.Queries = append(req.Queries, query)
			continue
		}

		op = "<"
		split = strings.Split(q, op)
		if len(split) == 2 {
			query := &api.Query{
				LabelKey: split[0],
				ShouldBe: op,
				Value:    split[1],
			}
			req.Queries = append(req.Queries, query)
			continue
		}

		op = "!="
		split = strings.Split(q, op)
		if len(split) == 2 {
			query := &api.Query{
				LabelKey: split[0],
				ShouldBe: op,
				Value:    split[1],
			}
			req.Queries = append(req.Queries, query)
			continue
		}

		op = "="
		split = strings.Split(q, op)
		if len(split) == 2 {
			query := &api.Query{
				LabelKey: split[0],
				ShouldBe: op,
				Value:    split[1],
			}
			req.Queries = append(req.Queries, query)
			continue
		}
	}
	resp, err := magnetar.QueryNodes(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	renderNodes(resp.Nodes)
}

func renderNodes(nodes []*api.NodeStringified) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Node ID", "Labels"})
	for _, node := range nodes {
		t.AppendRow(table.Row{node.Id, formatLabels(node.Labels)})
		t.AppendSeparator()
	}
	t.Render()
}

func formatLabels(labels []*api.LabelStringified) (formatted string) {
	for i, label := range labels {
		if i > 0 {
			formatted += "\n"
		}
		formatted += fmt.Sprintf("%s=%s", label.Key, label.Value)
	}
	return
}
