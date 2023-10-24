package put

import (
	"context"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/kuiper/pkg/api"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const (
	fileFlag      = "file"
	fileFlagShort = "f"
)

func init() {
	ConfigCmd.Flags().StringP(fileFlag, fileFlagShort, "", "request file")
	err := ConfigCmd.MarkFlagRequired(fileFlag)
	if err != nil {
		log.Fatalln(err)
	}
}

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Put new config version",
	Run: func(cmd *cobra.Command, args []string) {
		filepath, err := cmd.Flags().GetString(fileFlag)
		if err != nil {
			render.Error(err)
			return
		}

		yamlReq, err := os.ReadFile(filepath)
		if err != nil {
			render.Error(err)
			return
		}

		req := ConfigReq{}
		err = yaml.Unmarshal(yamlReq, &req)
		if err != nil {
			render.Error(err)
			return
		}

		resp, err := clients.Kuiper.PutConfigGroup(context.Background(), &api.PutConfigGroupReq{
			Group: &api.ConfigGroup{
				Name:    req.Group.Name,
				OrgId:   req.Group.OrgId,
				Version: req.Group.Version,
				Configs: model.ConfigsFromDomain(req.Group.Configs),
			},
			SubId:   req.SubId,
			SubKind: req.SubKind,
		})
		if err != nil {
			render.Error(err)
			return
		}

		config := model.ConfigGroup{
			Name:    resp.Group.Name,
			OrgId:   resp.Group.OrgId,
			Version: resp.Group.Version,
			Configs: model.ConfigsToDomain(resp.Group.Configs),
		}
		render.Config(config)
	},
}

type ConfigReq struct {
	Group   model.ConfigGroup `yaml:"Group"`
	SubId   string            `yaml:"SubId"`
	SubKind string            `yaml:"SubKind"`
}
