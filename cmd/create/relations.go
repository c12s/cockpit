package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/aliases"
	"github.com/c12s/cockpit/clients"
	"os"
	"time"

	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

const (
	createRelationsShortDescription = "Create relations between entities"
	createRelationsLongDescription  = `This command creates relations between entities specified by their IDs and kinds.
Relations help to establish a hierarchical or dependency structure between different entities within the organization. This can include relationships between organizations, namespaces, and other resources.

Example:
- cockpit create relations --ids 'myOrg|dev' --kinds 'org|namespace'`

	// Flag Constants
	idsFlag   = "ids"
	kindsFlag = "kinds"

	// Flag Shorthand Constants
	idsShorthandFlag   = "i"
	kindsShorthandFlag = "k"

	// Flag Descriptions
	idsDescription   = "IDs of the entities separated by '|' (required)"
	kindsDescription = "Kinds of the entities separated by '|' (required)"
)

var (
	ids   string
	kinds string
)

var CreateRelationsCmd = &cobra.Command{
	Use:     "relations",
	Aliases: aliases.RelationsAliases,
	Short:   createRelationsShortDescription,
	Long:    createRelationsLongDescription,
	Run:     executeCreateRelations,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return utils.ValidateRequiredFlags(cmd, []string{idsFlag, kindsFlag})
	},
}

func executeCreateRelations(cmd *cobra.Command, args []string) {
	idsList, kindsList, err := utils.ParseIDsAndKinds(ids, kinds)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		os.Exit(1)
	}

	relation := model.Relation{
		From: model.Entity{ID: idsList[0], Kind: kindsList[0]},
		To:   model.Entity{ID: idsList[1], Kind: kindsList[1]},
	}

	if err := sendCreateRelationsRequest(relation); err != nil {
		fmt.Println("Error sending relations  request:", err)
		os.Exit(1)
	}
	fmt.Println("Relations created successfully")
	fmt.Println()
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
	CreateRelationsCmd.Flags().StringVarP(&ids, idsFlag, idsShorthandFlag, "", idsDescription)
	CreateRelationsCmd.Flags().StringVarP(&kinds, kindsFlag, kindsShorthandFlag, "", kindsDescription)

	CreateRelationsCmd.MarkFlagRequired(idsFlag)
	CreateRelationsCmd.MarkFlagRequired(kindsFlag)
}
