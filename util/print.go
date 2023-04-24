package util

import (
	"fmt"
	"os"
	"time"
)

func Println(format string, a ...interface{}) {
	fmt.Println(toStr(format, a...))
}

func PrintArgsLn(a ...interface{}) {
	fmt.Println(a...)
}

func PrintTimeLn(format string, a ...interface{}) {
	fmt.Println(time.Now().Format("[15:04:05.000] ") + toStr(format, a...))
}

func PrintErrln(format string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, toStr(format, a...))
}

func PrintErrArgsLn(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func PrintErrTimeLn(format string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, time.Now().Format("[15:04:05.000] ")+toStr(format, a...))
}

func PrintErrLongTimeLn(format string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, time.Now().Format("[2006-01-02 15:04:05.000] ")+toStr(format, a...))
}

func toStr(format string, a ...interface{}) string {
	if len(a) != 0 {
		return fmt.Sprintf(format, a...)
	}
	return format
}
