package ruler

import (
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"math"

	"github.com/fatih/structs"
)

func injectLogic(rules *Rules, logic string) (*Rules, error) {
	formatLogic := formatLogicExpression(logic)
	if formatLogic == Space || formatLogic == EmptyStr {
		return rules, nil
	}
	// validate the formatLogic string
	// 1. only contain legal symbol
	isValidSymbol := isFormatLogicExpressionAllValidSymbol(formatLogic)
	if !isValidSymbol {
		return nil, errors.New("invalid logic expression: invalid symbol")
	}

	// 2. check logic expression by trying to  calculate result with random bool
	err := tryToCalculateResultByFormatLogicExpressionWithRandomProbe(formatLogic)
	if err != nil {
		return nil, errors.New("invalid logic expression: can not calculate")
	}

	// 3. all ids in logic string must be in rules ids
	isValidIds := isFormatLogicExpressionAllIdsExist(formatLogic, rules)
	if !isValidIds {
		return nil, errors.New("invalid logic expression: invalid id")
	}
	rules.Logic = formatLogic

	return rules, nil
}

func injectExtractInfo(rules *Rules, extractInfo map[string]string) *Rules {
	if name, ok := extractInfo["name"]; ok {
		rules.Name = name
	}
	if msg, ok := extractInfo["msg"]; ok {
		rules.Msg = msg
	}
	return rules
}

// NewRulesSet RulesSet的构造方法，["name": "规则集的名称", "msg": "规则集的简述"]
func NewRulesSet(listRules []*Rules, extractInfo map[string]string) *RulesSet {
	// check if every rules has name, if not give a index as name
	for index, rules := range listRules {
		if rules.Name == EmptyStr {
			rules.Name = strconv.Itoa(index + 1)
		}
	}
	name := extractInfo["name"]
	msg := extractInfo["msg"]
	return &RulesSet{
		RulesSet: listRules,
		Name:     name,
		Msg:      msg,
	}
}

func newRulesWithJSON(jsonStr []byte) (*Rules, error) {
	var rules []*Rule
	err := json.Unmarshal(jsonStr, &rules)
	if err != nil {
		return nil, err
	}
	return newRulesWithArray(rules), nil
}

func newRulesWithArray(rules []*Rule) *Rules {
	// give rule an id
	var maxID = 1
	for _, rule := range rules {
		if rule.ID > maxID {
			maxID = rule.ID
		}
	}
	for index := range rules {
		if rules[index].ID == 0 {
			maxID++
			rules[index].ID = maxID
		}
	}
	return &Rules{
		Rules: rules,
	}
}

// FitSet RulesSet匹配结构体
func (rst *RulesSet) FitSet(o interface{}) []string {
	m := structs.Map(o)
	return rst.FitSetWithMap(m)
}

