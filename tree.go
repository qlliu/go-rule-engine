package go_rule_engine

import (
	"errors"
	"github.com/satori/go.uuid"
	"regexp"
	"strconv"
	"strings"
)

/**
  利用树来计算规则引擎
  输入：子规则ID和逻辑值map
  输出：规则匹配结果，导致匹配false的子规则ID
*/
func (rs *Rules) calculateExpressionByTree(values map[int]bool) (bool, []int, error) {
	var unfitIDs []int
	head := logicToTree(rs.Logic)
	leafs := head.traverseTreeInLayerAskForAllLeafs()
	for _, leaf := range leafs {
		if leaf.Blamed {
			ruleId, err := strconv.Atoi(leaf.Expr)
			if err != nil {
				return false, nil, err
			}
			actual, ok := values[ruleId]
			if !ok {
				return false, nil, errors.New("not found this rule value: " + leaf.Expr)
			}
			if leaf.Should != actual {
				// 此规则导致整体不匹配
				unfitIDs = append(unfitIDs, ruleId)
			}
		}
	}
	return len(unfitIDs) == 0, unfitIDs, nil
}

/**
  将逻辑表达式转化为树，返回树的根节点
*/
func logicToTree(logic string) *Node {
	if logic == "" || logic == " " {
		return nil
	}
	var head = &Node{
		Expr:   logic,
		Should: true,
		Blamed: true,
	}
	head.Leaf = head.isLeaf()
	propagateTree(head)
	return head
}

/**
  层序遍历获取树里面的所有叶子节点
*/
func (node *Node) traverseTreeInLayerAskForAllLeafs() []*Node {
	var leafs []*Node
	var buf []*Node
	var i int
	buf = append(buf, node)
	for {
		if i < len(buf) {
			if buf[i].Leaf {
				leafs = append(leafs, buf[i])
			}
			if buf[i].Children != nil {
				buf = append(buf, buf[i].Children...)
			}
			i++
		} else {
			break
		}
	}
	return leafs
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
