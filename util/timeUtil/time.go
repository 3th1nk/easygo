package timeUtil

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	Day  = time.Hour * 24
	Week = Day * 7

	DateTime     = "2006-01-02 15:04:05"
	DateTimeNano = "2006-01-02 15:04:05.999999999"
	DateOnly     = "2006-01-02"
)

var (
	DefaultFormat      = DateTime
	DefaultBeginOfWeek = time.Monday

	ZeroTime          = time.Unix(0, 0)
	Greenwich1970     = ZeroTime
	MinValidTimestamp = time.Unix(1, 0)                                // timestamp 的最小有效值
	MaxValidTimestamp = time.Unix(2147483647, 0)                       // timestamp 的最大有效值
	Local2000         = time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)  // 本地时区时间 2000-01-01 00:00:00.00000000
	EndOf2037         = time.Date(2038, 1, 1, 0, 0, 0, -1, time.Local) // 本地时区时间 2037-12-31 23:59:59.99999999
	Local2000Sec      = Local2000.Unix()                               // Unix of 本地时区时间 2000-01-01 00:00:00.00000000
	EndOf2037Sec      = EndOf2037.Unix()                               // Unix of 本地时区时间 2037-12-31 23:59:59.99999999
	Local2000Ms       = ToMs(Local2000)                                // Millisecond of 本地时区时间 2000-01-01 00:00:00.00000000
	EndOf2037Ms       = ToMs(EndOf2037)                                // Millisecond of 本地时区时间 2037-12-31 23:59:59.99999999
	Local2000Ns       = Local2000.UnixNano()                           // Nanosecond of 本地时区时间 2000-01-01 00:00:00.00000000
)

// 判断一个时间是否在 timestamp 范围内。
// return:
//   0: 在范围内
//  -1: 不在范围内，低于最小值
//   1: 不在范围内，高于最大值
func IsValidTimestamp(t time.Time) int {
	if t.Before(MinValidTimestamp) {
		return -1
	} else if t.After(MaxValidTimestamp) {
		return 1
	}
	return 0
}

func Max(a ...time.Time) time.Time {
	if len(a) != 0 {
		t := a[0]
		for i, v := range a {
			if i != 0 && v.After(t) {
				t = v
			}
		}
		return t
	}
	return time.Date(10000, 1, 1, 0, 0, 0, -1, time.Local)
}

func Min(a ...time.Time) time.Time {
	if len(a) != 0 {
		t := a[0]
		for i, v := range a {
			if i != 0 && v.Before(t) {
				t = v
			}
		}
		return t
	}
	return time.Time{}
}

// 将 time.Time 转换为以秒为单位的浮点数
func ToFloat(t time.Time, decimals ...int) float64 {
	theDecimals := -1
	if len(decimals) != 0 {
		theDecimals = decimals[0]
	}
	if theDecimals == 0 {
		return float64(t.Unix())
	} else {
		v, _ := strconv.ParseFloat(strconv.FormatFloat(float64(t.UnixNano())/float64(time.Second), 'f', theDecimals, 64), 64)
		return v
	}
}

func ToFloatStr(t time.Time, decimals ...int) string {
	theDecimals := -1
	if len(decimals) != 0 {
		theDecimals = decimals[0]
	}
	if theDecimals == 0 {
		return strconv.FormatInt(t.Unix(), 10)
	} else {
		str := strconv.FormatFloat(float64(t.UnixNano())/float64(time.Second), 'f', theDecimals, 64)
		if theDecimals > 0 {
			str = strings.TrimRight(strings.TrimRight(str, "0"), ".")
		}
		return str
	}
}

// 将 time.Time 转换为以毫秒为单位的时间戳整数
func ToMs(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// 将 time.Time 转换为以毫秒为单位的时间戳浮点数
func ToMsFloat(t time.Time, decimals ...int) float64 {
	theDecimals := -1
	if len(decimals) != 0 {
		theDecimals = decimals[0]
	}
	if theDecimals == 0 {
		return float64(t.UnixNano() / int64(time.Millisecond))
	} else {
		v, _ := strconv.ParseFloat(strconv.FormatFloat(float64(t.UnixNano())/float64(time.Millisecond), 'f', theDecimals, 64), 64)
		return v
	}
}

func ToMsFloatStr(t time.Time, decimals ...int) string {
	theDecimals := -1
	if len(decimals) != 0 {
		theDecimals = decimals[0]
	}
	if theDecimals == 0 {
		return strconv.FormatInt(t.UnixNano()/int64(time.Millisecond), 10)
	} else {
		str := strconv.FormatFloat(float64(t.UnixNano())/float64(time.Millisecond), 'f', theDecimals, 64)
		if theDecimals > 0 {
			str = strings.TrimRight(strings.TrimRight(str, "0"), ".")
		}
		return str
	}
}

func FromNS(n int64) time.Time {
	return time.Unix(n/int64(time.Second), n%int64(time.Second))
}

// 将以毫秒为单位的时间戳整数转换为 time.Time
func FromMs(n int64) time.Time {
	return time.Unix(0, n*1000000)
}

// 将以毫秒为单位的时间戳整数转换为 time.Time
func FromMsFloat(f float64) time.Time {
	return FromSecondFloat(f * 1000)
}

// 将以毫秒为单位的时间戳整数转换为 time.Time
func FromSecondFloat(f float64) time.Time {
	s, ns := math.Modf(f)
	return time.Unix(int64(s), int64(ns*float64(time.Second)))
}

// 获取一天开始的时间（00:00:00）
func BeginOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 获取一天结束的时间（23:59:59.999999999）
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, -1, t.Location())
}

