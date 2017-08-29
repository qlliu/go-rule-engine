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
	assert.False(t, fit)
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
	assert.True(t, fit)
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

func TestRules_Fit6(t *testing.T) {
	// ExpectTimeOkRuleJSON 预期时间ok的规则
	var ExpectTimeOkRuleJSON = []byte(`[
	{"op": "<", "key": "SecondsAfterOnShelf", "val": 21600, "id": 1, "msg": "新上架<6h"},
	{"op": "=", "key": "CustomerType", "val": "new", "id": 2, "msg": "新客户"},
	{"op": ">", "key": "SecondsBetweenWatchAndOnShelf", "val": 21600, "id": 3, "msg": "需要带看在上架6h以后"},
	{"op": "=", "key": "FinanceAuditPass", "val": 1, "id": 4, "msg": "需要预审通过"},
	{"op": "!=", "key": "IsDealer", "val": 1, "id": 5, "msg": "不能是车商"}
	]`)
	// ExpectTimeOkRuleLogic 判断预期时间ok的逻辑
	var ExpectTimeOkRuleLogic = "(1 and ((2 and 3) or (2 and 4 and 5) or not 2)) or not 1"
	rule, err := NewRulesWithJSONAndLogic(ExpectTimeOkRuleJSON, ExpectTimeOkRuleLogic)
	if err != nil {
		t.Error(err)
	}
	t.Log(rule)

	// wrap data
	type A struct {
		SecondsAfterOnShelf           int
		CustomerType                  string
		SecondsBetweenWatchAndOnShelf int
		FinanceAuditPass              int
		IsDealer                      int
	}
	a := &A{
		SecondsAfterOnShelf:           2160,
		CustomerType:                  "new",
		SecondsBetweenWatchAndOnShelf: 2160,
		FinanceAuditPass:              0,
		IsDealer:                      1,
	}

	fit, msg := rule.Fit(a)
	t.Log(fit)
	t.Log(msg)
	assert.False(t, fit)
	assert.Equal(t, "需要带看在上架6h以后", msg[3])
}

func TestRules_Fit7(t *testing.T) {
	var jsonIn = []byte(`[
	{"op": "@", "key": "A", "val": "11, 2, 3, 1", "id": 1, "msg": "error 1"},
	{"op": "!@", "key": "B", "val": "11, 2, 3, 1", "id": 2, "msg": "error 2"},
	{"op": "@", "key": "C", "val": "11, 2, 3, 1", "id": 3, "msg": "error 3"},
	{"op": "!@", "key": "D", "val": "11, 2, 3, 1", "id": 4, "msg": "error 4"},
	{"op": "@", "key": "E", "val": "11, 2, 3, 1,  88.1", "id": 5, "msg": "error 5"},
	{"op": "@", "key": "F", "val": "11, 2, 3, 1,  88.1", "id": 6, "msg": "error 6"},
	{"op": "@", "key": "G", "val": "11, 2, 3, 1,  88.1, 0.001", "id": 7, "msg": "error 7"},
	{"op": "@", "key": "H", "val": "11, 2, 3, 1,  88.1, 0.001,    ab c", "id": 8, "msg": "error 8"}
	]`)
	r, err := NewRulesWithJSONAndLogic(jsonIn, "")
	if err != nil {
		t.Error(err)
	}

	// to fit
	type obj struct {
		A string
		B string
		C int
		D int32
		E float32
		F float64
		G float64
		H string
	}
	o := &obj{
		A: "3",
		B: "4",
		C: 1,
		D: 8,
		E: 88.1,
		F: 88.12,
		G: 1e-3,
		H: "ab c",
	}
	fit, msg := r.Fit(o)
	t.Logf("result: %v", fit)
	t.Logf("msg: %+v", msg)
	assert.False(t, fit)
	assert.Equal(t, "error 6", msg[6])
}

func TestGetRuleIDsByLogicExpression(t *testing.T) {
	logic := "1 and (2 or (4or not5))"
	ids, err := GetRuleIDsByLogicExpression(logic)
	if err != nil {
		t.Error(err)
	}
	t.Log(ids)
	assert.Equal(t, []int{1, 2, 4, 5}, ids)
}
