package go_rule_engine

import (
	"strings"
	"container/list"
	"regexp"
	"errors"
	"strconv"
)

func (*Rules) CalculateExpression(expr string, values map[int]bool) (bool, error) {
	listExpr := strings.Split(expr, " ")
	stackNum := list.New()
	stackOp	:= list.New()

	// regex
	patternNum := "^\\d*$"
	regexNum, err := regexp.Compile(patternNum)
	if err != nil {
		return false, err
	}

	for _, c := range listExpr {
		// judge if number
		if regexNum.MatchString(c) {
			stackNum.PushBack(c)
		} else {
			var lastOp string
			var ok bool
			if stackOp.Back() != nil {
				lastOpRaw := stackOp.Back().Value
				lastOp, ok = lastOpRaw.(string)
				if !ok {
					return false, errors.New("error type of operator")
				}
			}
			if isOpBiggerInLogic(c, lastOp) || c == "(" {
				stackOp.PushBack(c)
			} else {
				iterMax := stackOp.Len()
				for i := 0; i < iterMax; i++ {
					lastOpRaw := stackOp.Back()
					if lastOpRaw == nil {
						break
					}
					lastOp, ok := lastOpRaw.Value.(string)
					if !ok {
						return false, errors.New("error type of operator")
					}
					if (isOpBiggerInLogic(c, lastOp)) {
						break;
					} else {
						stackNum.PushBack(lastOpRaw.Value)
						stackOp.Remove(lastOpRaw)
					}
				}
				if (c == ")") {
					// delete "("
					stackOp.Remove(stackOp.Back())
				} else {
					stackOp.PushBack(c)
				}
			}
		}
	}

	// dump op to num stack
	iterMax := stackOp.Len()
	for i := 0; i < iterMax; i++ {
		if stackOp.Back() == nil {
			break
		}
		stackNum.PushBack(stackOp.Back().Value)
		stackOp.Remove(stackOp.Back())
	}

	// count
	iterMax = stackNum.Len()
	for i:= 0; i < iterMax; i++ {
		itemRaw := stackNum.Front().Value
		item, ok := itemRaw.(string)
		if !ok {
			return false, errors.New("error type in stack number")
		}
		if (regexNum.MatchString(item)) {
			index, err := strconv.Atoi(item)
			if err != nil {
				return false, err
			}
			if val, ok := values[index]; ok {
				stackOp.PushBack(val)
			} else {
				return false, errors.New("empty operand value in map: "+item)
			}

		} else {
			// choose operands and operate
			if numOfOperandInLogic(item) == 2 {
				operandBRaw := stackOp.Back().Value
				operandB, ok := operandBRaw.(bool)
				if !ok {
					return false, errors.New("error type of operator")
				}
				stackOp.Remove(stackOp.Back())
				operandARaw := stackOp.Back().Value
				operandA, ok := operandARaw.(bool)
				if !ok {
					return false, errors.New("error type of operator")
				}
				stackOp.Remove(stackOp.Back())
				computeOutput, err := computeOneInLogic(item, []bool{operandA, operandB})
				if err != nil {
					return false, errors.New("error in one compute")
				}
				stackOp.PushBack(computeOutput)
			}
			if numOfOperandInLogic(item) == 1 {
				operandBRaw := stackOp.Back().Value
				operandB, ok := operandBRaw.(bool)
				if !ok {
					return false, errors.New("error type of operator")
				}
				stackOp.Remove(stackOp.Back())
				computeOutput, err := computeOneInLogic(item, []bool{operandB})
				if err != nil {
					return false, errors.New("error in one compute")
				}
				stackOp.PushBack(computeOutput)
			}
		}
		stackNum.Remove(stackNum.Front())
	}

	result, ok := stackOp.Back().Value.(bool)
	if !ok {
		return false, errors.New("error type in final result")
	}
	return result, nil
}

func isOpBiggerInLogic(obj, base string) bool {
	if obj == "" {
		return false
	}
	if base == "" {
		return true
	}
	mapPriority := map[string]int8{"or": 2, "and": 3, "not": 5, "(": 0, ")": 1}
	return mapPriority[obj] > mapPriority[base]
}

func numOfOperandInLogic(op string) int8 {
	mapOperand := map[string]int8{"or": 2, "and": 2, "not": 1}
	return mapOperand[op]
}

func computeOneInLogic(op string, v []bool) (bool, error) {
	switch op {
	case "or":
		return v[0] || v[1], nil
	case "and":
		return v[0] && v[1], nil
	case "not":
		return !v[0], nil
	default:
		return false, errors.New("unrecognized op")
	}
}