// 获取一天指定小时数的时间（xx:00:00）
func HourOfDay(t time.Time, hour int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, t.Location())
}

// 获取指定时间位于一周的第几天（从 1 开始）
func DayOfWeek(t time.Time, beginAtMonday ...bool) int {
	if len(beginAtMonday) == 0 {
		beginAtMonday = append(beginAtMonday, DefaultBeginOfWeek == time.Monday)
	}
	if beginAtMonday[0] {
		if wd := t.Weekday(); wd != time.Sunday {
			return int(wd)
		} else {
			return 7
		}
	} else {
		return int(t.Weekday()) + 1
	}
}

// 获取一周的开始时间（周日为一周第一天）
func BeginOfWeek(t time.Time, beginAtMonday ...bool) time.Time {
	if len(beginAtMonday) == 0 {
		beginAtMonday = append(beginAtMonday, DefaultBeginOfWeek == time.Monday)
	}
	if beginAtMonday[0] {
		if t.Weekday() == time.Sunday {
			return time.Date(t.Year(), t.Month(), t.Day()-6, 0, 0, 0, 0, t.Location())
		} else {
			return time.Date(t.Year(), t.Month(), t.Day()+1-int(t.Weekday()), 0, 0, 0, 0, t.Location())
		}
	} else {
		return time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()), 0, 0, 0, 0, t.Location())
	}
}

// 获取一周的结束时间（周六为一周最后一天）
func EndOfWeek(t time.Time, beginAtMonday ...bool) time.Time {
	if len(beginAtMonday) == 0 {
		beginAtMonday = append(beginAtMonday, DefaultBeginOfWeek == time.Sunday)
	}
	if beginAtMonday[0] {
		if t.Weekday() == time.Sunday {
			return time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, -1, t.Location())
		} else {
			return time.Date(t.Year(), t.Month(), t.Day()+8-int(t.Weekday()), 0, 0, 0, -1, t.Location())
		}
	} else {
		return time.Date(t.Year(), t.Month(), t.Day()+7-int(t.Weekday()), 0, 0, 0, -1, t.Location())
	}
}

func LastWeekday(t time.Time, w time.Weekday) time.Time {
	if tw := t.Weekday(); tw == w {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	} else {
		days := (7 + tw - w) % 7
		return time.Date(t.Year(), t.Month(), t.Day()-int(days), 0, 0, 0, 0, t.Location())
	}
}

func NextWeekday(t time.Time, w time.Weekday) time.Time {
	if tw := t.Weekday(); tw == w {
		return time.Date(t.Year(), t.Month(), t.Day()+7, 0, 0, 0, 0, t.Location())
	} else {
		days := (7 + w - tw) % 7
		return time.Date(t.Year(), t.Month(), t.Day()+int(days), 0, 0, 0, 0, t.Location())
	}
}

func BeginOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, -1, t.Location())
}

func BeginOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January+12, 1, 0, 0, 0, -1, t.Location())
}

func Milliseconds(d time.Duration, decimals ...int) float64 {
	if len(decimals) == 0 || decimals[0] == 0 {
		return float64(d) / float64(time.Millisecond)
	} else {
		v, _ := strconv.ParseFloat(strconv.FormatFloat(float64(d)/float64(time.Millisecond), 'f', decimals[0], 64), 64)
		return v
	}
}

func Format(t time.Time, format ...string) string {
	var realFormat string
	if len(format) != 0 {
		realFormat = format[0]
	}
	if realFormat == "" {
		realFormat = DefaultFormat
	}
	return t.Format(realFormat)
}

func FormatDate(t time.Time, format ...string) string {
	var realFormat string
	if len(format) != 0 {
		realFormat = format[0]
	}
	if realFormat == "" {
		realFormat = DateOnly
	}
	return t.Format(realFormat)
}

func Rand(d time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(d)))
}

func CurrentTimeStr(format ...string) string {
	return Format(time.Now(), format...)
}

func CurrentDateStr() string {
	return FormatDate(time.Now())
}
