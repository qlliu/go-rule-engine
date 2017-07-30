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
	var head = shipChildren([]string{logic}, true, nil)[0]
	propagateTree(head)
	return head, nil
}

func propagateTree(head *Node) {
	children := splitExprToChildren(head.Expr)
	if children != nil {
		head.Children = children
	} else {
		return
	}
	for index := range head.Children {
		propagateTree(head.Children[index])
	}
}

func splitExprToChildren(expr string) []*Node {
	// wrap biggest (***) block
	exprWrap, mapReplace := replaceBiggestBracketContent(expr)
	// or layer
	ors := strings.Split(exprWrap, " or ")
	if len(ors) > 1 {
		return shipChildren(ors, true, mapReplace)
	}
	// and layer
	ands := strings.Split(exprWrap, " and ")
	if len(ands) > 1 {
		return shipChildren(ands, true, mapReplace)
	}
	// not layer
	not := strings.Split(exprWrap, "not ")
	if len(not) > 1 {
		return shipChildren([]string{not[1]}, false, mapReplace)
	}
	return nil
}

func shipChildren(splits []string, should bool, mapReplace map[string]string) []*Node {
	var children = make([]*Node, 0)
	for _, o := range splits {
		for k, v := range mapReplace {
			if o == k {
				o = mapReplace[k]
			} else if strings.Contains(o, k) {
				o = strings.Replace(o, k, "( "+v+" )", -1)
			}
		}
		// judge if leaf
		var leaf bool
		if flag, _ := regexp.MatchString("^\\d+$", o); flag {
			leaf = true
		}

		child := &Node{
			Expr:   o,
			Leaf:   leaf,
			Should: should,
		}
		children = append(children, child)
	}
	return children
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
