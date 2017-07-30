package go_rule_engine

import (
	"github.com/docker/docker/pkg/testutil/assert"
	"testing"
)

func TestNewRulesWithJson(t *testing.T) {
	jsonStr := []byte(`[{"op": "=", "key": "status", "val": 1}]`)
	rules, err := NewRulesWithJson(jsonStr)
	if err != nil {
		t.Error(err)
	}
	t.Log(rules.Rules[0])
}

func TestPluck(t *testing.T) {
	obj := map[string]interface{}{"op": map[string]interface{}{"deep": 1}, "key": "status", "val": 1}
	t.Log(obj)
	result := pluck("op.deep", obj)
	t.Log(result)
	t.Log(obj)
}

func TestRule_Fit(t *testing.T) {
	rule := &Rule{
		Op:  "=",
		Key: "status",
		Val: 0,
	}
	result := rule.fit(0)
	t.Log(result)
}

func TestRules_Fit(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "status", "val": "abcd", "id": 13, "msg": "状态不对"},
	{"op": "=", "key": "name", "val": "peter", "id": 15},
	{"op": ">=", "key": "data.deep", "val": 1, "id": 17, "msg": "deep 数值不对"}
	]`)
	rules, err := NewRulesWithJson(jsonStr)
	if err != nil {
		t.Error(err)
	}
	rules.Logic = "( 15 or 13 ) and 17 and not 13"

	obj := map[string]interface{}{"data": map[string]interface{}{"deep": 1}, "name": "peter", "status": "abc"}
	result, msg := rules.FitWithMap(obj)
	t.Log(result)
	t.Log(msg)
}

func TestRules_Fit2(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	rules, err := NewRulesWithJson(jsonStr)
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
	t.Log(result)
}

func TestNewRulesWithJsonAndLogic(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	logic := "     13       and (15     )"
	rules, err := NewRulesWithJsonAndLogic(jsonStr, logic)
	if err != nil {
		t.Error(err)
	}
	t.Log(rules.Rules[0])
	t.Log(rules.Logic)
}

func TestNewRulesWithJsonAndLogic2(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	logic := "     13    or   (15     )"
	rules, err := NewRulesWithJsonAndLogic(jsonStr, logic)
	if err != nil {
		t.Error(err)
	}
	t.Log(rules.Rules[0])
	t.Log(rules.Logic)
}

func TestNewRulesWithJsonAndLogic3(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	logic := "     13     and  (15or13    )"
	rules, err := NewRulesWithJsonAndLogic(jsonStr, logic)
	if err != nil {
		t.Error(err)
	}
	t.Log(rules.Rules[0])
	t.Log(rules.Logic)
}

func TestNewRulesWithArrayAndLogic(t *testing.T) {
	jsonStr := []byte(`[
	{"op": "@", "key": "Status", "val": "abcd", "id": 13},
	{"op": "=", "key": "Name", "val": "peter", "id": 15},
	{"op": ">=", "key": "Data.Deep", "val": 1, "id": 17}
	]`)
	logic := "     13     and  (15or13    )"
	rules, err := NewRulesWithJsonAndLogic(jsonStr, logic)
	if err != nil {
		t.Error(err)
	}
	rules, err = NewRulesWithArrayAndLogic(rules.Rules, logic)
	if err != nil {
		t.Error(err)
	}
	t.Log(rules.Rules[0])
	t.Log(rules.Logic)
}

func TestNewRulesWithJsonAndLogicAndInfo(t *testing.T) {
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
	rules, err := NewRulesWithJsonAndLogicAndInfo(jsonStr, logic, extractInfo)
	if err != nil {
		t.Error(err)
	}
	t.Log(rules.Rules[0])
	t.Log(rules.Name)
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
	rules, err := NewRulesWithJsonAndLogic(jsonStr, logic)
	if err != nil {
		t.Error(err)
	}
	rules, err = NewRulesWithArrayAndLogicAndInfo(rules.Rules, logic, extractInfo)
	if err != nil {
		t.Error(err)
	}
	t.Log(rules.Rules[0])
	t.Log(rules.Msg)
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
	rules, err := NewRulesWithJsonAndLogic(jsonStr, logic)
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
	rules, err := NewRulesWithJsonAndLogicAndInfo(jsonStr, logic, extractInfo)
	if err != nil {
		t.Error(err)
	}

	rulesSet := NewRulesSet([]*Rules{rules}, extractInfo)

	obj := map[string]interface{}{"Name": "peter", "Status": "abcd", "Key": 0}
	fitRules, _ := rules.FitWithMap(obj)
	t.Log(fitRules)

	result := rulesSet.FitSetWithMap(obj)
	t.Log(result)
	t.Log(result == nil)
	t.Log(len(result) == 0)
}

func TestRules_FitWithMap(t *testing.T) {
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
	rules, err := NewRulesWithJsonAndLogicAndInfo(jsonStr, logic, extractInfo)
	if err != nil {
		t.Error(err)
	}
	objMap := map[string]interface{}{"Status": "abcd"}
	fit, _ := rules.FitWithMap(objMap)
	t.Log(fit)
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
	rules, err := NewRulesWithJsonAndLogicAndInfo(jsonStr, logic, extractInfo)
	if err != nil {
		t.Error(err)
	}
	objMap := map[string]interface{}{"Status": "abcd"}
	fit, _, val := rules.FitWithMapAskVal(objMap)
	t.Log(fit)
	t.Log(val)
}
