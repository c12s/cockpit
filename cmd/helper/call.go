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

func GetJson(timeout time.Duration, url string) (error, *request.Response) {
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

	fmt.Println(rsp)

	s := &request.Response{}
	err = json.Unmarshal([]byte(rsp), &s)
	if err != nil {
		return err, nil
	}

	return nil, s
}

func Pprint(resp string) {
	rez := map[string]string{}
	val := strings.Split(resp, ",")
	if len(val) > 1 {
		for _, v := range val {
			r := strings.Replace(strings.Replace(v, "{", "", -1), "}", "", -1)
			kv := strings.Split(r, ":")
			k := strings.Replace(strings.Replace(kv[0], "\"", "", -1), "\\", "", -1)
			v := strings.Replace(strings.Replace(kv[1], "\"", "", -1), "\\", "", -1)
			rez[k] = v
		}
		// initialize tabwriter
		w := new(tabwriter.Writer)
		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)
		defer w.Flush()
		fmt.Fprintf(w, "\n %s\t%s\t%s\t", "Namespace", "Name", "Age")
		fmt.Fprintf(w, "\n %s\t%s\t%s\t", "----", "----", "----")
		fmt.Fprintf(w, "\n %s\t%s\t%s\t", rez["namespace"], rez["name"], rez["age"])
		fmt.Fprintf(w, "\n")
	} else {
		fmt.Println(resp)
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
