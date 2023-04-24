/**
这是一个动态表达式处理类，可以将字符串中的动态表达式替换为真实的值。
用法参考单元测试代码。
*/

package dynamicExpr

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/arrUtil"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"regexp"
	"strings"
	"unicode"
)

// 默认的格式化器。
func Default() *Formatter { return _default }

// 严格匹配大括号的格式化器，要求表达式必须用 {} 包裹，即 {$表达式}
func BraceStrict() *Formatter { return _braceStrict }

var (
	_default     = mustNew(`\$[a-zA-Z_][a-zA-Z0-9-_.*]*`)
	_braceStrict = mustNew(`{\$[a-zA-Z_][a-zA-Z0-9-_.*]*}`)
)

func New(pattern string) (*Formatter, error) {
	obj := &Formatter{}
	if err := obj.SetPattern(pattern); err != nil {
		return nil, err
	}
	return obj, nil
}

func mustNew(pattern string) *Formatter {
	obj, _ := New(pattern)
	return obj
}

type ValueProvider interface {
	Value(path string) (val interface{}, err error)
}

type Options struct {
	VariablePattern string
}

type Formatter struct {
	pattern         string
	exprRegex       *regexp.Regexp
	simpleExprRegex *regexp.Regexp
}

type Pattern struct {
	fmt       *Formatter
	exclude   []string
	checkJson int // 0=未设置; 1=check; 2=no-check
}

func (this *Formatter) GetPattern() string {
	return this.pattern
}

func (this *Formatter) SetPattern(pattern string) error {
	reg1, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	reg2, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return err
	}
	this.pattern, this.exprRegex, this.simpleExprRegex = pattern, reg1, reg2
	return nil
}

// 判断指定的字符串是否是一个简单表达式（包含且只包含一个变量表达式）
// 例：
//   '$obj.name' 是一个简单表达式，因为整体个字符串都是变量表达式。
//   'abc{$obj.name}' 不是简单表达式，因为除了变量表达式之外，还包含了一个 abc 前缀。
func (this *Formatter) IsSimpleExpr(s string) bool {
	return this.simpleExprRegex.MatchString(s)
}

// 获取所有变量表达式的起止位置。
func (this *Formatter) FindAllExprIndex(s string) [][]int {
	matches := this.exprRegex.FindAllStringIndex(s, -1)
	// 遍历每个匹配到的表达式，检查如果表达式前后被大括号包括着，则将大括号一起视为表达式一部分
	exprEnd := len(s) - 1
	for _, arr := range matches {
		i, j := arr[0], arr[1]-1
		for i != 0 {
			if unicode.IsSpace(rune(s[i-1])) {
				i--
			} else {
				break
			}
		}
		for j != exprEnd {
			if unicode.IsSpace(rune(s[j+1])) {
				j++
			} else {
				break
			}
		}
		if i > 0 && j < exprEnd {
			if i, j := i-1, j+1; s[i] == '{' && s[j] == '}' {
				arr[0], arr[1] = i, j+1
			}
		}
	}
	return matches
}

// 替换 expr 中的表达式
func (this *Formatter) Format(expr string, provider ValueProvider) (string, error) {
	return (&Pattern{fmt: this}).Format(expr, provider)
}

// 替换 expr 中的表达式
func (this *Formatter) Exclude(path ...string) *Pattern {
	return (&Pattern{fmt: this, exclude: path})
}

// 替换 expr 中的表达式
func (this *Formatter) CheckJson(val bool) *Pattern {
	return (&Pattern{fmt: this, checkJson: util.IfInt(val, 1, 2)})
}

// 替换 expr 中的表达式
func (this *Pattern) Format(expr string, provider ValueProvider) (string, error) {
	if expr = strings.TrimSpace(expr); expr == "" {
		return "", nil
	}

	if this.fmt.IsSimpleExpr(expr) {
		if len(this.exclude) != 0 && -1 != arrUtil.IndexOfString(this.exclude, expr, false) {
			return expr, nil
		} else {
			val, err := provider.Value(expr[1:])
			if err != nil {
				return "", fmt.Errorf("invalid expr '%v'", expr)
			} else {
				return convertor.ToString(val)
			}
		}
	}

	checkJson := true
	if this.checkJson == 2 {
		checkJson = false
	} else if this.checkJson == 0 {
		trimExpr := strings.TrimSpace(expr)
		if !strings.HasPrefix(trimExpr, "[") && !strings.HasPrefix(trimExpr, "{") {
			checkJson = false
		} else if strings.HasPrefix(trimExpr, "{$") {
			checkJson = false
		}
	}

	matches := this.fmt.FindAllExprIndex(expr)
	exprLen, matchCount := len(expr), len(matches)

	// 判断 matchStart, matchEnd 的前后是否是双引号
	isInQuotation := func(matchStart, matchEnd int) bool {
		for i := matchStart - 1; i >= 0; i-- {
			if c := expr[i]; unicode.IsSpace(rune(c)) {
				continue
			} else if c != '"' {
				return false
			}

			for j := matchEnd; j < exprLen; j++ {
				if c := expr[j]; unicode.IsSpace(rune(c)) {
					continue
				} else {
					return c == '"'
				}
			}
			break
		}
		return false
	}

	// 判断 matchStart, matchEnd 的前后不是 Json 格式
	maybeJsonMember := func(matchStart, matchEnd int) bool {
		for i := matchStart; i >= 0; i-- {
			if c := expr[i]; unicode.IsSpace(rune(c)) {
				continue
			} else if c != ':' && c != '[' && c != ',' {
				return false
			}
		}
		for i := matchEnd; i < exprLen; i++ {
			if c := expr[i]; unicode.IsSpace(rune(c)) {
				continue
			} else if c != ',' && c != ']' {
				return false
			}
		}
		return true
	}

	segCount := matchCount*2 + 1
	segs, copyEnd := make([]string, segCount), exprLen
	for i, j := matchCount-1, segCount-1; i >= 0; i, j = i-1, j-2 {
		matchStart, matchEnd := matches[i][0], matches[i][1]
		segs[j] = expr[matchEnd:copyEnd]
		path := expr[matchStart:matchEnd]
		if path[0] == '{' && path[len(path)-1] == '}' {
			path = path[1 : len(path)-1]
		}
		path = strings.TrimLeft(strings.TrimSpace(path), "$")
		if len(this.exclude) != 0 && -1 != arrUtil.IndexOfString(this.exclude, path, false) {
			segs[j-1] = path
		} else {
			val, err := provider.Value(path)
			if err != nil {
				segs[j-1] = expr[matchStart:matchEnd]
			} else if checkJson {
				if maybeJsonMember(matchStart, matchEnd) {
					// 如果表达式前后的格式表明这有可能是一个 JsonMember，则直接 Json 序列化
					segs[j-1] = jsonUtil.MustMarshalToString(val)
				} else if isInQuotation(matchStart, matchEnd) {
					// 如果表达式前后被双引号包括，则将 val 转化为字符串。 注意：Json Marshal 然后去掉首尾双引号是为了处理转义字符以便拼起来之后还是个合法的字符串
					str := jsonUtil.MustMarshalToString(convertor.ToStringNoError(val))
					segs[j-1] = str[1 : len(str)-1]
				} else {
					segs[j-1] = jsonUtil.MustMarshalToString(val)
				}
			} else {
				segs[j-1] = convertor.ToStringNoError(val)
			}
		}
		copyEnd = matchStart
	}
	segs[0] = expr[:copyEnd]
	return strings.Join(segs, ""), nil
}
