package timeUtil

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
	"math"
	"strings"
	"testing"
	"time"
)

func TestMaxMin(t *testing.T) {
	for _, item := range []time.Time{
		ZeroTime,
		Local2000,
		EndOf2037,
		MinValidTimestamp,
		MaxValidTimestamp,
		time.Date(1, 1, 1, 0, 0, 0, 0, time.Local),
		time.Unix(math.MinInt32, 0),
		time.Unix(math.MaxUint32, 0),
		time.Date(10000, 1, 1, 0, 0, 0, -1, time.Local),
	} {
		t.Log(Format(item), item.Unix(), item.Day(), item.Weekday())
	}
}

func TestAbc(t *testing.T) {
	now := time.Date(1980, 1, 1, 0, 0, 0, 0, time.Local)
	t.Log(math.MaxInt32, now.Unix(), ToMs(now))
}

func TestFormat(t *testing.T) {
	t.Log(Format(time.Unix(math.MaxInt32, 0)))
	t.Log(Format(time.Unix(0, 0)))
	t.Log(Format(time.Time{}))
	t.Log(Format(time.Unix(math.MinInt32, 0)))
	t.Log(Format(time.Date(9999, 12, 31, 23, 59, 60, -1, time.Local), DateTimeNano))
	t9999 := time.Date(9999, 12, 31, 23, 59, 60, -1, time.Local)
	t.Log(t9999.Unix())
	t.Log(time.Now().UnixNano())
	t.Log(math.MaxInt32)
	t.Log(time.Time{}.Unix())
}

func TestMax(t *testing.T) {
	now := time.Now()
	val := Max(now.Add(-time.Minute), now, now.Add(time.Minute))
	if val != now.Add(time.Minute) {
		t.Errorf("assert faild: %v, now=%v", val, now)
	}
}

func TestMin(t *testing.T) {
	now := time.Now()
	val := Min(now.Add(-time.Minute), now, now.Add(time.Minute))
	if val != now.Add(-time.Minute) {
		t.Errorf("assert faild: %v, now=%v", val, now)
	}
}

func TestVar(t *testing.T) {
	var str string

	str = Local2000.Format("2006-01-02 15:04:05.000000000")
	if str != "2000-01-01 00:00:00.000000000" {
		t.Error(fmt.Sprintf("assert faild: %v", str))
	}

	str = EndOf2037.Format("2006-01-02 15:04:05.000000000")
	if str != "2037-12-31 23:59:59.999999999" {
		t.Error(fmt.Sprintf("assert faild: %v", str))
	}
}

func TestToFloatStr(t *testing.T) {
	for i := 0; i < 10; i++ {
		now := time.Now()
		for d := -1; d < 10; d++ {
			t.Logf("%02d: %v", d, ToFloatStr(now, d))
		}
		time.Sleep(time.Millisecond * 50)
	}
}

func TestToMsStr(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(ToMsFloatStr(time.Now(), 6))
		time.Sleep(time.Microsecond * 10)
	}
}

func TestGeneralParse(t *testing.T) {
	str := "2019-6-4    12:35:36"
	val, err := Parse(str)
	if !reflect2.IsNil(err) {
		t.Errorf("error: %v", err)
	} else {
		t.Logf("%s --> %s", str, val.Format(DateTimeNano))
	}

	for _, str := range []string{
		"6-4",
		"6-04",
		"06-4",
		"06-04",
		"2019-6-4",
		"2019-06-4",
		"2019-6-04",
		"2019-06-04",
		"2019-06-0",
		"2019-06-00",
		"2019-6-4 12:13",
		"2019-6-4    12:13:36",
		"2019-6-4 12:35:36.111112000000000",
	} {
		val, err := Parse(str)
		if !reflect2.IsNil(err) {
			t.Errorf("error: %v", err)
		} else {
			t.Logf("%s --> %s", str, val.Format(DateTimeNano))
		}
	}
}

func TestLastWeekday(t *testing.T) {
	today := BeginOfDay(time.Now())
	for i := 0; i < 7; i++ {
		w := time.Weekday(i)
		last, next := LastWeekday(today, w), NextWeekday(today, w)
		if last.After(today) {
			t.Errorf("assert faild: last=%v", FormatDate(last))
		} else if !next.After(today) {
			t.Errorf("assert faild: last=%v", FormatDate(last))
		} else {
			t.Logf("%9v: last=%v, next=%v", w, FormatDate(last), FormatDate(next))
		}
	}
}

func TestBeginOfWeek(t *testing.T) {
	for i := 0; i < 10; i++ {
		now := time.Now().Add(time.Duration(i) * 24 * time.Hour)
		t.Logf("beginAtSunday:  date=%v, begin=%v, end=%v", Format(now, DateOnly), Format(BeginOfWeek(now, false)), Format(EndOfWeek(now, false)))
		t.Logf(strings.Repeat(" ", 15)+"                  begin=%v, end=%v", Format(BeginOfWeek(now, true)), Format(EndOfWeek(now, true)))
	}
}

func TestRoundAndTruncate(t *testing.T) {
	ts, _ := Parse("2022-01-01 00:00:00")
	rd := 5 * time.Second
	assert.Equal(t, ts.Format(DateTime), ts.Add(2*time.Second).Round(5*time.Second).Format(DateTime))
	assert.Equal(t, ts.Add(rd).Format(DateTime), ts.Add(3*time.Second).Round(5*time.Second).Format(DateTime))
	assert.Equal(t, ts.Format(DateTime), ts.Add(2*time.Second).Truncate(5*time.Second).Format(DateTime))
	assert.Equal(t, ts.Format(DateTime), ts.Add(3*time.Second).Truncate(5*time.Second).Format(DateTime))
}
