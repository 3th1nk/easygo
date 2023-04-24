package _gen

import (
	"fmt"
	"github.com/3th1nk/easygo/util/strUtil"
	"strings"
	"testing"
)

func TestStrSliceTo(t *testing.T) {
	for _, typ := range strings.Split("int,int64,int32", ",") {
		code := strUtil.Replace(`
func StrSliceTo{FuncSuffix}(in []string) ([]{TypeName}, error) {
	arr, err := make([]{TypeName}, len(in)), error(nil)
	for i, s := range in {
		if arr[i], err = StrTo{FuncSuffix}(s); err != nil {
			return nil, err
		}
	}
	return arr, nil
}

func StrSliceTo{FuncSuffix}NoError(in []string) []{TypeName} {
	arr, _ := StrSliceTo{FuncSuffix}(in)
	return arr
}
`, map[string]string{
			"{TypeName}":   typ,
			"{FuncSuffix}": strUtil.UcFirst(typ),
		})
		fmt.Println(code)
	}
}
