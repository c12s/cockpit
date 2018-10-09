package helper

import (
	"bufio"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"github.com/c12s/cockpit/cmd/model/request"
	"os"
	"strings"
)

func toBASE64(data string) string {
	b64Data := b64.StdEncoding.EncodeToString([]byte(data))
	return b64Data
}

func arrToBase(data []string) []string {
	conv := []string{}
	for _, item := range data {
		conv = append(conv, toBASE64(item))
	}
	return conv
}

func readFile(file, kind string) (error, []string) {
	// Open the file.
	f, err := os.Open(file)
	if err != nil {
		return err, nil
	}

	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)

	data := []string{}
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && len(line) > 0 {
			if !strings.Contains(line, "=") {
				return errors.New(fmt.Sprintf("Error: Line %s is not in right format KEY=VALUE.", line)), nil
			} else {
				if kind == SECRETS {
					bLine := toBASE64(line)
					data = append(data, bLine)
				} else if kind == CONFIGS {
					data = append(data, line)
				}
			}
		}
	}

	return nil, data
}

func constructPayload(payload map[string][]string, kind string) []request.Payload {
	retVal := []request.Payload{}
	if len(payload) != 0 {
		for k, v := range payload {
			if k == FILES {
				for _, file := range v {
					err, data := readFile(file, kind)
					if err != nil {
						fmt.Println(err)
						return nil
					}
					p := request.Payload{
						Kind:    constructFileKey(file),
						Content: data,
					}
					retVal = append(retVal, p)
				}
			} else {
				p := request.Payload{Kind: k}
				if kind == SECRETS {
					p.Content = arrToBase(v)
				} else if kind == CONFIGS {
					p.Content = v
				}
				retVal = append(retVal, p)
			}
		}
	} else {
		return nil
	}

	return retVal
}

func extractPayload(file, region, cluster map[string][]string, kind string) []request.Payload {
	p := constructPayload(file, kind)
	if p != nil {
		return p
	}

	p = constructPayload(region, kind)
	if p != nil {
		return p
	}

	p = constructPayload(cluster, kind)
	if p != nil {
		return p
	}

	return nil
}

func constructStrategy(strategy map[string]string) *request.Strategy {
	s := request.Strategy{}
	if len(strategy) != 0 {
		if val, ok := strategy[TYPE]; ok {
			s.Type = val
		} else {
			return nil //TODO: should throw some error
		}

		if val, ok := strategy[UPDATE]; ok {
			s.Kind = val
		} else {
			return nil //TODO: should throw some error
		}
	} else {
		return nil
	}

	return &s
}

func extractStrategy(file, region, cluster map[string]string) *request.Strategy {
	s := constructStrategy(file)
	if s != nil {
		return s
	}

	s = constructStrategy(region)
	if s != nil {
		return s
	}

	s = constructStrategy(cluster)
	if s != nil {
		return s
	}

	return nil
}

func constructSelector(selector map[string]map[string]string) *request.Selector {
	s := request.Selector{}
	if len(selector) != 0 {
		if val, ok := selector[LABELS]; ok {
			s.Labels = val
		} else {
			return nil
		}

		if val, ok := selector[COMPARE]; ok {
			s.Compare = val
		} else {
			s.Compare = map[string]string{KIND: ALL}
		}
	} else {
		return nil
	}

	return &s
}

func extractSelector(file, region, cluster map[string]map[string]string) *request.Selector {
	s := constructSelector(file)
	if s != nil {
		return s
	}

	s = constructSelector(region)
	if s != nil {
		return s
	}

	s = constructSelector(cluster)
	if s != nil {
		return s
	}

	return nil
}
