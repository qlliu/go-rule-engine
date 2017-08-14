package go_rule_engine

type Rule struct {
	Op  string      `json:"op"`  // 预算符
	Key string      `json:"key"` // 目标变量键名
	Val interface{} `json:"val"` // 目标变量子规则存值
	Id  int         `json:"id"`  // 子规则ID
	Msg string      `json:"msg"` // 该规则抛出的负提示
}

type Rules struct {
	Rules []*Rule // 子规则集合
	Logic string  // 逻辑表达式，使用子规则ID运算表达
	Name  string  // 规则名称
	Msg   string  // 规则抛出的负提示
}

type RulesSet struct {
	RulesSet []*Rules
	Name     string
	Msg      string
}

var VALID_OPERATORS = []string{"and", "or", "not"}

type Node struct {
	Expr       string  // 分割的logic表达式
	Val        bool    // 节点值
	Computed   bool    // 节点值被计算过
	Children   []*Node // 孩子树
	ChildrenOp string  // 孩子树之间的运算符: and, or, not
	Leaf       bool    // 是否叶子节点
	Should     bool    // 为了Fit为true，此节点必须的值
	Blamed     bool    // 此节点为了Fit为true，有责任必须为某值
}

type Operator string

const (
	OperatorAnd Operator = "and"
	OperatorOr  Operator = "or"
	OperatorNot Operator = "not"
)
