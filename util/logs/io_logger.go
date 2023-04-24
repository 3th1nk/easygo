package logs

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"io"
	"os"
)

func Stdout(level ...int) *IOLogger {
	return NewIOLogger(os.Stdout, level...)
}

func Stderr(level ...int) *IOLogger {
	return NewIOLogger(os.Stderr, level...)
}

func NewIOLogger(w io.Writer, level ...int) *IOLogger {
	return &IOLogger{W: w, level: util.IfEmptyIntSlice(level, 0)}
}

type IOLogger struct {
	W     io.Writer
	level int
}

func (this *IOLogger) GetLevel() int {
	return this.level
}

func (this *IOLogger) SetLevel(level int) {
	this.level = level
}

func (this *IOLogger) Fatal(format string, a ...interface{}) {
	_, _ = fmt.Fprintln(this.W, timePrefix(LevelFatal)+fmtStr(format, a...))
}

func (this *IOLogger) Error(format string, a ...interface{}) {
	_, _ = fmt.Fprintln(this.W, timePrefix(LevelError)+fmtStr(format, a...))
}

func (this *IOLogger) Warn(format string, a ...interface{}) {
	_, _ = fmt.Fprintln(this.W, timePrefix(LevelWarn)+fmtStr(format, a...))
}

func (this *IOLogger) Info(format string, a ...interface{}) {
	_, _ = fmt.Fprintln(this.W, timePrefix(LevelInfo)+fmtStr(format, a...))
}

func (this *IOLogger) Debug(format string, a ...interface{}) {
	_, _ = fmt.Fprintln(this.W, timePrefix(LevelDebug)+fmtStr(format, a...))
}

func (this *IOLogger) Write(level int, format string, a ...interface{}) {
	_, _ = fmt.Fprintln(this.W, timePrefix(level)+fmtStr(format, a...))
}
