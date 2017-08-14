package ruler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRules_Fit3(t *testing.T) {
	jsonRules := []byte(`[
	{"op": "=", "key": "A", "val": 3, "id": 1, "msg": "A fail"},
	{"op": ">", "key": "B", "val": 1, "id": 2, "msg": "B fail"},
	{"op": "<", "key": "C", "val": 5, "id": 3, "msg": "C fail"}
	]`)
	logic := "1 and 2 and ( not (1 or 2) or not 3)"
	rs, err := NewRulesWithJSONAndLogic(jsonRules, logic)
	if err != nil {
		t.Error(err)
	}
	type Obj struct {
		A int
		B int
		C int
	}
	o := &Obj{
		A: 3,
		B: 3,
		C: 3,
	}
	fit, msg := rs.Fit(o)
	t.Log(fit)
	t.Log(msg)
}

func TestRules_Fit4(t *testing.T) {
	jsonRules := []byte(`[
	{"op": "=", "key": "A", "val": 3, "id": 1, "msg": "A fail"},
	{"op": ">", "key": "B", "val": 1, "id": 2, "msg": "B fail"},
	{"op": "<", "key": "C", "val": 5, "id": 3, "msg": "C fail"}
	]`)
	logic := "1 or 2"
	rs, err := NewRulesWithJSONAndLogic(jsonRules, logic)
	if err != nil {
		t.Error(err)
	}
	type Obj struct {
		A int
		B int
		C int
	}
	o := &Obj{
		A: 3,
		B: 3,
		C: 3,
	}
	fit, msg := rs.Fit(o)
	t.Log(fit)
	t.Log(msg)

	head := logicToTree(logic)
	err = head.traverseTreeInPostOrderForCalculate(map[int]bool{1: true, 2: true})
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("\n%+v\n", head)
}

func TestRules_Fit5(t *testing.T) {
	jsonRules := []byte(`[
	{"op": "=", "key": "A", "val": 3, "id": 1, "msg": "A fail"},
	{"op": ">", "key": "B", "val": 1, "id": 2, "msg": "B fail"},
	{"op": "<", "key": "C", "val": 5, "id": 3, "msg": "C fail"}
	]`)
	logic := "1 2"
	_, err := NewRulesWithJSONAndLogic(jsonRules, logic)
	assert.NotNil(t, err)
}
