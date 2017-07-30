package go_rule_engine

import (
	"strings"
	"github.com/satori/go.uuid"
)

func logicToTree(logic string) (head *Node, err error) {

	return nil, nil
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
				o = strings.Replace(o, k, "( " + v + " )", -1)
			}
		}

		child := &Node{
			Expr: o,
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
		} else if flag {
			// add to buffer
			toReplace = append(toReplace, v)
		}
	}

	if flag {
		key := uuid.NewV1().String()
		result = strings.Replace(result, "("+string(toReplace)+")", key, 1)
		mapReplace[key] = strings.Trim(string(toReplace), " ")
	}
	return result, mapReplace
}
