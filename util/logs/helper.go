package logs

import (
	"fmt"
	"github.com/3th1nk/easygo/util/mathUtil"
	"github.com/modern-go/reflect2"
	"time"
)

// IsIntLevelEnable 判断是否需要记录日志。
//   expect: 需要记录的日志等级
//   actual: 实际的日志等级
func IsIntLevelEnable(expect, actual int) bool {
	return expect <= actual
}

func IsLevelEnable(logger Logger, level int) bool {
	if reflect2.IsNil(logger) {
		return false
	} else if val := logger.GetLevel(); val != LevelNotSet {
		return val >= level
	}
	return true
}

func IsErrorEnable(logger Logger) bool {
	return IsLevelEnable(logger, LevelError)
}

func IsWarnEnable(logger Logger) bool {
	return IsLevelEnable(logger, LevelWarn)
}

func IsInfoEnable(logger Logger) bool {
	return IsLevelEnable(logger, LevelInfo)
}

func IsDebugEnable(logger Logger) bool {
	return IsLevelEnable(logger, LevelDebug)
}

func Write(logger Logger, level int, format string, a ...interface{}) {
	if v, _ := logger.(WritableLogger); v != nil {
		v.Write(level, format, a...)
		return
	}

	switch mathUtil.MinMaxInt(level, LevelError, LevelDebug) {
	case LevelDebug:
		logger.Debug(format, a...)
	case LevelInfo:
		logger.Info(format, a...)
	case LevelWarn:
		logger.Warn(format, a...)
	case LevelError:
		logger.Error(format, a...)
	}
}

func fmtStr(format string, a ...interface{}) string {
	if len(a) != 0 {
		return fmt.Sprintf(format, a...)
	}
	return format
}

func timePrefix(level int) string {
	return time.Now().Format("2006-01-02 15:04:05") + " [" + levelPrefix(level) + "] "
}

func levelPrefix(level int) string {
	switch level {
	case LevelDebug:
		return "D"
	case LevelInfo:
		return "I"
	case LevelWarn:
		return "W"
	case LevelError:
		return "E"
	case LevelFatal:
		return "C"
	default:
		if level > LevelDebug {
			return "D"
		} else if level < LevelFatal {
			return "C"
		}
	}
	return fmt.Sprintf("??? %v", level)
}
