package apply

import (
	"context"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/kuiper/pkg/api"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
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
	Short: "Apply config version to a namespace",
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

		_, err = clients.Kuiper.ApplyConfigGroup(context.Background(), &api.ApplyConfigGroupReq{
			GroupName: req.Group.Name,
			OrgId:     req.Group.OrgId,
			Version:   req.Group.Version,
			Namespace: req.Namespace,
			Queries:   queriesFromDomain(req.Queries),
			SubId:     req.SubId,
			SubKind:   req.SubKind,
		})
		if err != nil {
			render.Error(err)
			return
		}
	},
}

type ConfigReq struct {
	Group     model.ConfigGroup `yaml:"Group"`
	Namespace string            `yaml:"Namespace"`
	Queries   []model.Query     `yaml:"Queries"`
	SubId     string            `yaml:"SubId"`
	SubKind   string            `yaml:"SubKind"`
}

func queriesFromDomain(queries []model.Query) []*magnetarapi.Query {
	resp := make([]*magnetarapi.Query, len(queries))
	for i, query := range queries {
		resp[i] = &magnetarapi.Query{
			LabelKey: query.Key,
			ShouldBe: query.ShouldBe,
			Value:    query.Value,
		}
	}
	return resp
}
