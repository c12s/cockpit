package cmd

import (
	"errors"
	"fmt"
	"github.com/c12s/cockpit/cmd/model"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"strings"
)

var ContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Init empty CLI context environment, to interact with region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}

func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Context already exists.")
	}

	return nil
}

func initContext(path, address string) error {
	s := []string{path, "ccontext.yml"}
	filename := strings.Join(s, "/")

	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	c := &model.CContext{
		Context: &model.Content{
			Version:   "v1",
			Address:   address,
			Namespace: "default",
			User:      "",
		},
	}

	nerr, data := model.Marshall(c)
	if nerr != nil {
		return nerr
	}

	fmt.Fprintf(file, data)

	return nil
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init empty CLI context environment, to interact with region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		a := cmd.Flag("address").Value.String()
		usr, err := user.Current()
		if err != nil {
			fmt.Println(err)
		}
		s := []string{usr.HomeDir, "ccontext"}
		contextPath := strings.Join(s, "/")
		fmt.Printf("Empty context initialized in %s. run 'cockpit context login'\n", contextPath)
		err = createDir(contextPath)
		if err != nil {
			fmt.Println(err)
		}

		initContext(contextPath, a)
	},
}

func doLogin(username, password string) {}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login user on inited context environment, to interact with region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		u := cmd.Flag("username").Value.String()
		p := cmd.Flag("password").Value.String()

		doLogin(u, p)
	},
}

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout user on inited context environment, to interact with region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}
