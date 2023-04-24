package strUtil

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSplitN(t *testing.T) {
	str := "a, b, c, d"
	a, b := Split2(str, ",", true)
	assert.Equal(t, "a", a)
	assert.Equal(t, "b, c, d", b)
	t.Logf("Split2: [%v], [%v]", a, b)
	a, b, c := Split3(str, ",", true)
	assert.Equal(t, "a", a)
	assert.Equal(t, "b", b)
	assert.Equal(t, "c, d", c)
	t.Logf("Split2: [%v], [%v], [%v]", a, b, c)
	arr := Split(str, ",", true)
	assert.Equal(t, 4, len(arr))
	t.Logf("Split: [%v]", strings.Join(arr, ","))
}

func TestSplitToInt(t *testing.T) {
	var arr []int
	var err error

	arr, err = SplitToInt("", ",", true)
	if !reflect2.IsNil(err) {
		t.Error(fmt.Errorf("error occured: %v", err))
	} else if len(arr) != 0 {
		t.Error(fmt.Sprintf("assert faild"))
	}

	arr, err = SplitToInt("1,2,,3,", ",", true)
	if !reflect2.IsNil(err) {
		t.Error(fmt.Errorf("error occured: %v", err))
	} else if len(arr) != 3 || arr[0] != 1 || arr[1] != 2 || arr[2] != 3 {
		str, _ := jsonApi.MarshalToString(arr)
		t.Error(fmt.Sprintf("assert faild: %v", str))
	}
}
