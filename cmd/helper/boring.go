package helper

import (
	"github.com/c12s/cockpit/cmd/model/request"
)

func constructPayload(payload map[string][]string) []request.Payload {
	retVal := []request.Payload{}
	if len(payload) != 0 {
		for k, v := range payload {
			if inside(k) {
				p := request.Payload{}
				if k == "file" {
					//TODO: Open file, read file and construct ENV variables
					p.Kind = constructFileKey(v[0])
					p.Content = v
				} else {
					p.Kind = k
					p.Content = v
				}
				retVal = append(retVal, p)
			} else {

			}
		}
	} else {
		return nil
	}

	return retVal
}

func extractPayload(file, region, cluster map[string][]string) []request.Payload {
	p := constructPayload(file)
	if p != nil {
		return p
	}

	p = constructPayload(region)
	if p != nil {
		return p
	}

	p = constructPayload(cluster)
	if p != nil {
		return p
	}

	return nil
}

func constructStrategy(strategy map[string]string) *request.Strategy {
	s := request.Strategy{}
	if len(strategy) != 0 {
		if val, ok := strategy["type"]; ok {
			s.Type = val
		} else {
			return nil //TODO: should throw some error
		}

		if val, ok := strategy["update"]; ok {
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
		if val, ok := selector["labels"]; ok {
			s.Labels = val
		} else {
			return nil
		}

		if val, ok := selector["compare"]; ok {
			s.Compare = val
		} else {
			s.Compare = map[string]string{"kind": "all"}
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
