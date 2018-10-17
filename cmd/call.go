package cmd

import (
	"bytes"
	"fmt"
	"github.com/c12s/cockpit/cmd/model"
	"io/ioutil"
	"net/http"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func formCall(artifact, path string, c *model.CContext, query map[string]string) string {
	s := []string{"http:/", c.Context.Address, "api", c.Context.Version, artifact, path}
	pathSegment := strings.Join(s, "/")

	if len(query) == 0 {
		return pathSegment
	}

	q := []string{}
	for k, v := range query {
		pair := strings.Join([]string{k, v}, "=")
		q = append(q, pair)
	}
	querySegment := strings.Join(q, "&")

	return strings.Join([]string{pathSegment, querySegment}, "?")
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
