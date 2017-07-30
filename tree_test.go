package go_rule_engine

import (
	"testing"
	"github.com/satori/go.uuid"
)

func TestLogicToTree(t *testing.T) {
	logic := "1 or 2"
	head, err := logicToTree(logic)
	if err != nil {
		t.Error(err)
	}
	t.Log(head)
}

func TestSplitExprToChildren(t *testing.T) {
	expr := "2 and 3 and not 1"
	children := splitExprToChildren(expr)
	printChildren(children, t)
}

func TestSplitExprToChildren3(t *testing.T) {
	expr := " "
	children := splitExprToChildren(expr)
	t.Log(children == nil)
	printChildren(children, t)
}

func printChildren(children []*Node, t *testing.T) {
	for _, v := range children {
		t.Log(v.Expr)
	}
}

func TestSplitExprToChildren2(t *testing.T) {
	uuid1 := uuid.NewV1()
	t.Log(uuid1)
	uuid4 := uuid.NewV4()
	t.Log(uuid4)
}

func TestReplaceBiggestBracketContent3(t *testing.T) {
	expr := "1 or 2 and 3 or ( 2 and 4 )"
	result, _ := replaceBiggestBracketContentAtOnce(expr, make(map[string]string, 0))
	t.Log(result)
}

func TestReplaceBiggestBracketContent(t *testing.T) {
	expr := "( 1 or 2 ) and 3 or ( 2 and 4 )"
	result, _ := replaceBiggestBracketContentAtOnce(expr, make(map[string]string, 0))
	t.Log(result)
}

func TestReplaceBiggestBracketContent2(t *testing.T) {
	expr := "( 1 or 2 ) and 3 or ( 2 and 4 )"
	result, mapReplace := replaceBiggestBracketContent(expr)
	t.Log(result)
	t.Log(mapReplace)
}

func TestSplitExprToChildren4(t *testing.T) {
	expr := "( 1 or 2 ) and 3 or ( 2 and 4 )"
	children := splitExprToChildren(expr)
	printChildren(children, t)
}