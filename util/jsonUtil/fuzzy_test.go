package jsonUtil

import (
	"github.com/3th1nk/easygo/util"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestFuzzyTime(t *testing.T) {
	type temp struct {
		T1 time.Time `json:"t1,omitempty"`
		T2 time.Time `json:"t2"`
	}

	a := temp{T1: time.Time{}, T2: time.Time{}}
	str := MustMarshalToString(a)
	util.Println(str)
	assert.False(t, strings.Contains(str, "t1"), str)

	a.T1 = time.Now()
	a.T2 = time.Now()
	str = MustMarshalToString(a)

	b := a
	b.T1 = b.T1.Add(time.Hour)
	b.T2 = b.T2.Add(time.Hour)
	if assert.NoError(t, UnmarshalFromString(str, &b), str) {
		assert.Equal(t, a.T1.Format(DefaultTimeFormat), b.T1.Format(DefaultTimeFormat))
		assert.Equal(t, a.T2.Format(DefaultTimeFormat), b.T2.Format(DefaultTimeFormat))
	}
}
