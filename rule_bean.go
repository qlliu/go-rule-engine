package go_rule_engine

type Rule struct {
	Op    string      `json:"op"`
	Key   string      `json:"key"`
	Val   interface{} `json:"val"`
	Id    int         `json:"id"`
	Msg   string 	  `json:"msg"`
}

type Rules struct {
	Rules []*Rule
	Logic string
	Name  string
	Msg   string
}

type RulesSet struct {
	RulesSet []*Rules
	Name     string
	Msg 	 string
}
