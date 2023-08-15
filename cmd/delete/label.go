package delete

import (
	"context"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/magnetar/pkg/api"
	"github.com/spf13/cobra"
	"log"
)

const (
	idFlag       = "nodeid"
	idFlagShort  = "i"
	keyFlag      = "key"
	keyFlagSHort = "k"
)

func init() {
	LabelCmd.Flags().StringP(idFlag, idFlagShort, "", "node id")
	err := LabelCmd.MarkFlagRequired(idFlag)
	if err != nil {
		log.Fatalln(err)
	}
	LabelCmd.Flags().StringP(keyFlag, keyFlagSHort, "", "label key")
	err = LabelCmd.MarkFlagRequired(keyFlag)
	if err != nil {
		log.Fatalln(err)
	}
}

var LabelCmd = &cobra.Command{
	Use:   "nodelabel",
	Short: "Delete node label",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString(idFlag)
		if err != nil {
			render.Error(err)
			return
		}
		key, err := cmd.Flags().GetString(keyFlag)
		if err != nil {
			render.Error(err)
			return
		}
		resp, err := clients.Magnetar.DeleteLabel(context.Background(), &api.DeleteLabelReq{
			NodeId:   id,
			LabelKey: key,
		})
		if err != nil {
			render.Error(err)
			return
		}
		render.Node(resp.Node)
	},
}
