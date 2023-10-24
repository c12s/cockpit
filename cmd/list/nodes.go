package list

import (
	"context"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/magnetar/pkg/api"
	"github.com/spf13/cobra"
	"strings"
)

const (
	labelFlag      = "label"
	labelFlagShort = "l"
)

func init() {
	NodesCmd.Flags().StringArrayP(labelFlag, labelFlagShort, []string{}, "node query selector <key>[=|!=|>|<]<value>")
}

var NodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "List or query nodes",
	Run: func(cmd *cobra.Command, args []string) {
		selector, err := cmd.Flags().GetStringArray(labelFlag)
		var nodes []*api.NodeStringified
		if err != nil || len(selector) == 0 {
			resp, err := listNodesReq()
			if err != nil {
				render.Error(err)
				return
			}
			nodes = resp
		} else {
			resp, err := queryNodesReq(selector)
			if err != nil {
				render.Error(err)
				return
			}
			nodes = resp
		}
		render.Nodes(nodes)
	},
}

func listNodesReq() ([]*api.NodeStringified, error) {
	resp, err := clients.Magnetar.ListNodes(context.Background(), &api.ListNodesReq{})
	if err != nil {
		return nil, err
	}
	return resp.Nodes, nil
}

func queryNodesReq(selector []string) ([]*api.NodeStringified, error) {
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
	resp, err := clients.Magnetar.QueryNodes(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return resp.Nodes, nil
}
