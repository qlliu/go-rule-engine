package go_rule_engine

import (
	"encoding/json"
	"strings"
)

func NewRulesWithJson(jsonStr []byte) (*Rules, error) {
	var rules []*Rule
	err := json.Unmarshal(jsonStr, &rules)
	if err != nil {
		return nil, err
	}
	return &Rules{
		Rules: rules,
	}, nil
}

func (r *Rules) Fit(o map[string]interface{}) {

}

func pluck(key string, o map[string]interface{}) interface{} {
	if o == nil || key == "" {
		return nil
	}
	paths := strings.Split(key, ".")
	var ok bool
	for index, step := range paths {
		// last step is object key
		if index == len(paths) - 1 {
			return o[step]
		}
		// explore deeper
		if o, ok = o[step].(map[string]interface{}); !ok {
			return nil
		}
	}
	return nil
}