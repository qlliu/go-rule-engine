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

// Fit 匹配, 返回true/false, 如果为false的一些提示
func (m *Molecule) Fit(o interface{}) (bool, map[int]string) {
	rs := (*Rules)(m)
	return rs.Fit(o)
}

// FitWithMap Rules匹配map
func (rs *Rules) FitWithMap(o map[string]interface{}) (bool, map[int]string) {
	fit, tips, _ := rs.fitWithMapInFact(o)
	return fit, tips
}

// FitWithMap Rules匹配map, 返回true/false, 如果为false的一些提示
func (m *Molecule) FitWithMap(o map[string]interface{}) (bool, map[int]string) {
	rs := (*Rules)(m)
	return rs.FitWithMap(o)
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

// Fit RulesSet's fit, means hitting first time in array
func (rst *RulesSet) Fit(o interface{}) *Rules {
	m := structs.Map(o)
	return rst.FitWithMap(m)
}

// FitWithMap RulesSet's fit, means hitting first time in array
func (rst *RulesSet) FitWithMap(o map[string]interface{}) *Rules {
	for _, rs := range rst.RulesSet {
		if flag, _ := rs.FitWithMap(o); flag {
			return rs
		}
	}
	return nil
}

// Fit Compound's fit, means hitting first time in array
func (c *Compound) Fit(o interface{}) *Molecule {
	m := structs.Map(o)
	return c.FitWithMap(m)
}

// FitWithMap Compound's fit, means hitting first time in array
func (c *Compound) FitWithMap(o map[string]interface{}) *Molecule {
	for _, rs := range c.RulesSet {
		if flag, _ := rs.FitWithMap(o); flag {
			return (*Molecule)(rs)
		}
	}
	return nil
}

// FitGetStr Compound return hit, string
func (c *Compound) FitGetStr(o interface{}) (bool, string) {
	m := c.Fit(o)
	if m == nil {
		return false, EmptyStr
	}
	if str, ok := m.Val.(string); ok {
		return true, str
	}
	return true, EmptyStr
}

// FitGetNum Compound return hit, float64
func (c *Compound) FitGetNum(o interface{}) (bool, float64) {
	m := c.Fit(o)
	if m == nil {
		return false, EmptyFloat64
	}
	if num, ok := m.Val.(float64); ok {
		return true, num
	}
	return true, EmptyFloat64
}

// CheckLogicExpressionAndFormat 检查逻辑表达式正确性，并返回formatted
func CheckLogicExpressionAndFormat(logic string) (string, error) {
	return validLogic(logic)
}
