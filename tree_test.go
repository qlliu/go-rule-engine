package ruler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogicToTree(t *testing.T) {
	logic := "1 or 2"
	head := logicToTree(logic)
	t.Log(head)
	assert.NotNil(t, head)
	assert.Equal(t, "1 or 2", head.Expr)
}

func TestReplaceBiggestBracketContent3(t *testing.T) {
	expr := "1 or 2 and 3 or ( 2 and 4 )"
	result, _ := replaceBiggestBracketContentAtOnce(expr, make(map[string]string))
	t.Log(result)
	assert.Contains(t, result, "1 or 2 and 3")
}

func TestReplaceBiggestBracketContent(t *testing.T) {
	expr := "( 1 or 2 ) and 3 or ( 2 and 4 )"
	result, _ := replaceBiggestBracketContentAtOnce(expr, make(map[string]string))
	t.Log(result)
	assert.Contains(t, result, "3 or ( 2 and 4 )")
}

func TestReplaceBiggestBracketContent2(t *testing.T) {
	expr := "( 1 or 2 ) and 3 or ( 2 and 4 )"
	result, mapReplace := replaceBiggestBracketContent(expr)
	t.Log(result)
	t.Log(mapReplace)
	assert.Contains(t, result, "and 3 or")
}

func TestReplaceBiggestBracketContent4(t *testing.T) {
	expr := "( 1 or 2 ) and ( 3 or ( 2 and 4 ) )"
	result, mapReplace := replaceBiggestBracketContent(expr)
	t.Log(result)
	t.Log(mapReplace)
	assert.Contains(t, result, "and")
}

func TestReplaceBiggestBracketContent5(t *testing.T) {
	expr := "1 or 2 and ( 3 or ( 2 and 4 ) )"
	result, _ := replaceBiggestBracketContentAtOnce(expr, make(map[string]string))
	t.Log(result)
	assert.Contains(t, result, "1 or 2 and")
}

func TestLogicToTree2(t *testing.T) {
	logic := "1 and 2 and ( 3 or not ( 2 and 4 ) )"
	head := logicToTree(logic)
	traverseTreeInLayer(head)
	assert.NotNil(t, head)
}

func traverseTreeInLayer(head *Node) {
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

func traverseTreeInPostOrder(head *Node) {
	children := head.Children
	for _, child := range children {
		traverseTreeInPostOrder(child)
	}

	if head.Leaf {
		fmt.Printf("%+v\n", head)
		return
	}
	fmt.Printf("%+v\n", head)
}

func TestLogicToTree3(t *testing.T) {
	logic := "1 and 2 and ( 3 or not ( 2 and 4 ) )"
	head := logicToTree(logic)
	traverseTreeInPostOrder(head)
	assert.NotNil(t, head)
}
