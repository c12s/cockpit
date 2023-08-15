package get

import (
	"context"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/magnetar/pkg/api"
	"github.com/spf13/cobra"
	"log"
)

const (
	idFlag      = "id"
	idFlagShort = "i"
)

func init() {
	NodeCmd.Flags().StringP(idFlag, idFlagShort, "", "node id")
	err := NodeCmd.MarkFlagRequired(idFlag)
	if err != nil {
		log.Fatalln(err)
	}
}

var NodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Get node by id",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString(idFlag)
		if err != nil {
			render.Error(err)
			return
		}
		resp, err := clients.Magnetar.GetNode(context.Background(), &api.GetNodeReq{
			NodeId: id,
		})
		if err != nil {
			render.Error(err)
			return
		}
		render.Node(resp.Node)
	},
}