// FitSetWithMap RulesSet匹配Map
func (rst *RulesSet) FitSetWithMap(o map[string]interface{}) []string {
	result := make([]string, 0)
	for _, rules := range rst.RulesSet {
		if fit, _ := rules.FitWithMap(o); fit {
			// hit this rules
			result = append(result, rules.Name)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func (rs *Rules) fitWithMapInFact(o map[string]interface{}) (bool, map[int]string, map[int]interface{}) {
	var results = make(map[int]bool)
	var tips = make(map[int]string)
	var values = make(map[int]interface{})
	var hasLogic = false
	var allRuleIDs []int
	if rs.Logic != EmptyStr {
		hasLogic = true
	}
	for _, rule := range rs.Rules {
		v := pluck(rule.Key, o)
		if v != nil && rule.Val != nil {
			typeV := reflect.TypeOf(v)
			typeR := reflect.TypeOf(rule.Val)
			if !typeV.Comparable() || !typeR.Comparable() {
				return false, nil, nil
			}
		}
		values[rule.ID] = v

		flag := rule.fit(v)
		results[rule.ID] = flag
		if !flag {
			// fit false, record msg, for no logic expression usage
			tips[rule.ID] = rule.Msg
		}
		allRuleIDs = append(allRuleIDs, rule.ID)
	}
	// compute result by considering logic

	if !hasLogic {
		for _, flag := range results {
			if !flag {
				return false, tips, values
			}
		}
		return true, rs.getTipsByRuleIDs(allRuleIDs), values
	}
	answer, ruleIDs, err := rs.calculateExpressionByTree(results)
	// tree can return fail reasons in fact
	tips = rs.getTipsByRuleIDs(ruleIDs)
	if err != nil {
		return false, nil, values
	}
	return answer, tips, values
}

func (rs *Rules) getTipsByRuleIDs(ids []int) map[int]string {
	var tips = make(map[int]string)
	var allTips = make(map[int]string)
	for _, rule := range rs.Rules {
		allTips[rule.ID] = rule.Msg
	}
	for _, id := range ids {
		tips[id] = allTips[id]
	}
	return tips
}

func (r *Rule) fit(v interface{}) bool {
	op := r.Op
	// judge if need convert to uniform type
	var ok bool
	// index-0 actual, index-1 expect
	var pairStr = make([]string, 2)
	var pairNum = make([]float64, 2)
	var isStr, isNum, isObjStr, isRuleStr bool
	pairStr[0], ok = v.(string)
	if !ok {
		pairNum[0] = formatNumber(v)
		isStr = false
		isNum = true
		isObjStr = false
	} else {
		isStr = true
		isNum = false
		isObjStr = true
	}
	pairStr[1], ok = r.Val.(string)
	if !ok {
		pairNum[1] = formatNumber(r.Val)
		isStr = false
		isRuleStr = false
	} else {
		isNum = false
		isRuleStr = true
	}

	var flagOpIn bool
	// if in || nin
	if op == "@" || op == "in" || op == "!@" || op == "nin" || op == "<<" || op == "between" {
		flagOpIn = true
		if !isObjStr && isRuleStr {
			pairStr[0] = strconv.FormatFloat(pairNum[0], 'f', -1, 64)
		}
	}

	// if types different, ignore in & nin
	if !isStr && !isNum && !flagOpIn {
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
	case "@", "in":
		return isIn(pairStr[0], pairStr[1], !isObjStr)
	case "!@", "nin":
		return !isIn(pairStr[0], pairStr[1], !isObjStr)
	case "^$", "regex":
		return checkRegex(pairStr[1], pairStr[0])
	case "0", "empty":
		return v == nil
	case "1", "nempty":
		return v != nil
	case "<<", "between":
		return isBetween(pairNum[0], pairStr[1])
	case "@@", "intersect":
		return isIntersect(pairStr[1], pairStr[0])
	default:
		return false
	}
}

func pluck(key string, o map[string]interface{}) interface{} {
	if o == nil || key == EmptyStr {
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
		return t
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
	var flagPre, flagNow string
	strBracket := "bracket"
	strSpace := "space"
	strNotSpace := "notSpace"
	strOrigin := strings.ToLower(strRawExpr)
	runesPretty := make([]rune, 0)

	for _, c := range strOrigin {
		if c <= '9' && c >= '0' {
			flagNow = "num"
		} else if c <= 'z' && c >= 'a' {
			flagNow = "char"
		} else if c == '(' || c == ')' {
			flagNow = strBracket
		} else {
			flagNow = flagPre
		}
		if flagNow != flagPre || flagNow == strBracket && flagPre == strBracket {
			// should insert space here
			runesPretty = append(runesPretty, []rune(Space)[0])
		}
		runesPretty = append(runesPretty, c)
		flagPre = flagNow
	}
	// remove redundant space
	flagPre = strNotSpace
	runesTrim := make([]rune, 0)
	for _, c := range runesPretty {
		if c == []rune(Space)[0] {
			flagNow = strSpace
		} else {
			flagNow = strNotSpace
		}
		if flagNow == strSpace && flagPre == strSpace {
			// continuous space
			continue
		} else {
			runesTrim = append(runesTrim, c)
		}
		flagPre = flagNow
	}
	strPrettyTrim := string(runesTrim)
	strPrettyTrim = strings.Trim(strPrettyTrim, Space)

	return strPrettyTrim
}

func isFormatLogicExpressionAllValidSymbol(strFormatLogic string) bool {
	listSymbol := strings.Split(strFormatLogic, Space)
	for _, symbol := range listSymbol {
		flag := false
		regex, err := regexp.Compile(PatternNumber)
		if err != nil {
			return false
		}
		if regex.MatchString(symbol) {
			// is number ok
			continue
		}
		for _, op := range ValidOperators {
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
		if !flag {
			return false
		}
	}
	return true
}

func isFormatLogicExpressionAllIdsExist(strFormatLogic string, rules *Rules) bool {
	mapExistIds := make(map[string]bool)
	for _, eachRule := range rules.Rules {
		mapExistIds[strconv.Itoa(eachRule.ID)] = true
	}
	listSymbol := strings.Split(strFormatLogic, Space)
	regex, err := regexp.Compile(PatternNumber)
	if err != nil {
		return false
	}
	for _, symbol := range listSymbol {
		if regex.MatchString(symbol) {
			// is id, check it
			if _, ok := mapExistIds[symbol]; ok {
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func tryToCalculateResultByFormatLogicExpressionWithRandomProbe(strFormatLogic string) error {
	listSymbol := strings.Split(strFormatLogic, Space)
	regex, err := regexp.Compile(PatternNumber)
	if err != nil {
		return err
	}
	// random probe
	mapProbe := make(map[int]bool)
	for _, symbol := range listSymbol {
		if regex.MatchString(symbol) {
			id, iErr := strconv.Atoi(symbol)
			if iErr != nil {
				return iErr
			}
			randomInt := rand.Intn(10)
			randomBool := randomInt < 5
			mapProbe[id] = randomBool
		}
	}
	// calculate still use reverse_polish_notation
	r := &Rules{}
	_, err = r.calculateExpression(strFormatLogic, mapProbe)
	return err
}

func numOfOperandInLogic(op string) int8 {
	mapOperand := map[string]int8{"or": 2, "and": 2, "not": 1}
	return mapOperand[op]
}

func computeOneInLogic(op string, v []bool) (bool, error) {
	switch op {
	case "or":
		return v[0] || v[1], nil
	case "and":
		return v[0] && v[1], nil
	case "not":
		return !v[0], nil
	default:
		return false, errors.New("unrecognized op")
	}
}

func isIn(needle, haystack string, isNeedleNum bool) bool {
	// get number of needle
	var iNum float64
	var err error
	if isNeedleNum {
		if iNum, err = strconv.ParseFloat(needle, 64); err != nil {
			return false
		}
	}
	// compatible to "1, 2, 3" and "1,2,3"
	li := strings.Split(haystack, ",")
	for _, o := range li {
		trimO := strings.TrimLeft(o, " ")
		if isNeedleNum {
			oNum, err := strconv.ParseFloat(trimO, 64)
			if err != nil {
				continue
			}
			if math.Abs(iNum-oNum) < 1E-5 {
				// 考虑浮点精度问题
				return true
			}
		} else if needle == trimO {
			return true
		}
	}
	return false
}

func isIntersect(objStr string, ruleStr string) bool {
	// compatible to "1, 2, 3" and "1,2,3"
	vl := strings.Split(objStr, ",")
	li := strings.Split(ruleStr, ",")
	for _, o := range li {
		trimO := strings.Trim(o, " ")
		for _, v := range vl {
			trimV := strings.Trim(v, " ")
			if trimV == trimO {
				return true
			}
		}
	}
	return false
}

func isBetween(obj float64, scope string) bool {
	scope = strings.Trim(scope, " ")
	var equalLeft, equalRight bool
	// [] 双闭区间
	result := regexp.MustCompile("^\\[ *(-?\\d*.?\\d*) *, *(-?\\d*.?\\d*) *]$").FindStringSubmatch(scope)
	if len(result) > 2 {
		equalLeft = true
		equalRight = true
		return calculateBetween(obj, result, equalLeft, equalRight)
	}
	// [) 左闭右开区间
	result = regexp.MustCompile("^\\[ *(-?\\d*.?\\d*) *, *(-?\\d*.?\\d*) *\\)$").FindStringSubmatch(scope)
	if len(result) > 2 {
		equalLeft = true
		equalRight = false
		return calculateBetween(obj, result, equalLeft, equalRight)
	}
	// (] 左开右闭区间
	result = regexp.MustCompile("^\\( *(-?\\d*.?\\d*) *, *(-?\\d*.?\\d*) *]$").FindStringSubmatch(scope)
	if len(result) > 2 {
		equalLeft = false
		equalRight = true
		return calculateBetween(obj, result, equalLeft, equalRight)
	}
	// () 双开区间
	result = regexp.MustCompile("^\\( *(-?\\d*.?\\d*) *, *(-?\\d*.?\\d*) *\\)$").FindStringSubmatch(scope)
	if len(result) > 2 {
		equalLeft = false
		equalRight = false
		return calculateBetween(obj, result, equalLeft, equalRight)
	}
	return false
}

func calculateBetween(obj float64, result []string, equalLeft, equalRight bool) bool {
	var hasLeft, hasRight bool
	var left, right float64
	var err error
	if result[1] != "" {
		hasLeft = true
		left, err = strconv.ParseFloat(result[1], 64)
		if err != nil {
			return false
		}
	}
	if result[2] != "" {
		hasRight = true
		right, err = strconv.ParseFloat(result[2], 64)
		if err != nil {
			return false
		}
	}
	// calculate
	if !hasLeft && !hasRight {
		return false
	}
	flag := true
	if hasLeft {
		if equalLeft {
			flag = flag && obj >= left
		} else {
			flag = flag && obj > left
		}
	}
	if hasRight {
		if equalRight {
			flag = flag && obj <= right
		} else {
			flag = flag && obj < right
		}
	}
	return flag
}
