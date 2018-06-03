package ruler

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/structs"
)

// NewRulesWithJSONAndLogicAndInfo 用json串构造Rules的完全方法，logic表达式如果没有则传空字符串, ["name": "规则名称", "msg": "规则不符合的提示"]
func NewRulesWithJSONAndLogicAndInfo(jsonStr []byte, logic string, extractInfo map[string]string) (*Rules, error) {
	rulesObj, err := NewRulesWithJSONAndLogic(jsonStr, logic)
	if err != nil {
		return nil, err
	}
	return injectExtractInfo(rulesObj, extractInfo), nil
}

// NewRulesWithArrayAndLogicAndInfo 用rule数组构造Rules的完全方法，logic表达式如果没有则传空字符串, ["name": "规则名称", "msg": "规则不符合的提示"]
func NewRulesWithArrayAndLogicAndInfo(rules []*Rule, logic string, extractInfo map[string]string) (*Rules, error) {
	rulesObj, err := NewRulesWithArrayAndLogic(rules, logic)
	if err != nil {
		return nil, err
	}
	return injectExtractInfo(rulesObj, extractInfo), nil
}

// NewRulesWithJSONAndLogic 用json串构造Rules的标准方法，logic表达式如果没有则传空字符串
func NewRulesWithJSONAndLogic(jsonStr []byte, logic string) (*Rules, error) {
	if logic == "" {
		// empty logic
		return newRulesWithJSON(jsonStr)
	}
	rulesObj, err := newRulesWithJSON(jsonStr)
	if err != nil {
		return nil, err
	}
	rulesObj, err = injectLogic(rulesObj, logic)
	if err != nil {
		return nil, err
	}

	return rulesObj, nil
}

// NewRulesWithArrayAndLogic 用rule数组构造Rules的标准方法，logic表达式如果没有则传空字符串
func NewRulesWithArrayAndLogic(rules []*Rule, logic string) (*Rules, error) {
	if logic == "" {
		// empty logic
		return newRulesWithArray(rules), nil
	}
	rulesObj := newRulesWithArray(rules)
	rulesObj, err := injectLogic(rulesObj, logic)
	if err != nil {
		return nil, err
	}

	return rulesObj, nil
}

// Fit Rules匹配传入结构体
func (rs *Rules) Fit(o interface{}) (bool, map[int]string) {
	m := structs.Map(o)
	return rs.FitWithMap(m)
}

// FitWithMap Rules匹配map
func (rs *Rules) FitWithMap(o map[string]interface{}) (bool, map[int]string) {
	fit, tips, _ := rs.fitWithMapInFact(o)
	return fit, tips
}

// FitAskVal Rules匹配结构体，同时返回所有子规则key值
func (rs *Rules) FitAskVal(o interface{}) (bool, map[int]string, map[int]interface{}) {
	m := structs.Map(o)
	return rs.FitWithMapAskVal(m)
}

// FitWithMapAskVal Rules匹配map，同时返回所有子规则key值
func (rs *Rules) FitWithMapAskVal(o map[string]interface{}) (bool, map[int]string, map[int]interface{}) {
	return rs.fitWithMapInFact(o)
}

// GetRuleIDsByLogicExpression 根据逻辑表达式得到规则id列表
func GetRuleIDsByLogicExpression(logic string) ([]int, error) {
	var result []int
	formatLogic, errLogic := validLogic(logic)
	if errLogic != nil {
		return nil, errLogic
	}
	if formatLogic == EmptyStr {
		return result, nil
	}

	// return rule id list
	var mapGot = make(map[int]bool)
	listSymbol := strings.Split(formatLogic, Space)
	regex := regexp.MustCompile(PatternNumber)
	for _, symbol := range listSymbol {
		if regex.MatchString(symbol) {
			// is id, check it
			id, err := strconv.Atoi(symbol)
			if err != nil {
				return nil, err
			}
			// keep unique
			if _, ok := mapGot[id]; !ok {
				result = append(result, id)
				mapGot[id] = true
			}
		}
	}
	return result, nil
}

// NewRulesList RulesList的构造方法，["name": "规则集的名称", "msg": "规则集的简述"]
func NewRulesList(listRules []*Rules, extractInfo map[string]string) *RulesList {
	// check if every rules has name, if not give a index as name
	for index, rules := range listRules {
		if rules.Name == EmptyStr {
			rules.Name = strconv.Itoa(index + 1)
		}
	}
	name := extractInfo["name"]
	msg := extractInfo["msg"]
	return &RulesList{
		RulesList: listRules,
		Name:      name,
		Msg:       msg,
	}
}

// Fit RulesList's fit, means hitting first rules in array
func (rst *RulesList) Fit(o interface{}) *Rules {
	m := structs.Map(o)
	return rst.FitWithMap(m)
}

// FitWithMap RulesList's fit, means hitting first rules in array
func (rst *RulesList) FitWithMap(o map[string]interface{}) *Rules {
	for _, rs := range rst.RulesList {
		if flag, _ := rs.FitWithMap(o); flag {
			return rs
		}
	}
	return nil
}

// FitGetStr return hit rules value, string
func (rst *RulesList) FitGetStr(o interface{}) (bool, string) {
	rs := rst.Fit(o)
	if rs == nil {
		return false, EmptyStr
	}
	if str, ok := rs.Val.(string); ok {
		return true, str
	}
	return true, EmptyStr
}

// FitGetFloat64 return hit value, float64
func (rst *RulesList) FitGetFloat64(o interface{}) (bool, float64) {
	rs := rst.Fit(o)
	if rs == nil {
		return false, EmptyFloat64
	}
	return true, formatNumber(rs.Val)
}

// FitGetInt64 return hit value, int64
func (rst *RulesList) FitGetInt64(o interface{}) (bool, int64) {
	rs := rst.Fit(o)
	if rs == nil {
		return false, EmptyFloat64
	}
	var result int64
	switch t := rs.Val.(type) {
	case uint:
		result = int64(t)
	case uint8:
		result = int64(t)
	case uint16:
		result = int64(t)
	case uint32:
		result = int64(t)
	case uint64:
		result = int64(t)
	case int:
		result = int64(t)
	case int8:
		result = int64(t)
	case int16:
		result = int64(t)
	case int32:
		result = int64(t)
	case int64:
		result = t
	case float32:
		result = int64(t + 1e-5)
	case float64:
		result = int64(t + 1e-5)
	}
	return true, result
}

// CheckLogicExpressionAndFormat 检查逻辑表达式正确性，并返回formatted
func CheckLogicExpressionAndFormat(logic string) (string, error) {
	return validLogic(logic)
}
