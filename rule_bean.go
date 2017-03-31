package go_rule_engine

type Rule struct {
	Op    string      `json:"op"`
	Key   string      `json:"key"`
	Val   interface{} `json:"val"`
	Id    int         `json:"id"`
}

type Rules struct {
	Rules []*Rule
	Logic string
	Name  string
}

type RulesSet struct {
	RulesSet []*Rules
	Name     string
}
