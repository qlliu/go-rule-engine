package go_rule_engine

import (
	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestLogicToTree(t *testing.T) {
	logic := "1 or 2"
	head := logicToTree(logic)
	t.Log(head)
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

func TestReplaceBiggestBracketContent4(t *testing.T) {
	expr := "( 1 or 2 ) and ( 3 or ( 2 and 4 ) )"
	result, mapReplace := replaceBiggestBracketContent(expr)
	t.Log(result)
	t.Log(mapReplace)
}

func TestReplaceBiggestBracketContent5(t *testing.T) {
	expr := "1 or 2 and ( 3 or ( 2 and 4 ) )"
	result, _ := replaceBiggestBracketContentAtOnce(expr, make(map[string]string, 0))
	t.Log(result)
}

func TestLogicToTree2(t *testing.T) {
	logic := "1 and 2 and ( 3 or not ( 2 and 4 ) )"
	head := logicToTree(logic)
	traverseTreeInLayer(head, t)
}

func traverseTreeInLayer(head *Node, t *testing.T) {
	var buf []*Node
	var i int
	buf = append(buf, head)
	for {
		if i < len(buf) {
			fmt.Printf("%+v\n", buf[i])
			if buf[i].Children != nil {
				buf = append(buf, buf[i].Children...)
			}
			i++
		} else {
			break
		}
	}
}

func TestTraverseTreeInLayerAskForAllLeafs(t *testing.T) {
	logic := "1 and 2 and ( 3 or not ( 2 and 4 ) )"
	head := logicToTree(logic)
	leafs := head.traverseTreeInLayerAskForAllLeafs()
	for _, v := range leafs {
		fmt.Printf("%+v\n", v)
	}
}
