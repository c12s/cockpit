package helper

import (
	"bytes"
	"encoding/json"
	"errors"
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

func checkRedirectFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", via[0].Header.Get("Authorization"))
	return nil
}

func newClient(t time.Duration) *http.Client {
	var netClient = &http.Client{
		Timeout:       t,
		CheckRedirect: checkRedirectFunc,
	}
	return netClient
}

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

func Get(timeout time.Duration, url string, h map[string]string) (error, interface{}) {
	if strings.Contains(url, "namespaces") {
		return GetNSJson(timeout, url, h)
	} else if strings.Contains(url, "configs") {
		return GetConfigsJson(timeout, url, h)
	} else if strings.Contains(url, "secrets") {
		return GetSecretsJson(timeout, url, h)
	} else if strings.Contains(url, "actions") {
		return GetActionsJson(timeout, url, h)
	} else if strings.Contains(url, "trace/get") {
		return GetTraceJson(timeout, url, h)
	} else if strings.Contains(url, "trace/list") {
		return GetListTraceJson(timeout, url, h)
	}
	return errors.New("undefined kind"), nil
}

func GetListTraceJson(timeout time.Duration, url string, h map[string]string) (error, *request.Traces) {
	var netClient = newClient(timeout)
	req, err := http.NewRequest("GET", url, nil)
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	rsp := string(body)
	if resp.StatusCode == http.StatusOK {
		s := &request.Traces{}
		err = json.Unmarshal([]byte(rsp), &s)
		if err != nil {
			return err, nil
		}
		return nil, s
	}

	var s map[string]string
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, nil
	}

	fmt.Println(fmt.Sprintf("Statuss code: %d Message: %s", resp.StatusCode, s["message"]))
	return nil, nil
}

func GetTraceJson(timeout time.Duration, url string, h map[string]string) (error, *request.Trace) {
	var netClient = newClient(timeout)
	req, err := http.NewRequest("GET", url, nil)
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	rsp := string(body)
	if resp.StatusCode == http.StatusOK {
		rsp = strings.Replace(rsp, "\\", "", -1)
		rsp = strings.TrimSuffix(rsp, "\"")
		rsp = strings.TrimPrefix(rsp, "\"")

		s := &request.Trace{}
		err = json.Unmarshal([]byte(rsp), &s)
		if err != nil {
			return err, nil
		}
		return nil, s
	}

	var s map[string]string
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, nil
	}

	fmt.Println(fmt.Sprintf("Statuss code: %d Message: %s", resp.StatusCode, s["message"]))
	return nil, nil
}

func GetNSJson(timeout time.Duration, url string, h map[string]string) (error, *request.NSResponse) {
	var netClient = newClient(timeout)
	req, err := http.NewRequest("GET", url, nil)
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	rsp := string(body)

	if resp.StatusCode == http.StatusOK {
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

	var s map[string]string
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, nil
	}

	fmt.Println(fmt.Sprintf("Statuss code: %d Message: %s", resp.StatusCode, s["message"]))
	return nil, nil
}

func GetConfigsJson(timeout time.Duration, url string, h map[string]string) (error, *request.ConfigResponse) {
	var netClient = newClient(timeout)
	// var netClient = &http.Client{
	// 	Timeout: timeout,
	// }
	req, err := http.NewRequest("GET", url, nil)
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	rsp := string(body)

	if resp.StatusCode == http.StatusOK {
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

	var s map[string]string
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, nil
	}

	fmt.Println(fmt.Sprintf("Statuss code: %d Message: %s", resp.StatusCode, s["message"]))
	return nil, nil
}

func GetActionsJson(timeout time.Duration, url string, h map[string]string) (error, *request.ActionsResponse) {
	var netClient = newClient(timeout)
	// var netClient = &http.Client{
	// 	Timeout: timeout,
	// }
	req, err := http.NewRequest("GET", url, nil)
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	rsp := string(body)

	if resp.StatusCode == http.StatusOK {
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

	var s map[string]string
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, nil
	}

	fmt.Println(fmt.Sprintf("Statuss code: %d Message: %s", resp.StatusCode, s["message"]))
	return nil, nil
}

