package go_rule_engine

type Rule struct {
	Op    string      `json:"op"`
	Key   string      `json:"key"`
	Val   interface{} `json:"val"`
	Id    int         `json:"id"`
	Logic string      `json:"logic"`
}

type Rules struct {
	Rules []*Rule
	Name  string
}

type RulesSet struct {
	RulesSet []*Rules
	Name     string
}
