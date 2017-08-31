### Intro

go-rule-engine是用Golang实现的小型规则引擎，可以使用json串构造子规则，以及逻辑表达式表征子规则间逻辑运算关系。也可以通过直接构造子规则对象传入，这样方便子规则的持久化和复用。

### Usage

```go
    // build rules
	jsonRules := []byte(`[
	{"op": "=", "key": "Grade", "val": 3, "id": 1, "msg": "Grade not match"},
	{"op": "=", "key": "Sex", "val": "male", "id": 2, "msg": "not male"},
	{"op": ">=", "key": "Score.Math", "val": 90, "id": 3, "msg": "Math not so well"},
	{"op": ">=", "key": "Score.Physic", "val": 90, "id": 4, "msg": "Physic not so well"}
	]`)
	logic := "1 and not 2 and (3 or 4)"
	ruleToFit, err := NewRulesWithJSONAndLogic(jsonRules, logic)
	if err != nil {
		t.Error(err)
	}

	// prepare obj
	type Exams struct {
		Math int
		Physic int
	}
	type Student struct {
		Name string
		Grade int
		Sex string
		Score *Exams
	}
	Chris := &Student{
		Name: "Chris",
		Grade: 3,
		Sex: "female",
		Score: &Exams{Math: 88, Physic: 91},
	}

	// fit
	fit, msg := ruleToFit.Fit(Chris)
	t.Log(fit)
	t.Log(msg)
```

```go
	// result
	true    // fit
	map[]   // msg
```

```go
	Helen := &Student{
		Name: "Helen",
		Grade: 4,
		Sex: "female",
		Score: &Exams{Math: 96, Physic: 93},
	}
```

```go
	// result
	false   				// fit
	map[1:Grade not match]   // msg
```



### 概念

1. 子规则Rule

```go
{"op": "=", "key": "Grade", "val": 3, "id": 1, "msg": "Grade not match"}

// Rule 最小单元，子规则
type Rule struct {
	Op  string      `json:"op"`  // 预算符
	Key string      `json:"key"` // 目标变量键名
	Val interface{} `json:"val"` // 目标变量子规则存值
	ID  int         `json:"id"`  // 子规则ID
	Msg string      `json:"msg"` // 该规则抛出的负提示
}
```

2. 逻辑表达式logic

```go
"1 and 2 and (3 or 4)"   // 数字是rule的ID，当需要所有rule都为true，可以缺省写法：logic=""
```

3. 规则Rules

```go
// 即创建的ruleToFit对象
ruleToFit, err := NewRulesWithJSONAndLogic(jsonRules, logic)

// Rules 规则，拥有逻辑表达式
type Rules struct {
	Rules []*Rule // 子规则集合
	Logic string  // 逻辑表达式，使用子规则ID运算表达
	Name  string  // 规则名称
	Msg   string  // 规则抛出的负提示
}
```

4. 匹配结果Fit

```go
fit, msg := ruleToFit.Fit(Chris)
t.Log(fit)

false
```

5. 不匹配原因msg

```go
fit, msg := ruleToFit.Fit(Chris)
t.Log(msg)

map[1:Grade not match]   // 键1是导致fit为false的那个rule的ID，值是那个rule的msg字段，用于提示
```



### API

```go
// NewRulesWithJSONAndLogic 用json串构造Rules的标准方法，logic表达式如果没有则传空字符串
func NewRulesWithJSONAndLogic(jsonStr []byte, logic string) (*Rules, error)

// NewRulesWithArrayAndLogic 用rule数组构造Rules的标准方法，logic表达式如果没有则传空字符串
func NewRulesWithArrayAndLogic(rules []*Rule, logic string) (*Rules, error) 

// Fit Rules匹配传入结构体
func (rs *Rules) Fit(o interface{}) (bool, map[int]string) 

// FitWithMap Rules匹配map
func (rs *Rules) FitWithMap(o map[string]interface{}) (bool, map[int]string) 

// FitAskVal Rules匹配结构体，同时返回所有子规则key对应实际值
func (rs *Rules) FitAskVal(o interface{}) (bool, map[int]string, map[int]interface{}) 

// FitWithMapAskVal Rules匹配map，同时返回所有子规则key对应实际值
func (rs *Rules) FitWithMapAskVal(o map[string]interface{}) (bool, map[int]string, map[int]interface{}) 

// GetRuleIDsByLogicExpression 根据逻辑表达式得到规则id列表
func GetRuleIDsByLogicExpression(logic string) ([]int, error) 
```



### 支持的算符

```go
// 等于
case "=", "eq":

// 大于
case ">", "gt":

// 小于		
case "<", "lt":

// 大于或等于	
case ">=", "gte":

// 小于或等于		
case "<=", "lte":

// 不等于		
case "!=", "neq":

// 取值在...之中（val用逗号,分隔取值）		
case "@", "in":

// 取值不能在...之中		
case "!@", "nin":

// 正则表达式
case "^$", "regex":

// 为空不存在
case "0", "empty":

// 不为空
case "1", "nempty":

```

### 支持的逻辑

```go
// 与
and

// 或
or

// 非
not

// 括号，可嵌套
()
```