func GetSecretsJson(timeout time.Duration, url string, h map[string]string) (error, *request.SecretsResponse) {
	var netClient = newClient(timeout)
	// var netClient = &http.Client{
	// 	Timeout: timeout,
	// }
	req, err := http.NewRequest("GET", url, nil)
	for k, v := range h {
		req.Header.Set(k, v)
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	rsp := string(body)

	if resp.StatusCode == http.StatusOK {
		rsp = strings.Replace(rsp, "\\", "", -1)
		rsp = strings.TrimSuffix(rsp, "\"")
		rsp = strings.TrimPrefix(rsp, "\"")

		s := &request.SecretsResponse{}
		err = json.Unmarshal([]byte(rsp), &s)
		if err != nil {
			return err, nil
		}
		return nil, s
	}

	var s map[string]string
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, nil
	}

	fmt.Println(fmt.Sprintf("Statuss code: %d Message: %s", resp.StatusCode, s["message"]))
	return nil, nil
}

func Print(kind string, data interface{}) {
	if kind == "namespaces" {
		if data.(*request.NSResponse) != nil {
			NSPrint(data.(*request.NSResponse))
		}
	} else if kind == "configs" {
		if data.(*request.ConfigResponse) != nil {
			ConfigsPrint(data.(*request.ConfigResponse))
		}
	} else if kind == "secrets" {
		if data.(*request.SecretsResponse) != nil {
			SecretsPrint(data.(*request.SecretsResponse))
		}
	} else if kind == "actions" {
		if data.(*request.ActionsResponse) != nil {
			ActionsPrint(data.(*request.ActionsResponse))
		}
	} else if kind == "trace/get" {
		if data.(*request.Trace) != nil {
			TracePrint(data.(*request.Trace))
		}
	} else if kind == "trace/list" {
		if data.(*request.Traces) != nil {
			TracesPrint(data.(*request.Traces))
		}
	}
}

func TracePrint(data *request.Trace) {
	fmt.Println("DEMO")
	fmt.Println(data)
}

func TracesPrint(data *request.Traces) {
	fmt.Println("DEMO")
	fmt.Println(data)
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

func ConfigsPrint(resp *request.ConfigResponse) {
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

func SecretsPrint(resp *request.SecretsResponse) {
	if len(resp.Result) > 0 {
		// initialize tabwriter
		w := new(tabwriter.Writer)
		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)
		defer w.Flush()

		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "RegionID", "ClusterID", "NodeID", "Secrets")
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", "----", "----", "----", "----")
		for _, rez := range resp.Result {
			fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t", rez.RegionId, rez.ClusterId, rez.NodeId, rez.Secrets)
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

func PostCallExtractToken(timeout time.Duration, address, data string) (error, string, string) {
	var netClient = &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", address, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := netClient.Do(req)
	if err != nil {
		return err, "", ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, "", ""
	}
	rsp := string(body)

	var s map[string]string
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, "", ""
	}

	return nil, fmt.Sprintf("Statuss code: %d Message: %s",
		resp.StatusCode, s["message"]), resp.Header.Get("Auth-Token")
}

func Post(timeout time.Duration, address, data string, headers map[string]string) (error, string) {
	var netClient = newClient(timeout)
	// var netClient = &http.Client{
	// 	Timeout: timeout,
	// }

	req, err := http.NewRequest("POST", address, bytes.NewBuffer([]byte(data)))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return err, ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, ""
	}
	rsp := string(body)

	var s map[string]string
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, ""
	}

	return nil, fmt.Sprintf("Statuss code: %d Message: %s", resp.StatusCode, s["message"])
}
