package go_rule_engine

import "testing"

func TestNewRulesWithJson(t *testing.T) {
	jsonStr := []byte(`[{"op": "=", "key": "status", "val": 1}]`)

	rule, err := NewRulesWithJson(jsonStr)
	if err != nil {
		t.Error(err)
	}
	t.Log(rule.Name)
}

func TestPluck(t *testing.T) {
	result := pluck("op.deep", map[string]interface{}{"op": map[string]interface{}{"deep": 1}, "key": "status", "val": 1})
	t.Log(result)
}