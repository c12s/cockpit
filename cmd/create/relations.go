package cmd

import (
	"fmt"
	"github.com/c12s/cockpit/clients"
	"log"
	"os"
	"strings"
	"time"

	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

const (
	idsFlag   = "ids"
	kindsFlag = "kinds"
	idsDesc   = "IDs of the entities separated by '|' (required)"
	kindsDesc = "Kinds of the entities separated by '|' (required)"
)

var (
	ids   string
	kinds string
)

var CreateRelationsCmd = &cobra.Command{
	Use:   "relations",
	Short: "Create relations between entities",
	Long: `This command creates relations between entities specified by their IDs and kinds.
For example:

create relations --ids 'id1|id2' --kinds 'kind1|kind2'`,
	Run: executeCreateRelations,
}

func executeCreateRelations(cmd *cobra.Command, args []string) {
	idsList, kindsList, err := parseIDsAndKinds(ids, kinds)
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
}

func parseIDsAndKinds(ids, kinds string) ([]string, []string, error) {
	idsList := strings.Split(ids, "|")
	kindsList := strings.Split(kinds, "|")

	if len(idsList) != 2 || len(kindsList) != 2 {
		return nil, nil, fmt.Errorf("invalid ids or kinds format. Please provide exactly two values separated by '|'")
	}
	return idsList, kindsList, nil
}

func sendCreateRelationsRequest(relation model.Relation) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading token: %v\n", err)
		os.Exit(1)
	}

	url := clients.BuildURL("core", "v1", "CreateInheritanceRel")

	config := model.HTTPRequestConfig{
		Method:      "POST",
		URL:         url,
		Token:       token,
		Timeout:     10 * time.Second,
		RequestBody: relation,
		Response:    nil,
	}

	err = utils.SendHTTPRequest(config)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	fmt.Println("Relations created successfully")
	return nil
}

func init() {
	CreateRelationsCmd.Flags().StringVarP(&ids, idsFlag, "i", "", idsDesc)
	CreateRelationsCmd.Flags().StringVarP(&kinds, kindsFlag, "k", "", kindsDesc)

	CreateRelationsCmd.MarkFlagRequired(idsFlag)
	CreateRelationsCmd.MarkFlagRequired(kindsFlag)
}
