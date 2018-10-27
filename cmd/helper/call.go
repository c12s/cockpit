package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/c12s/cockpit/cmd/model"
	"github.com/c12s/cockpit/cmd/model/request"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/tabwriter"
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

func GetNSJson(timeout time.Duration, url string) (error, *request.NSResponse) {
	var myClient = &http.Client{
		Timeout: timeout,
	}
	r, err := myClient.Get(url)
	if err != nil {
		return err, nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err, nil
	}
	r.Body.Close()

	rsp := string(body)

	rsp = strings.Replace(rsp, "\\", "", -1)
	rsp = strings.TrimSuffix(rsp, "\"")
	rsp = strings.TrimPrefix(rsp, "\"")

	s := &request.NSResponse{}
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, nil
	}

	return nil, s
}

func GetConfigsJson(timeout time.Duration, url string) (error, *request.ConfigResponse) {
	var myClient = &http.Client{
		Timeout: timeout,
	}
	r, err := myClient.Get(url)
	if err != nil {
		return err, nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err, nil
	}
	r.Body.Close()

	rsp := string(body)

	rsp = strings.Replace(rsp, "\\", "", -1)
	rsp = strings.TrimSuffix(rsp, "\"")
	rsp = strings.TrimPrefix(rsp, "\"")

	s := &request.ConfigResponse{}
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, nil
	}

	return nil, s
}

func GetActionsJson(timeout time.Duration, url string) (error, *request.ActionsResponse) {
	var myClient = &http.Client{
		Timeout: timeout,
	}
	r, err := myClient.Get(url)
	if err != nil {
		return err, nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err, nil
	}
	r.Body.Close()

	rsp := string(body)

	rsp = strings.Replace(rsp, "\\", "", -1)
	rsp = strings.TrimSuffix(rsp, "\"")
	rsp = strings.TrimPrefix(rsp, "\"")

	s := &request.ActionsResponse{}
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, nil
	}

	return nil, s
}

func ActionsPrint(resp *request.ActionsResponse) {
	if len(resp.Result) > 0 {
		for _, rez := range resp.Result {
			fmt.Printf("RegionID: %s\n", rez.RegionId)
			fmt.Printf("ClusterID: %s\n", rez.ClusterId)
			fmt.Printf("NodeID: %s\n", rez.NodeId)
			for k, v := range rez.Actions {
				rsp := strings.Replace(v, ",", "\n\t", -1)
				t := strings.Split(k, "_")
				fmt.Printf("Timestamp: %s\n", t[1])
				fmt.Printf("\t%s\n", rsp)
				fmt.Println("")
			}
			fmt.Println("")
		}
	} else {
		fmt.Println("No results")
	}
}

func NSPrint(resp *request.NSResponse) {
	if len(resp.Result) > 0 {
		// initialize tabwriter
		w := new(tabwriter.Writer)
		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)
		defer w.Flush()

		fmt.Fprintf(w, "\n %s\t%s\t%s\t", "Namespace", "Name", "Age")
		fmt.Fprintf(w, "\n %s\t%s\t%s\t", "----", "----", "----")
		for _, rez := range resp.Result {
			fmt.Fprintf(w, "\n %s\t%s\t%s\t", rez.Namespace, rez.Name, rez.Age)
		}
		fmt.Fprintf(w, "\n")
	} else {
		fmt.Println("No results")
	}
}

func ConfigPrint(resp *request.ConfigResponse) {
	if len(resp.Result) > 0 {
		// initialize tabwriter
		w := new(tabwriter.Writer)
		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)
		defer w.Flush()

		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "RegionID", "ClusterID", "NodeID", "Configs")
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "----", "----", "----", "----")
		for _, rez := range resp.Result {
			fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", rez.RegionId, rez.ClusterId, rez.NodeId, rez.Configs)
		}
		fmt.Fprintf(w, "\n")
	} else {
		fmt.Println("No results")
	}
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
