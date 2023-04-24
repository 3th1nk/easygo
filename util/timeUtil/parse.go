package timeUtil

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var knownFormat = map[string]time.Duration{
	"year":        0,
	"month":       0,
	"day":         24 * time.Hour,
	"d":           24 * time.Hour,
	"hour":        time.Hour,
	"h":           time.Hour,
	"minute":      time.Minute,
	"min":         time.Minute,
	"m":           time.Minute,
	"second":      time.Second,
	"sec":         time.Second,
	"s":           time.Second,
	"millisecond": time.Millisecond,
	"ms":          time.Millisecond,
	"microsecond": time.Microsecond,
	"us":          time.Microsecond,
	"nanosecond":  time.Nanosecond,
	"ns":          time.Nanosecond,
}

var (
	timeStrPattern     = regexp.MustCompile(`^[\d-+.:tz\s]+$`)
	knownFormatPattern = regexp.MustCompile(`([+-]?\d+)\s*([a-z-]+)`)
)

func ParseNoErr(str string) (t time.Time) {
	t, _ = Parse(str)
	return
}

// 将字符串转换为 time.Time 类型。支持的格式如下：
//   数字:
//      如果小于 MaxInt32，则以秒为单位转换； 否则以毫秒为单位转换。
//   预定义格式:
//      now|today|yesterday|tomorrow: 当前时间、当天开始时间、昨天开始时间、明天开始时间
//      {n}{unit}: 当前时间 增加 {n} 个 {unit} 后的结果。 n 可以为负数、 unit 支持 year|month|day(d)|hour(h)|minute(min)|second(sec)|millisecond(ms)|microsecond(us)|nanosecond(ns)
//      以上格式可以通过英文逗号串联，比如 “tomorrow,2hour,-3min” 表示 “明天01:57”
//   [yyyy-]MM-dd 格式
//   [yyyy-]MM-dd HH:mm[:ss] 格式
func Parse(tmStr string, now ...time.Time) (time.Time, error) {
	str := strings.ToLower(strings.TrimSpace(tmStr))
	if str == "" {
		return ZeroTime, fmt.Errorf("empty str")
	}

	// 处理常量字符串 以及 (+-)(N)(Unit) 格式的字符串，例如 5day、-2hour
	var t time.Time
	if len(now) != 0 {
		t = now[0]
	} else {
		t = time.Now()
	}
	for _, s := range strings.Split(str, ",") {
		if s = strings.TrimSpace(s); s == "" {
			continue
		}

		// 如果是数字，则按照秒或者毫秒处理
		if n, err := strconv.ParseInt(s, 10, 64); reflect2.IsNil(err) {
			// 4102416000 是 2100-01-01 对应的秒数。
			// 根据数值范围，自动用 秒|毫秒|微妙|纳秒 为单位进行转换。
			if n < 4102416000 {
				// 如果小于 ‘2100-1-1’，按秒处理
				t = time.Unix(n, 0)
			} else if n < 4102416000000 {
				t = time.Unix(0, n*1000000)
			} else if n < 4102416000000000 {
				t = time.Unix(0, n*1000)
			} else {
				t = time.Unix(0, n)
			}
			continue
		}

		switch s {
		case "now":
			t = time.Now()
		case "today":
			t = BeginOfDay(t)
		case "yesterday":
			t = BeginOfDay(t).Add(-24 * time.Hour)
		case "tomorrow":
			t = BeginOfDay(t).Add(24 * time.Hour)
		case "month-start":
			t = BeginOfMonth(t)
		case "month-end":
			t = EndOfMonth(t)
		case "year-start":
			t = BeginOfYear(t)
		case "year-end":
			t = EndOfYear(t)
		case "week-start":
			t = LastWeekday(t, DefaultBeginOfWeek)
		case "week-end":
			t = NextWeekday(t, DefaultBeginOfWeek).Add(-1)
		case "week-start-sunday":
			t = LastWeekday(t, time.Sunday)
		case "week-start-monday":
			t = LastWeekday(t, time.Monday)
		case "week-end-saturday":
			t = NextWeekday(t, time.Sunday).Add(-1)
		case "week-end-sunday":
			t = NextWeekday(t, time.Monday).Add(-1)
		default:
			// 如果是时间字符串，则 parseTimeStr; 否则判断是否是 knownFormatPattern
			if timeStrPattern.MatchString(s) {
				var err error
				if t, err = parseTimeStr(s); !reflect2.IsNil(err) {
					return time.Time{}, fmt.Errorf("无法识别的格式: '%v'", tmStr)
				}
			} else {
				matches := knownFormatPattern.FindStringSubmatch(s)
				if len(matches) == 0 {
					return time.Time{}, fmt.Errorf("无法识别的格式: '%v'", tmStr)
				}

				unit, _ := strconv.ParseInt(matches[1], 10, 64)
				switch matches[2] {
				case "year":
					t = time.Date(t.Year()+int(unit), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
				case "month":
					t = time.Date(t.Year(), t.Month()+time.Month(unit), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
				default:
					if duration, ok := knownFormat[matches[2]]; !ok {
						return time.Time{}, fmt.Errorf("无法识别的格式: '%v'", tmStr)
					} else {
						t = t.Add(time.Duration(unit) * duration)
					}
				}
			}
		}
	}
	return t, nil
}

func parseTimeStr(str string) (time.Time, error) {
	formatError := fmt.Errorf("格式不正确: %v", str)

	// 兼容RFC3339，转换成Go标准时间格式 2006-01-02 15:04:05
	// 2006-01-02T15:04:05.00Z		// 带T
	// 2006-01-02 15:04:05.00Z		// 不带T
	// 2006-01-02T15:04:05+07:00    // 时区指示符: Z或±hh:mm
	var loc = time.Local
	str = strings.Replace(str, "t", " ", 1)
	if strings.HasSuffix(str, "z") {
		str = strings.TrimSuffix(str, "z")
		loc = time.UTC
	}

	var dateStr, timeStr, msStr string
	if pos := strings.Index(str, " "); pos != -1 {
		dateStr = str[:pos]
		ss := strings.TrimSpace(str[pos+1:])
		// 处理timezone
		// ±hh:mm
		for _, sign := range []string{"+", "-"} {
			if pos = strings.Index(ss, sign); pos != -1 {
				tzStr := strings.Split(ss[pos+1:], ":")

				var hour, min int
				hour, err := strconv.Atoi(tzStr[0])
				if err != nil {
					return time.Time{}, formatError
				}
				if len(tzStr) > 1 && tzStr[1] != "" {
					min, err = strconv.Atoi(tzStr[1])
					if err != nil {
						return time.Time{}, formatError
					}
				}
				offset := hour*3600 + min*60
				if sign == "-" {
					offset *= -1
				}
				loc = time.FixedZone("", offset)

				ss = strings.TrimSpace(ss[:pos])
				break
			}
		}

		if pos = strings.Index(ss, "."); pos != -1 {
			timeStr = ss[:pos]
			msStr = ss[pos+1:]
		} else {
			timeStr = ss
		}
	} else {
		dateStr = str
	}

	var yearStr, monthStr, dayStr, hourStr, minStr, secStr, nsecStr string
	tmp := strings.Split(dateStr, "-")
	switch len(tmp) {
	case 3:
		yearStr, monthStr, dayStr = tmp[0], tmp[1], tmp[2]
	case 2:
		monthStr, dayStr = tmp[0], tmp[1]
	}
	if timeStr != "" {
		tmp = strings.Split(timeStr, ":")
		switch len(tmp) {
		case 3:
			hourStr, minStr, secStr = tmp[0], tmp[1], tmp[2]
			if msStr != "" {
				n := len(msStr)
				if n == 9 {
					nsecStr = msStr
				} else if n > 9 {
					nsecStr = msStr[:9]
				} else {
					nsecStr = msStr + strings.Repeat("0", 9-n)
				}
			}
		case 2:
			hourStr, minStr = tmp[0], tmp[1]
		}
	}

	var year, month, day, hour, min, sec, nsec int
	var err error
	if yearStr != "" {
		if year, err = strconv.Atoi(yearStr); !reflect2.IsNil(err) {
			return time.Time{}, formatError
		}
	}
	if monthStr == "" || dayStr == "" {
		return time.Time{}, formatError
	} else {
		if month, err = strconv.Atoi(monthStr); !reflect2.IsNil(err) {
			return time.Time{}, formatError
		}
		if day, err = strconv.Atoi(dayStr); !reflect2.IsNil(err) {
			return time.Time{}, formatError
		}
	}
	if hourStr != "" || minStr != "" || secStr != "" {
		if hour, err = strconv.Atoi(hourStr); !reflect2.IsNil(err) {
			return time.Time{}, formatError
		}
		if min, err = strconv.Atoi(minStr); !reflect2.IsNil(err) {
			return time.Time{}, formatError
		}
		if secStr != "" {
			if sec, err = strconv.Atoi(secStr); !reflect2.IsNil(err) {
				return time.Time{}, formatError
			}
			if nsecStr != "" {
				if nsec, err = strconv.Atoi(nsecStr); !reflect2.IsNil(err) {
					return time.Time{}, formatError
				}
			}
		}
	}
	if year == 0 {
		year = time.Now().Year()
	}

	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.Local).In(loc), nil
}
