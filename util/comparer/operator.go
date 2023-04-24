package comparer

type Operator string

const (
	Operator_Exist       Operator = "exist"
	Operator_Eq          Operator = "eq"
	Operator_Ueq         Operator = "ueq"
	Operator_Gt          Operator = "gt"
	Operator_Egt         Operator = "egt"
	Operator_Lt          Operator = "lt"
	Operator_Elt         Operator = "elt"
	Operator_Like        Operator = "like"
	Operator_NotLike     Operator = "not-like"
	Operator_Regex       Operator = "regex"
	Operator_Contains    Operator = "contains"
	Operator_NotContains Operator = "not-contains"
	Operator_In          Operator = "in"
	Operator_NotIn       Operator = "not-in"
	Operator_Empty       Operator = "empty"
	Operator_NotEmpty    Operator = "not-empty"
)

const (
	Option_None            = iota
	Option_CaseSensitive   = 1
	Option_CaseInsensitive = 2
)

var (
	operatorAlias = map[Operator]Operator{
		"=":         Operator_Eq,
		"==":        Operator_Eq,
		"!=":        Operator_Ueq,
		">":         Operator_Gt,
		"≥":         Operator_Egt,
		">=":        Operator_Egt,
		"<":         Operator_Lt,
		"≤":         Operator_Elt,
		"<=":        Operator_Elt,
		"!like":     Operator_NotLike,
		"!contains": Operator_NotContains,
		"!in":       Operator_NotIn,
		"nil":       Operator_Empty,
		"!empty":    Operator_NotEmpty,
		"!nil":      Operator_NotEmpty,
	}
)
