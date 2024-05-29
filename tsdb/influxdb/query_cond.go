package influxdb

import (
	"fmt"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/toolkits/slice"
	"strings"
)

type ICond interface {
	IsValid() bool
	And(...ICond) ICond
	Or(...ICond) ICond
	String() string
}

var _ ICond = condEmpty{}

type condEmpty struct{}

func NewCond() ICond {
	return condEmpty{}
}

func (condEmpty) String() string {
	return ""
}

func (condEmpty) And(arr ...ICond) ICond {
	return And(arr...)
}

func (condEmpty) Or(arr ...ICond) ICond {
	return Or(arr...)
}

func (condEmpty) IsValid() bool {
	return false
}

var _ ICond = condAnd{}

type condAnd []ICond

func And(arr ...ICond) ICond {
	result := make(condAnd, 0, len(arr))
	for _, cond := range arr {
		if cond == nil || !cond.IsValid() {
			continue
		}
		result = append(result, cond)
	}
	return result
}

func (and condAnd) String() string {
	arr := make([]string, 0, len(and))
	for _, cond := range and {
		var needQuote bool
		switch cond.(type) {
		case condOr, condRawExpr:
			needQuote = true
		}

		if needQuote {
			arr = append(arr, "("+cond.String()+")")
		} else {
			arr = append(arr, cond.String())
		}
	}

	return strings.Join(arr, " AND ")
}

func (and condAnd) And(arr ...ICond) ICond {
	return And(and, And(arr...))
}

func (and condAnd) Or(arr ...ICond) ICond {
	return Or(and, Or(arr...))
}

func (and condAnd) IsValid() bool {
	return len(and) > 0
}

type condOr []ICond

var _ ICond = condOr{}

// Or 或条件
//	！！！注意时间字段不能使用 OR 连接，否则会返回空结果！！！
//	https://docs.influxdata.com/influxdb/v1/troubleshooting/frequently-asked-questions/#why-is-my-query-with-a-where-or-time-clause-returning-empty-results
func Or(arr ...ICond) ICond {
	result := make(condOr, 0, len(arr))
	for _, cond := range arr {
		if cond == nil || !cond.IsValid() {
			continue
		}
		result = append(result, cond)
	}
	return result
}

func (o condOr) String() string {
	arr := make([]string, 0, len(o))
	for _, cond := range o {
		var needQuote bool
		switch cond.(type) {
		case condAnd, condRawExpr:
			needQuote = true
		}

		if needQuote {
			arr = append(arr, "("+cond.String()+")")
		} else {
			arr = append(arr, cond.String())
		}
	}

	return strings.Join(arr, " OR ")
}

func (o condOr) And(arr ...ICond) ICond {
	return And(o, And(arr...))
}

func (o condOr) Or(arr ...ICond) ICond {
	return Or(o, Or(arr...))
}

func (o condOr) IsValid() bool {
	return len(o) > 0
}

var _ ICond = &condRawExpr{}

type condRawExpr struct {
	expr string
}

// RawExpr 表达式
func RawExpr(expr string) ICond {
	return condRawExpr{expr: expr}
}

func (c condRawExpr) String() string {
	return c.expr
}

func (c condRawExpr) And(arr ...ICond) ICond {
	return And(c, And(arr...))
}

func (c condRawExpr) Or(arr ...ICond) ICond {
	return Or(c, Or(arr...))
}

func (c condRawExpr) IsValid() bool {
	return len(c.expr) > 0
}

var _ ICond = &condExpr{}

type condExpr struct {
	col string // 仅支持tag字段（包含time）
	opr string
	val interface{}
}

// Expr 表达式
//	- col: 字段名，支持tag字段（包含time）
//	- opr: 操作符，支持 =, !=, >, <, >=, <=, <>, =~, !~
//	- val: 字段值，支持数值、字符串
func Expr(col, opr string, val interface{}) ICond {
	return condExpr{col, opr, val}
}

func (c condExpr) String() string {
	if !c.IsValid() {
		return ""
	}

	if c.col == "time" {
		// 时间字段值如果是数值类型，不需要带单引号
		if v, ok := c.val.(string); ok {
			c.val = SingleQuote(v)
		}
	} else {
		c.col = QuoteIfNeed(c.col)
		// TODO tag字段值是字符串类型，所以查询条件值总带单引号，暂不考虑field字段：
		//	1、有可能数值型的field字段值用作查询条件，此时值无需加单引号
		//	2、field字段值使用EscapeCondValue时需要指定isFieldVal=true
		c.val = SingleQuote(EscapeCondValue(convertor.ToStringNoError(c.val)))
	}

	return fmt.Sprintf(`%s %s %v`, c.col, c.opr, c.val)
}

func (c condExpr) And(arr ...ICond) ICond {
	return And(c, And(arr...))
}

func (c condExpr) Or(arr ...ICond) ICond {
	return Or(c, Or(arr...))
}

func (c condExpr) IsValid() bool {
	return len(c.col) > 0 && c.val != nil &&
		slice.ContainsString([]string{"=", "!=", ">", "<", ">=", "<=", "<>", "=~", "!~"}, c.opr)
}

// Between 左右均为闭区间
func Between(col string, less, more interface{}) ICond {
	return And(condExpr{col, ">=", less}, condExpr{col, "<=", more})
}

func BetweenOpen(col string, less, more interface{}) ICond {
	return And(condExpr{col, ">", less}, condExpr{col, "<", more})
}

func BetweenOpenR(col string, less, more interface{}) ICond {
	return And(condExpr{col, ">=", less}, condExpr{col, "<", more})
}

func BetweenOpenL(col string, less, more interface{}) ICond {
	return And(condExpr{col, ">", less}, condExpr{col, "<=", more})
}

func In(col string, values ...interface{}) ICond {
	if len(values) == 0 {
		return condEmpty{}
	}

	var arr []ICond
	for _, val := range values {
		arr = append(arr, condExpr{col, "=", val})
	}
	return Or(arr...)
}

func NotIn(col string, values ...interface{}) ICond {
	if len(values) == 0 {
		return condEmpty{}
	}

	var arr []ICond
	for _, val := range values {
		arr = append(arr, condExpr{col, "!=", val})
	}
	return And(arr...)
}

// Match 模糊匹配
//	https://docs.influxdata.com/influxdb/v1/query_language/explore-data/#regular-expressions
//	!!! 注意：正则表达式匹配性能较差，尽量避免使用 !!!
func Match(col, pattern string) ICond {
	return condRawExpr{fmt.Sprintf(`%s =~ /%s/`, QuoteIfNeed(col), pattern)}
}

// NotMatch 不匹配
//	!!! 注意：正则表达式匹配性能较差，尽量避免使用 !!!
func NotMatch(col, pattern string) ICond {
	return condRawExpr{fmt.Sprintf(`%s !~ /%s/`, QuoteIfNeed(col), pattern)}
}
