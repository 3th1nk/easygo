// ------------------------------------------------------------------------------
// 类型转换函数集
// ------------------------------------------------------------------------------
package convertor

import (
	"fmt"
	"github.com/modern-go/reflect2"
)

func ToError(a interface{}) (err error) {
	if reflect2.IsNil(a) {
		return
	}
	if err, _ = a.(error); err == nil {
		err = fmt.Errorf("%v", a)
	}
	return
}
