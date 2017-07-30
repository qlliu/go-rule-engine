package go_rule_engine

import (
	"github.com/satori/go.uuid"
	"regexp"
	"strings"
)

/**
  将逻辑表达式转化为树，返回树的根节点
*/
func logicToTree(logic string) (*Node, error) {
	if logic == "" || logic == " " {
		return nil, nil
	}
	var head = &Node{
		Expr:   logic,
		Should: true,
		Blamed: true,
	}
	head.Leaf = head.isLeaf()
	propagateTree(head)
	return head, nil
}

func propagateTree(head *Node) {
	children := head.splitExprToChildren()
	if children != nil {
		head.Children = children
	} else {
		return
	}
	for index := range head.Children {
		propagateTree(head.Children[index])
	}
}

func (node *Node) splitExprToChildren() []*Node {
	// wrap biggest (***) block
	exprWrap, mapReplace := replaceBiggestBracketContent(node.Expr)
	// or layer
	ors := strings.Split(exprWrap, " or ")
	if len(ors) > 1 {
		node.ChildrenOp = string(OperatorOr)
		return node.shipChildren(ors, mapReplace)
	}
	// and layer
	ands := strings.Split(exprWrap, " and ")
	if len(ands) > 1 {
		node.ChildrenOp = string(OperatorAnd)
		return node.shipChildren(ands, mapReplace)
	}
	// not layer
	not := strings.Split(exprWrap, "not ")
	if len(not) > 1 {
		node.ChildrenOp = string(OperatorNot)
		return node.shipChildren([]string{not[1]}, mapReplace)
	}
	return nil
}

func (node *Node) shipChildren(splits []string, mapReplace map[string]string) []*Node {
	var children = make([]*Node, 0)
	var isFirstChild = true
	for _, o := range splits {
		for k, v := range mapReplace {
			if o == k {
				o = mapReplace[k]
			} else if strings.Contains(o, k) {
				o = strings.Replace(o, k, "( "+v+" )", -1)
			}
		}
		var should bool
		var blamed bool
		switch node.ChildrenOp {
		case string(OperatorAnd):
			should = true
			// and和not的时候所有子树都有责任
			blamed = true
		case string(OperatorOr):
			should = true
			// or的时候只有第一个子树有责任
			if isFirstChild {
				blamed = true
			}
		case string(OperatorNot):
			should = false
			// and和not的时候所有子树都有责任
			blamed = true
		}
		// 父节点无责任，子树也无责任
		blamed = node.Blamed && blamed

		child := &Node{
			Expr:   o,
			Should: should,
			Blamed: blamed,
		}
		// judge if leaf
		child.Leaf = child.isLeaf()

		children = append(children, child)
		isFirstChild = false
	}
	return children
}

func (node *Node) isLeaf() bool {
	if flag, _ := regexp.MatchString("^\\d+$", node.Expr); flag {
		return true
	}
	return false
}

func replaceBiggestBracketContent(expr string) (string, map[string]string) {
	var result = expr
	var mapReplace = make(map[string]string, 0)
	for {
		before := result
		result, mapReplace = replaceBiggestBracketContentAtOnce(result, mapReplace)
		if before == result {
			// replace finished
			break
		}
	}
	return result, mapReplace
}

func replaceBiggestBracketContentAtOnce(expr string, mapReplace map[string]string) (string, map[string]string) {
	var result = expr
	var flag bool
	bracketStack := make([]rune, 0)
	toReplace := make([]rune, 0)
	runeExpr := []rune(expr)

	for _, v := range runeExpr {
		if flag {
			// add to buffer
			toReplace = append(toReplace, v)
		}
		if v == '(' {
			flag = true
			bracketStack = append(bracketStack, v)
		} else if v == ')' {
			// delete last ')'
			bracketStack = append(bracketStack[:len(bracketStack)-1])
			if len(bracketStack) == 0 {
				// it's one biggest (***)block, break to replace
				break
			}
		}
	}

	if flag {
		// delete last )
		toReplace = toReplace[:len(toReplace)-1]
		key := uuid.NewV1().String()
		result = strings.Replace(result, "("+string(toReplace)+")", key, 1)
		mapReplace[key] = strings.Trim(string(toReplace), " ")
	}
	return result, mapReplace
}
