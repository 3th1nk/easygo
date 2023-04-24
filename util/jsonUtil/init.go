package jsonUtil

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/json-iterator/go/extra"
)

var (
	DefaultTimeFormat = "2006-01-02 15:04:05"
)

func init() {
	// 默认使用 小写下划线 格式
	extra.SetNamingStrategy(strUtil.PascalToSnake)

	// 注册弱类型解析
	registerFuzzy()
}
