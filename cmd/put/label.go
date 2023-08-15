package put

import (
	"context"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/magnetar/pkg/api"
	"github.com/spf13/cobra"
	"log"
)

const (
	idFlag         = "nodeid"
	idFlagShort    = "i"
	keyFlag        = "key"
	keyFlagSHort   = "k"
	valueFlag      = "value"
	valueFlagShort = "v"
)

func init() {
	LabelCmd.PersistentFlags().StringP(idFlag, idFlagShort, "", "node id")
	err := LabelCmd.MarkPersistentFlagRequired(idFlag)
	if err != nil {
		log.Fatalln(err)
	}
	LabelCmd.PersistentFlags().StringP(keyFlag, keyFlagSHort, "", "label key")
	err = LabelCmd.MarkPersistentFlagRequired(keyFlag)
	if err != nil {
		log.Fatalln(err)
	}

	BoolLabelCmd.Flags().BoolP(valueFlag, valueFlagShort, false, "label value")
	err = BoolLabelCmd.MarkFlagRequired(valueFlag)
	if err != nil {
		log.Fatalln(err)
	}
	Float64LabelCmd.Flags().Float64P(valueFlag, valueFlagShort, 0, "label value")
	err = Float64LabelCmd.MarkFlagRequired(valueFlag)
	if err != nil {
		log.Fatalln(err)
	}
	StringLabelCmd.Flags().StringP(valueFlag, valueFlagShort, "", "label value")
	err = StringLabelCmd.MarkFlagRequired(valueFlag)
	if err != nil {
		log.Fatalln(err)
	}
}

var LabelCmd = &cobra.Command{
	Use:   "nodelabel",
	Short: "Put node label",
}

var BoolLabelCmd = &cobra.Command{
	Use:   "bool",
	Short: "Put node bool label",
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
		value, err := cmd.Flags().GetBool(valueFlag)
		if err != nil {
			render.Error(err)
			return
		}

		resp, err := clients.Magnetar.PutBoolLabel(context.Background(), &api.PutBoolLabelReq{
			NodeId: id,
			Label: &api.BoolLabel{
				Key:   key,
				Value: value,
			},
		})
		if err != nil {
			render.Error(err)
			return
		}
		render.Node(resp.Node)
	},
}

var Float64LabelCmd = &cobra.Command{
	Use:   "float64",
	Short: "Put node float64 label",
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
		value, err := cmd.Flags().GetFloat64(valueFlag)
		if err != nil {
			render.Error(err)
			return
		}

		resp, err := clients.Magnetar.PutFloat64Label(context.Background(), &api.PutFloat64LabelReq{
			NodeId: id,
			Label: &api.Float64Label{
				Key:   key,
				Value: value,
			},
		})
		if err != nil {
			render.Error(err)
			return
		}
		render.Node(resp.Node)
	},
}

var StringLabelCmd = &cobra.Command{
	Use:   "string",
	Short: "Put node string label",
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
		value, err := cmd.Flags().GetString(valueFlag)
		if err != nil {
			render.Error(err)
			return
		}

		resp, err := clients.Magnetar.PutStringLabel(context.Background(), &api.PutStringLabelReq{
			NodeId: id,
			Label: &api.StringLabel{
				Key:   key,
				Value: value,
			},
		})
		if err != nil {
			render.Error(err)
			return
		}
		render.Node(resp.Node)
	},
}
