package go_rule_engine

import "github.com/fatih/structs"

/**
  用json串构造Rules的完全方法，logic表达式如果没有则传空字符串, ["name": "规则名称", "msg": "规则不符合的提示"]
*/
func NewRulesWithJsonAndLogicAndInfo(jsonStr []byte, logic string, extractInfo map[string]string) (*Rules, error) {
	rulesObj, err := NewRulesWithJsonAndLogic(jsonStr, logic)
	if err != nil {
		return nil, err
	}
	return injectExtractInfo(rulesObj, extractInfo), nil
}

/**
  用rule数组构造Rules的完全方法，logic表达式如果没有则传空字符串, ["name": "规则名称", "msg": "规则不符合的提示"]
*/
func NewRulesWithArrayAndLogicAndInfo(rules []*Rule, logic string, extractInfo map[string]string) (*Rules, error) {
	rulesObj, err := NewRulesWithArrayAndLogic(rules, logic)
	if err != nil {
		return nil, err
	}
	return injectExtractInfo(rulesObj, extractInfo), nil
}

/**
  用json串构造Rules的标准方法，logic表达式如果没有则传空字符串
*/
func NewRulesWithJsonAndLogic(jsonStr []byte, logic string) (*Rules, error) {
	if logic == "" {
		// empty logic
		return newRulesWithJson(jsonStr)
	}
	rulesObj, err := newRulesWithJson(jsonStr)
	if err != nil {
		return nil, err
	}
	rulesObj, err = injectLogic(rulesObj, logic)
	if err != nil {
		return nil, err
	}

	return rulesObj, nil
}

/**
  用rule数组构造Rules的标准方法，logic表达式如果没有则传空字符串
*/
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

/**
  Rules匹配传入结构体
*/
func (rs *Rules) Fit(o interface{}) (bool, map[int]string) {
	m := structs.Map(o)
	return rs.FitWithMap(m)
}

/**
  Rules匹配map
*/
func (rs *Rules) FitWithMap(o map[string]interface{}) (bool, map[int]string) {
	fit, tips, _ := rs.fitWithMapInFact(o)
	return fit, tips
}

/**
  Rules匹配结构体，同时返回所有子规则key值
*/
func (rs *Rules) FitAskVal(o interface{}) (bool, map[int]string, map[int]interface{}) {
	m := structs.Map(o)
	return rs.FitWithMapAskVal(m)
}

/**
  Rules匹配map，同时返回所有子规则key值
*/
func (rs *Rules) FitWithMapAskVal(o map[string]interface{}) (bool, map[int]string, map[int]interface{}) {
	return rs.fitWithMapInFact(o)
}
