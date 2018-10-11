package cmd

import (
	"bytes"
	"fmt"
	"github.com/c12s/cockpit/cmd/model"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func getCall(timeout time.Duration, address string) {
	var netClient = &http.Client{
		Timeout: timeout,
	}
	response, err := netClient.Get(address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(response.Body)
		if err2 != nil {
			fmt.Println(err)
			return
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	} else {
		fmt.Printf("Resp: %d", response.StatusCode)
	}
}

func postCall(timeout time.Duration, address, data string) {
	var netClient = &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", address, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		fmt.Println("response Body:", string(body))
	}
}

var ConfigsCmd = &cobra.Command{
	Use:   "configs",
	Short: "Get the configurations from region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide some of avalible commands or type help for help")
	},
}

func formCall(path string, c *model.CContext) string {
	s := []string{"http:/", c.Context.Address, "api", c.Context.Version, "configs", path}
	return strings.Join(s, "/")
}

func getContext() (error, *model.CContext) {
	usr, err := user.Current()
	if err != nil {
		return err, nil
	}

	contextPath := filepath.Join(usr.HomeDir, ".constellations", "context.yml")
	err, ctx := model.Context(contextPath)
	if err != nil {
		return err, nil
	}

	return nil, ctx
}

var ConfigsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the configurations from region/s cluster/s node/s job/s",
	Long:  "change all data inside regions, clusters and nodes",
	Run: func(cmd *cobra.Command, args []string) {
		// regions := cmd.Flag("regions").Value.String()
		// clusters := cmd.Flag("clusters").Value.String()

		err, ctx := getContext()
		if err != nil {
			fmt.Println(err)
			return
		}

		callPath := formCall("", ctx)
		getCall(10*time.Second, callPath)
	},
}

var ConfigsMutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Mutate state of the configurations for the region, cluster, node and/or jobs",
	Long:  "change all data inside regions, clusters and nodes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()
		if _, err := os.Stat(file); err == nil {
			f, err := mutateFile(file)
			if err != nil {
				fmt.Println(err)
			}

			err2, data := kind(f)
			if err2 != nil {
				fmt.Println(err2)
				return
			}

			err3, ctx := getContext()
			if err != nil {
				fmt.Println(err3)
				return
			}
			callPath := formCall("new", ctx)
			postCall(10*time.Second, callPath, data)
		} else {
			fmt.Println("File not exists")
		}
	},
}
