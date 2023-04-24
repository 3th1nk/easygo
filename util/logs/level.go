package logs

import (
	"fmt"
	"strings"
)

// RFC5424 log message levels.
const (
	LevelOff    = -1
	LevelNotSet = 0
	LevelFatal  = 1
	LevelError  = 3
	LevelWarn   = 4
	LevelInfo   = 6
	LevelDebug  = 7
	LevelAll    = 9
)

func StrToLevel(str string) int {
	switch strings.ToLower(str) {
	case "all":
		return LevelAll
	case "debug":
		return LevelDebug
	case "info", "information":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	case "off", "none":
		return LevelOff
	}
	return LevelNotSet
}

func LevelToStr(level int) string {
	switch level {
	case LevelAll:
		return "all"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelOff:
		return "off"
	default:
		return fmt.Sprintf("level-%d", level)
	}
}
