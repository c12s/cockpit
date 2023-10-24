package put

import (
	"context"
	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/render"
	"github.com/c12s/oort/pkg/api"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func init() {
	PolicyCmd.Flags().StringP(fileFlag, fileFlagShort, "", "request file")
	err := PolicyCmd.MarkFlagRequired(fileFlag)
	if err != nil {
		log.Fatalln(err)
	}
}

var PolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Put policy",
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

		req := PolicyReq{}
		err = yaml.Unmarshal(yamlReq, &req)
		if err != nil {
			render.Error(err)
			return
		}

		_, err = clients.OortAdministrator.CreatePolicy(context.Background(), &api.CreatePolicyReq{
			SubjectScope: &api.Resource{
				Id:   req.Policy.SubjectScope.Id,
				Kind: req.Policy.SubjectScope.Kind,
			},
			ObjectScope: &api.Resource{
				Id:   req.Policy.ObjectScope.Id,
				Kind: req.Policy.ObjectScope.Kind,
			},
			Permission: &api.Permission{
				Name: req.Policy.Permission.Name,
				Kind: model.PermKindFromDomain(req.Policy.Permission.Kind),
				Condition: &api.Condition{
					Expression: req.Policy.Permission.Condition,
				},
			},
		})
		if err != nil {
			render.Error(err)
			return
		}
	},
}

type PolicyReq struct {
	Policy model.Policy `yaml:"Policy"`
}
