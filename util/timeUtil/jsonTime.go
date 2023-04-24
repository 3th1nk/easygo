package timeUtil

import (
	"database/sql/driver"
	"fmt"
	"github.com/modern-go/reflect2"
	"strings"
	"time"
)

type JsonTime struct {
	time.Time
}

func ToJsonTime(t time.Time) JsonTime {
	return JsonTime{Time: t}
}

func (this JsonTime) String() string {
	return this.Format(GeneralFormatNano)
}

// MarshalJSON on LocalTime format Time field with %Y-%m-%d %H:%M:%S
func (this JsonTime) MarshalJSON() ([]byte, error) {
	if !this.After(Local2000) {
		return []byte(`""`), nil
	}
	str := `"` + this.Format(GeneralFormatNano) + `"`
	return []byte(str), nil
}

func (this *JsonTime) UnmarshalJSON(data []byte) error {
	if len(data) > 0 {
		str := strings.Trim(string(data), "\"")
		tmp, err := Parse(str)
		if reflect2.IsNil(err) {
			this.Time = tmp
		}
		return err
	}
	return nil
}

// Value insert timestamp into mysql need this function.
func (this JsonTime) Value() (driver.Value, error) {
	return this.Format(GeneralFormatNano), nil
}

// Scan valueof time.Time
func (this *JsonTime) Scan(v interface{}) (err error) {
	switch t := v.(type) {
	case time.Time:
		this.Time = t
	case string:
		this.Time, err = Parse(t)
	case []byte:
		this.Time, err = Parse(string(t))
	default:
		err = fmt.Errorf("can not convert %v to timestamp", v)
	}
	return
}

func (this JsonTime) Format(format ...string) string {
	return Format(this.Time, format...)
}
