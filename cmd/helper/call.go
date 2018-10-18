package helper

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

func FormCall(artifact, path string, c *model.CContext, query map[string]string) string {
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

func GetContext() (error, *model.CContext) {
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

func GetCall(timeout time.Duration, address string) (error, string) {
	var netClient = &http.Client{
		Timeout: timeout,
	}
	response, err := netClient.Get(address)
	if err != nil {
		return err, ""
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(response.Body)
		if err2 != nil {
			return err, ""
		}
		bodyString := string(bodyBytes)
		return nil, bodyString
	} else {
		return nil, fmt.Sprintf("Resp: %d", response.StatusCode)
	}
}

func PostCall(timeout time.Duration, address, data string) (error, string) {
	var netClient = &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", address, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(req)
	if err != nil {
		return err, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return err2, ""
		}
		return nil, string(body)
	} else {
		return nil, fmt.Sprintf("Resp: %d", resp.StatusCode)
	}
}
