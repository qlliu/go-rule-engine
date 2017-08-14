package ruler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRulesWithJSON(t *testing.T) {
	jsonStr := []byte(`[{"op": "=", "key": "status", "val": 1}]`)
	rules, err := newRulesWithJSON(jsonStr)
	if err != nil {
		t.Error(err)
	}
	rule := &Rule{
		Op:  "=",
		Key: "status",
		Val: float64(1),
		ID:  2,
	}
	assert.Equal(t, rule, rules.Rules[0])
}

func TestPluck(t *testing.T) {
	obj := map[string]interface{}{"op": map[string]interface{}{"deep": 1}, "key": "status", "val": 1}
	result := pluck("op.deep", obj)
	assert.Equal(t, 1, result)
}

func TestRule_Fit(t *testing.T) {
	rule := &Rule{
		Op:  "=",
		Key: "status",
		Val: 0,
	}
	result := rule.fit(0)
	assert.True(t, result)
}

func TestRules_Fit(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "status", "val": "abcd", "id": 13, "msg": "状态不对"},
	{"op": "=", "key": "name", "val": "peter", "id": 15},
	{"op": ">=", "key": "data.deep", "val": 1, "id": 17, "msg": "deep 数值不对"}
	]`)
	rules, err := newRulesWithJSON(jsonStr)
	if err != nil {
		t.Error(err)
	}
	rules.Logic = "( 15 or 13 ) and 17 and not 13"

	obj := map[string]interface{}{"data": map[string]interface{}{"deep": 1}, "name": "peter", "status": "abc"}
	result, _ := rules.FitWithMap(obj)
	assert.True(t, result)
}

func TestRules_Fit2(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	rules, err := newRulesWithJSON(jsonStr)
	if err != nil {
		t.Error(err)
	}
	rules.Logic = "( 15 or 13 ) and 17 and not 13"
	type B struct {
		Deep int
	}
	type A struct {
		Data   B
		Name   string
		Status string
	}
	obj := A{
		Data: B{
			Deep: 1,
		},
		Name:   "peter",
		Status: "abc",
	}
	result, _ := rules.Fit(obj)
	assert.True(t, result)
}

func TestNewRulesWithJSONAndLogic(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	logic := "     13       and (15     )"
	rules, err := NewRulesWithJSONAndLogic(jsonStr, logic)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "13 and ( 15 )", rules.Logic)
}

func TestNewRulesWithJSONAndLogic2(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	logic := "     13    or   (15     )"
	rules, err := NewRulesWithJSONAndLogic(jsonStr, logic)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "13 or ( 15 )", rules.Logic)
}

func TestNewRulesWithJSONAndLogic3(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	logic := "     13     and  (15or13    )"
	rules, err := NewRulesWithJSONAndLogic(jsonStr, logic)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "13 and ( 15 or 13 )", rules.Logic)
}

func TestNewRulesWithArrayAndLogic(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	logic := "     13     and  (15or13    )"
	rules, err := NewRulesWithJSONAndLogic(jsonStr, logic)
	if err != nil {
		t.Error(err)
	}
	rules, err = NewRulesWithArrayAndLogic(rules.Rules, logic)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "13 and ( 15 or 13 )", rules.Logic)
}

func TestNewRulesWithJSONAndLogicAndInfo(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	logic := "     13     and  (15or13    )"
	extractInfo := map[string]string{
		"name": "名称",
		"msg":  "提示",
	}
	rules, err := NewRulesWithJSONAndLogicAndInfo(jsonStr, logic, extractInfo)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "名称", rules.Name)
}

func TestNewRulesWithArrayAndLogicAndInfo(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	logic := "     13     and  (15or13    )"
	extractInfo := map[string]string{
		"name": "名称",
		"msg":  "提示",
	}
	rules, err := NewRulesWithJSONAndLogic(jsonStr, logic)
	if err != nil {
		t.Error(err)
	}
	rules, err = NewRulesWithArrayAndLogicAndInfo(rules.Rules, logic, extractInfo)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "提示", rules.Msg)
}

func TestNewRulesSet(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	logic := "     13     and  (15or13    )"
	extractInfo := map[string]string{
		"name": "",
		"msg":  "提示",
	}
	rules, err := NewRulesWithJSONAndLogic(jsonStr, logic)
	if err != nil {
		t.Error(err)
	}
	rules, err = NewRulesWithArrayAndLogicAndInfo(rules.Rules, logic, extractInfo)
	if err != nil {
		t.Error(err)
	}

	rulesSet := NewRulesSet([]*Rules{rules}, extractInfo)
	assert.Equal(t, rulesSet.RulesSet[0].Name, "1")
}

func TestRulesSet_FitSetWithMap(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Key", "val": 1, "id": 17}
	]`)
	logic := "     "
	extractInfo := map[string]string{
		"name": "",
		"msg":  "提示",
	}
	rules, err := NewRulesWithJSONAndLogicAndInfo(jsonStr, logic, extractInfo)
	if err != nil {
		t.Error(err)
	}

	rulesSet := NewRulesSet([]*Rules{rules}, extractInfo)

	obj := map[string]interface{}{"Name": "peter", "Status": "abcd", "Key": 0}
	fitRules, _ := rules.FitWithMap(obj)
	assert.False(t, fitRules)

	result := rulesSet.FitSetWithMap(obj)
	assert.Nil(t, result)
}

func TestRules_FitWithMap(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "=", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Key", "val": 1, "id": 17}
	]`)
	logic := "13 and 15"
	extractInfo := map[string]string{
		"name": "",
		"msg":  "提示",
	}
	rules, err := NewRulesWithJSONAndLogicAndInfo(jsonStr, logic, extractInfo)
	if err != nil {
		t.Error(err)
	}
	objMap := map[string]interface{}{"Status": "abcd"}
	fit, _ := rules.FitWithMap(objMap)
	assert.False(t, fit)
}

func TestRules_FitWithMapAskVal(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "=", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Key", "val": 1, "id": 17}
	]`)
	logic := "13"
	extractInfo := map[string]string{
		"name": "",
		"msg":  "提示",
	}
	rules, err := NewRulesWithJSONAndLogicAndInfo(jsonStr, logic, extractInfo)
	if err != nil {
		t.Error(err)
	}
	objMap := map[string]interface{}{"Status": "abcd"}
	fit, _, val := rules.FitWithMapAskVal(objMap)
	valExpect := map[int]interface{}{17: nil, 13: "abcd", 15: nil}
	assert.True(t, fit)
	assert.Equal(t, valExpect, val)
}
