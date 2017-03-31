package go_rule_engine

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
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
			maxId++
			rules[index].Id = maxId
		}
	}
	return &Rules{
		Rules: rules,
	}, nil
}

func (rs *Rules) Fit(o map[string]interface{}) bool {
	var results = make(map[int]bool)
	var hasLogic = false
	if rs.Logic != "" {
		hasLogic = true
	}
	for _, rule := range rs.Rules {
		v := pluck(rule.Key, o)
		if v != nil && rule.Val != nil {
			typeV := reflect.TypeOf(v)
			typeR := reflect.TypeOf(rule.Val)
			if !typeV.Comparable() || !typeR.Comparable() {
				return false
			}
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

func (r *Rule) Fit(v interface{}) bool {
	op := r.Op
	// judge if need convert to uniform type
	var ok bool
	// index-0 actual, index-1 expect
	var pairStr = make([]string, 2)
	var pairNum = make([]float64, 2)
	var isStr, isNum bool
	pairStr[0], ok = v.(string)
	if !ok {
		pairNum[0] = formatNumber(v)
		isStr = false
		isNum = true
	} else {
		isStr = true
		isNum = false
	}
	pairStr[1], ok = r.Val.(string)
	if !ok {
		pairNum[1] = formatNumber(r.Val)
		isStr = false
	} else {
		isNum = false
	}
	// if types different
	if !isStr && !isNum {
		return false
	}

	switch op {
	case "=", "eq":
		if isNum {
			return pairNum[0] == pairNum[1]
		}
		if isStr {
			return pairStr[0] == pairStr[1]
		}
		return false
	case ">", "gt":
		if isNum {
			return pairNum[0] > pairNum[1]
		}
		if isStr {
			return pairStr[0] > pairStr[1]
		}
		return false
	case "<", "lt":
		if isNum {
			return pairNum[0] < pairNum[1]
		}
		if isStr {
			return pairStr[0] < pairStr[1]
		}
		return false
	case ">=", "gte":
		if isNum {
			return pairNum[0] >= pairNum[1]
		}
		if isStr {
			return pairStr[0] >= pairStr[1]
		}
		return false
	case "<=", "lte":
		if isNum {
			return pairNum[0] <= pairNum[1]
		}
		if isStr {
			return pairStr[0] <= pairStr[1]
		}
		return false
	case "!=", "neq":
		if isNum {
			return pairNum[0] != pairNum[1]
		}
		if isStr {
			return pairStr[0] != pairStr[1]
		}
		return false
	case "@", "contain":
		return checkRegex(pairStr[1], pairStr[0])
	case "!@", "ncontain":
		return !checkRegex(pairStr[1], pairStr[0])
	case "^$", "regex":
		return checkRegex(pairStr[1], pairStr[0])
	case "0", "empty":
		return v == nil
	case "1", "nempty":
		return v != nil
	default:
		return false
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
		if index == len(paths)-1 {
			return o[step]
		}
		// explore deeper
		if o, ok = o[step].(map[string]interface{}); !ok {
			return nil
		}
	}
	return nil
}

func formatNumber(v interface{}) float64 {
	switch t := v.(type) {
	case uint:
		return float64(t)
	case uint8:
		return float64(t)
	case uint16:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)
	case int:
		return float64(t)
	case int8:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case float32:
		return float64(t)
	case float64:
		return float64(t)
	default:
		return 0
	}
}

func checkRegex(pattern, o string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return regex.MatchString(o)
}
