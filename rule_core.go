package go_rule_engine

import (
	"encoding/json"
	"strings"
	"reflect"
)

func NewRulesWithJson(jsonStr []byte) (*Rules, error) {
	var rules []*Rule
	err := json.Unmarshal(jsonStr, &rules)
	if err != nil {
		return nil, err
	}
	// give rule an id
	var maxId = 1
	for _, rule := range rules {
		if rule.Id > maxId {
			maxId = rule.Id
		}
	}
	for index := range rules {
		if rules[index].Id == 0 {
			rules[index].Id = maxId
			maxId ++
		}
	}
	return &Rules{
		Rules: rules,
	}, nil
}

func (rs *Rules) Fit(o map[string]interface{}) bool {
	var results = make(map[int]bool)
	var hasLogic = false
	for _, rule := range rs.Rules {
		v := pluck(rule.Key, o)
		typeV := reflect.TypeOf(v)
		typeR := reflect.TypeOf(rule.Val)
		if !typeV.Comparable() || !typeR.Comparable() {
			return false
		}
		// seek logic
		if rule.Logic != "" {
			hasLogic = true
		}
		flag := rule.Fit(v)
		results[rule.Id] = flag
	}

	// compute result by considering logic
	var answer = true
	if !hasLogic {
		for _, flag := range results {
			answer = flag && answer
			if !answer {
				return false
			}
		}
	}

	return answer
}

func (r *Rule) Fit(actual interface{}) bool {
	op := r.Op
	expect := r.Val
	switch op {
	case "=":
		return expect == actual
	}
	return false
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