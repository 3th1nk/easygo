package timeUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestParseOne(t *testing.T) {
	str := "0001-01-01 00:00:00"
	val, err := Parse(str)
	if assert.NoError(t, err) {
		assert.Equal(t, str, val.Format(GeneralFormatNano))
		util.Println(val.Format(GeneralFormatNano))
	}
}

func TestParse(t *testing.T) {
	now := time.Now().Add(time.Hour * 24 * 365 * 100)
	t.Logf("now=%v, %v", now, now.UnixNano())
	for _, v := range []interface{}{
		1563362142,
		1563362142564,
		1563362142564325,
		1563362142564325365,
		"now",
		"today",
		"tomorrow",
		"2hour",
		"+2hour",
		"-2hour",
		"2 hour",
		"+2 hour",
		"-2 hour",
		"today,+2hour",
		"today,+2 hour",
		"today,2hour",
		"today,2 hour",
		"today,+2h",
		"today,+2 h",
		"today,2h",
		"today,2 h",
		"today,-2hour",
		"today,-2 hour",
		"today,-2h",
		"today,-2 h",
		"today,-2d,3h",
		"tomorrow,-1ns",
		"2019-07-17 19:15:42",
		"2019-07-17 19:15",
		"2019-07-17",
		"07-17 19:15",
		"2019-07-17 19:15:42,+2h,-1ns",
		"week-start",
		"week-end",
		"week-start-monday",
		"week-start-sunday",
		"week-end-sunday",
		"week-end-saturday",
		"month-start",
		"month-start,-1month",
		"month-end,-7month",
		"month-start,-6month",
		"2023-03-16T08:00:00Z",
		"2023-03-16t08:00:00z",
		"2023-03-16T08:00:00.000Z",
		"2023-03-16t08:00:00+07:31",
		"2023-03-16 08:00:00-03:00",
	} {
		if val, err := Parse(fmt.Sprint(v)); !reflect2.IsNil(err) {
			t.Error(err)
		} else {
			t.Logf("%v: %v", v, val.String())
		}
	}
}

func TestParse_2(t *testing.T) {
	for _, v := range [][]string{
		{"1580486400,", "2020-02-01 00:00:00"},
		{"1580486400,tomorrow", "2020-02-02 00:00:00"},
		{"1580486400,tomorrow,2h,3min", "2020-02-02 02:03:00"},
		{"1580486400,yesterday,2h,3min", "2020-01-31 02:03:00"},
	} {
		val, err := Parse(v[0])
		if assert.NoError(t, err) {
			assert.Equal(t, v[1], val.Format(GeneralFormat))

			val, err = Parse(strings.ReplaceAll(v[0], "1580486400", "2020-02-01"))
			assert.NoError(t, err)
			assert.Equal(t, v[1], val.Format(GeneralFormat))
		}
	}
}
