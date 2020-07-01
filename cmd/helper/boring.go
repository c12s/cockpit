package helper

import (
	"bufio"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"github.com/c12s/cockpit/cmd/model"
	"github.com/c12s/cockpit/cmd/model/request"
	"gopkg.in/yaml.v2"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func tobytes(n int64, unit string) (int64, error) {
	if val, ok := maper[unit]; ok {
		return n * val, nil
	}
	return 0, errors.New(fmt.Sprintf("%s not valid unit.valid units are b,kb,mb,gb,tb", unit))
}

func bytestostring(n int64, unit string) (string, error) {
	b, err := tobytes(n, unit)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(b, 10), nil

}

func PrepareLabels(l string) (string, error) {
	parts := []string{}
	re := regexp.MustCompile("([0-9]+)([a-z]+)")
	for _, item := range strings.Split(l, ",") {
		part := strings.Split(item, ":")
		if part[0] == "memory" || part[0] == "storage" {
			m := re.FindStringSubmatch(part[1])
			n, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				return "", err
			}
			newVal, err := bytestostring(n, m[2])
			if err != nil {
				return "", err
			}
			parts = append(parts, strings.Join([]string{part[0], newVal}, ":"))
		} else {
			parts = append(parts, item)
		}
	}
	return strings.Join(parts, ","), nil
}

func toBASE64(data string) string {
	b64Data := b64.StdEncoding.EncodeToString([]byte(data))
	return b64Data
}

func arrToBase(data yaml.MapSlice) map[string]string {
	conv := map[string]string{}
	for _, valItem := range data {
		conv[valItem.Key.(string)] = toBASE64(valItem.Value.(string))
	}
	return conv
}

func readFile(file, kind string) (error, map[string]string) {
	// Open the file.
	f, err := os.Open(file)
	if err != nil {
		return err, nil
	}

	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	data := map[string]string{}
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && len(line) > 0 {
			if !strings.Contains(line, "=") {
				return errors.New(fmt.Sprintf("Error: Line %s is not in right format KEY=VALUE.", line)), nil
			} else {
				bLine := strings.Split(line, "=")
				if kind == SECRETS {
					data[bLine[0]] = toBASE64(bLine[1])
				} else if kind == CONFIGS {
					data[bLine[0]] = bLine[1]
				}
			}
		}
	}

	return nil, data
}

func convert(data yaml.MapSlice) map[string]string {
	conv := map[string]string{}
	for _, valItem := range data {
		conv[valItem.Key.(string)] = valItem.Value.(string)
	}
	return conv
}

func constructPayload(payload map[string]yaml.MapSlice, kind string) []request.Payload {
	retVal := []request.Payload{}
	if len(payload) != 0 {
		for k, v := range payload {
			if k == FILES {
				for _, valItem := range v {
					if k == FILES {
						err, data := readFile(valItem.Value.(string), kind)
						if err != nil {
							return nil
						}

						data[FILE_NAME] = strings.Join([]string{valItem.Key.(string), "env"}, ".")
						p := request.Payload{
							Kind:    FILE,
							Content: data,
						}
						retVal = append(retVal, p)
					}
				}
			} else {
				p := request.Payload{Kind: k}
				if kind == SECRETS {
					p.Content = arrToBase(v)
				} else if kind == CONFIGS {
					p.Content = convert(v)
				} else if kind == ACTIONS {
					p.Content = convert(v)
					p.Index = []string{}

					// to preserve order of the operations
					// we need to know isertation time
					for _, item := range v {
						p.Index = append(p.Index, item.Key.(string))
					}
				}
				retVal = append(retVal, p)
			}
		}
	} else {
		return nil
	}
	return retVal
}

func extractPayload(file, region, cluster map[string]yaml.MapSlice, kind string) []request.Payload {
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

func constructStrategy(strategy *model.Strategy) *request.Strategy {
	s := request.Strategy{}
	if strategy == nil {
		return nil
	}
	// if len(strategy) != 0 {
	// 	if val, ok := strategy[TYPE]; ok {
	// 		s.Type = val
	// 	} else {
	// 		return nil //TODO: should throw some error
	// 	}

	// 	if val, ok := strategy[UPDATE]; ok {
	// 		s.Kind = val
	// 	} else {
	// 		return nil //TODO: should throw some error
	// 	}

	// 	if val, ok := strategy[INTERVAL]; ok {
	// 		s.Interval = val
	// 	} else {
	// 		return nil //TODO: should throw some error
	// 	}
	// } else {
	// 	return nil
	// }
	s.Type = strategy.Type
	s.Kind = strategy.Update
	s.Interval = strategy.Interval
	s.Retry = strategy.Retry

	return &s
}

func extractStrategy(file, region, cluster *model.Strategy) *request.Strategy {
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

func constructSelector(selector map[string]yaml.MapSlice) *request.Selector {
	s := request.Selector{}
	if len(selector) != 0 {
		if val, ok := selector[LABELS]; ok {
			s.Labels = convert(val)
		} else {
			return nil
		}

		if val, ok := selector[COMPARE]; ok {
			s.Compare = convert(val)
		} else {
			s.Compare = map[string]string{KIND: ALL}
		}
	} else {
		return nil
	}

	return &s
}

func extractSelector(file, region, cluster map[string]yaml.MapSlice) *request.Selector {
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
