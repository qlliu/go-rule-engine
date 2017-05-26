package go_rule_engine

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"github.com/fatih/structs"
	"errors"
	"strconv"
)

func NewRulesWithJsonAndLogic(jsonStr []byte, logic string) (*Rules, error) {
	if logic == "" {
		// empty logic
		return NewRulesWithJson(jsonStr)
	}

	formatLogic := formatLogicExpression(logic)
	// validate the formatLogic string
	// 1. only contain legal symbol
	isValidSymbol := isFormatLogicExpressionAllValidSymbol(formatLogic)
	if (!isValidSymbol) {
		return nil, errors.New("invalid logic expression: invalid symbol")
	}

	rules, err := NewRulesWithJson(jsonStr)
	if err != nil {
		return nil, err
	}
	// 2. all ids in logic string must be in rules ids
	isValidIds := isFormatLogicExpressionAllIdsExist(formatLogic, rules)
	if (!isValidIds) {
		return nil, errors.New("invalid logic expression: invalid id")
	}

	rules.Logic = formatLogic

	return rules, nil
}

func NewRulesWithJson(jsonStr []byte) (*Rules, error) {
	var rules []*Rule
	err := json.Unmarshal(jsonStr, &rules)
	if err != nil {
		return nil, err
	}
	return NewRulesWithArray(rules), nil
}

func NewRulesWithArray(rules []*Rule) *Rules {
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
	}
}

func (rs *Rules) Fit(o interface{}) (bool, map[int]string) {
	m := structs.Map(o)
	return rs.FitWithMap(m)
}

func (rs *Rules) FitWithMap(o map[string]interface{}) (bool, map[int]string) {
	var results = make(map[int]bool)
	var tips = make(map[int]string)
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
				return false, nil
			}
		}
		flag := rule.fit(v)
		results[rule.Id] = flag
		if (!flag) {
			// unfit, record msg
			tips[rule.Id] = rule.Msg
		}
	}
	// compute result by considering logic
	if !hasLogic {
		for _, flag := range results {
			if !flag {
				return false, tips
			}
		}
		return true, nil
	} else {
		answer, err := rs.calculateExpression(rs.Logic, results)
		if err != nil {
			return false, nil
		}
		return answer, tips
	}
}

func (r *Rule) fit(v interface{}) bool {
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

func formatLogicExpression(strRawExpr string) string {
	flagPre := ""
	flagNow := ""
	runesOrigin := []rune(strings.ToLower(strRawExpr))
	runesPretty := make([]rune, 0)
	for _, c := range runesOrigin {
		if c <= []rune("9")[0] && c >= []rune("0")[0] {
			flagNow = "num"
		} else if c <= []rune("z")[0] && c >= []rune("a")[0] {
			flagNow = "char"
		} else if c == []rune("(")[0] || c == []rune(")")[0] {
			flagNow = "bracket"
		} else {
			flagNow = flagPre
		}
		if flagNow != flagPre {
			// should insert space here
			runesPretty = append(runesPretty, []rune(" ")[0])
		}
		runesPretty = append(runesPretty, c)
		flagPre = flagNow
	}
	// remove redundant space
	flagPre = "notSpace"
	flagNow = ""
	runesTrim := make([]rune, 0)
	for _, c := range runesPretty {
		if c == []rune(" ")[0] {
			flagNow = "space"
		} else {
			flagNow = "notSpace"
		}
		if flagNow == "space" && flagPre == "space" {
			// continuous space
			continue
		} else {
			runesTrim = append(runesTrim, c)
		}
		flagPre = flagNow
	}
	strPrettyTrim := string(runesTrim)
	strPrettyTrim = strings.Trim(strPrettyTrim, " ")

	return strPrettyTrim
}

func isFormatLogicExpressionAllValidSymbol(strFormatLogic string) bool {
	listSymbol := strings.Split(strFormatLogic, " ")
	for _, symbol := range listSymbol {
		flag := false
		pattern := "^\\d*$"
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return false
		}
		if regex.MatchString(symbol) {
			// is number ok
			continue
		}
		for _, op := range VALID_OPERATORS {
			if op == symbol {
				// is operator ok
				flag = true
			}
		}
		for _, v := range []string{"(", ")"} {
			if v == symbol {
				// is bracket ok
				flag = true
			}
		}
		if (!flag) {
			return false
		}
	}
	return true
}

func isFormatLogicExpressionAllIdsExist(strFormatLogic string, rules *Rules) bool {
	mapExistIds := make(map[string]bool)
	for _, eachRule := range rules.Rules {
		mapExistIds[strconv.Itoa(eachRule.Id)] = true
	}
	listSymbol := strings.Split(strFormatLogic, " ")
	pattern := "^\\d*$"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	for _, symbol := range listSymbol {
		if regex.MatchString(symbol) {
			// is id, check it
			if _, ok := mapExistIds[symbol]; ok {
				continue
			} else {
				return false;
			}
		}
	}
	return true;
}
