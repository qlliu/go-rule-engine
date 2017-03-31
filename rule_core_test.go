package go_rule_engine

import "testing"

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
	{"op": "@", "key": "status", "val": "abcd", "id": 13},
	{"op": "=", "key": "name", "val": "peter", "id": 15},
	{"op": ">=", "key": "data.deep", "val": 1, "id": 17}
	]`)
	rules, err := NewRulesWithJson(jsonStr)
	if err != nil {
		t.Error(err)
	}
	rules.Logic = "( 15 or 13 ) and 17 and not 13"

	obj := map[string]interface{}{"data": map[string]interface{}{"deep": 1}, "name": "peter", "status": "abc"}
	result := rules.Fit(obj)
	t.Log(result)
}
