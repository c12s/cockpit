package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"log"
	"os"
	"time"

	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

const (
	createRelationsShortDescription = "Create relations between entities"
	createRelationsLongDescription  = "This command creates relations between entities specified by their IDs and kinds.\nFor example:\n\ncreate relations --ids 'id1|id2' --kinds 'kind1|kind2'\n\n" +
		"Example:\n" +
		"create-relations --ids 'c12s|dev' --kinds 'org|namespace'"

	// Flag Constants
	idsFlag   = "ids"
	kindsFlag = "kinds"

	// Flag Shorthand Constants
	idsFlagShorthand   = "i"
	kindsFlagShorthand = "k"

	// Flag Descriptions
	idsDesc   = "IDs of the entities separated by '|' (required)"
	kindsDesc = "Kinds of the entities separated by '|' (required)"
)

var (
	ids   string
	kinds string
)

var CreateRelationsCmd = &cobra.Command{
	Use:   "relations",
	Short: createRelationsShortDescription,
	Long:  createRelationsLongDescription,
	Run:   executeCreateRelations,
}

func executeCreateRelations(cmd *cobra.Command, args []string) {
	idsList, kindsList, err := utils.ParseIDsAndKinds(ids, kinds)
	if err != nil {
		log.Fatal(err)
	}

	relation := model.Relation{
		From: model.Entity{ID: idsList[0], Kind: kindsList[0]},
		To:   model.Entity{ID: idsList[1], Kind: kindsList[1]},
	}

	if err := sendCreateRelationsRequest(relation); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Relations created successfully")
}

func sendCreateRelationsRequest(relation model.Relation) error {
	config, err := prepareRelationsRequestConfig(relation)
	if err != nil {
		fmt.Printf("Error creating request config: %v\n", err)
		os.Exit(1)
	}

	err = utils.SendHTTPRequest(config)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	return nil
}

func prepareRelationsRequestConfig(relation model.Relation) (model.HTTPRequestConfig, error) {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return model.HTTPRequestConfig{}, fmt.Errorf("error reading token: %v", err)
	}

	url := clients.BuildURL("core", "v1", "CreateInheritanceRel")

	return model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		RequestBody: relation,
		Token:       token,
		Timeout:     10 * time.Second,
	}, nil
}

func init() {
	CreateRelationsCmd.Flags().StringVarP(&ids, idsFlag, idsFlagShorthand, "", idsDesc)
	CreateRelationsCmd.Flags().StringVarP(&kinds, kindsFlag, kindsFlagShorthand, "", kindsDesc)

	CreateRelationsCmd.MarkFlagRequired(idsFlag)
	CreateRelationsCmd.MarkFlagRequired(kindsFlag)
}
