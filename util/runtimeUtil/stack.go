package runtimeUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/util/logs"
	"io"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

var (
	StackFileFormatter     func(file string) (newFile string)
	StackFunctionFormatter func(funcName string) (newFuncName string)

	filePattern = regexp.MustCompile(`^\s*(.*):([\d]+)\s*(.+)?`)
	funcPattern = regexp.MustCompile(`^\s*(.*[\w-_]+)\(([\w0-9.,\s]+)?\)`)
)

func StackStr(skip int, f ...func(frame *Frame) string) string {
	return StackStrN(skip, -1, f...)
}

func StackStrN(skip, limit int, f ...func(frame *Frame) string) string {
	var theF func(frame *Frame) string
	if len(f) != 0 && f[0] != nil {
		theF = f[0]
	} else {
		theF = func(frame *Frame) string {
			return frame.String()
		}
	}

	stack := Stack(skip, limit)
	lines := make([]string, len(stack.Frames))
	for i, f := range stack.Frames {
		lines[i] = theF(f)
	}
	return strings.Join(lines, "\n")
}

func Caller(skip int) (fn, file string, line int, ok bool) {
	var pc uintptr
	if pc, file, line, ok = runtime.Caller(skip + 1); ok {
		fn = runtime.FuncForPC(pc).Name()
		if StackFileFormatter != nil {
			file = StackFileFormatter(file)
		} else {
			file = defaultStackFileFormatter(file)
		}
		if StackFunctionFormatter != nil {
			fn = StackFunctionFormatter(fn)
		} else {
			fn = defaultStackFunctionFormatter(fn)
		}
	}
	return
}

func CallerFunc(skip int) string {
	fn, _, _, _ := Caller(skip + 1)
	return fn
}

func CallerFile(skip int) string {
	_, file, _, _ := Caller(skip + 1)
	return file
}

func CallerFileLine(skip int) string {
	_, file, line, _ := Caller(skip + 1)
	return fmt.Sprintf("%v:%v", file, line)
}

func Stack(skip int, limit ...int) *StackInfo {
	return debugStringToStack(string(debug.Stack()), skip, limit...)
}

func WriteStack(w io.Writer, stack *StackInfo, msg ...string) {
	var theMsg string
	if len(msg) != 0 {
		theMsg = msg[0]
	}
	stack.FPrints(w, theMsg)
}

func debugStringToStack(debugString string, skip int, limit ...int) *StackInfo {
	lines := strings.Split(strings.TrimSpace(debugString), "\n")

	stack := &StackInfo{Frames: make([]*Frame, 0, len(lines)/2+1)}
	if strings.HasPrefix(lines[0], "goroutine ") {
		if arr := strings.Split(lines[0], " "); len(arr) > 2 {
			stack.Goroutine, _ = strconv.Atoi(arr[1])
		}
		lines = lines[1:]
	}

	start, end := 0, -1
	for i, str := range lines {
		if i%2 == 1 {
			if 0 == start && !strings.Contains(str, "/runtime/panic.go") && !strings.Contains(str, "/runtime/debug/stack.go") && !strings.Contains(str, "/util/runtimeUtil/stack.go") {
				start = i - 1
			}
			if -1 == end && -1 != strings.Index(str, "/testing/testing.go") {
				end = i - 1
			}
		} else {
			if -1 == end && -1 != strings.Index(str, "main.main()") {
				end = i
			}
		}
	}
	if start != 0 && end != -1 {
		lines = lines[start:end]
	} else if start != -1 {
		lines = lines[start:]
	} else {
		lines = lines[:end]
	}

	var f *Frame
	for i, line := range lines {
		if i%2 == 0 {
			if f != nil {
				stack.Frames = append(stack.Frames, f)
			}
			f = &Frame{}

			matches := funcPattern.FindStringSubmatch(line)
			if len(matches) >= 3 {
				fn := matches[1]
				if StackFunctionFormatter != nil {
					fn = StackFunctionFormatter(fn)
				} else {
					fn = defaultStackFunctionFormatter(fn)
				}
				f.Func = fn
				f.Args = matches[2]
			} else {
				f.Func = line
			}
		} else if f != nil {
			matches := filePattern.FindStringSubmatch(line)
			if n := len(matches); n >= 3 {
				file := matches[1]
				if StackFileFormatter != nil {
					file = StackFileFormatter(file)
				} else {
					file = defaultStackFileFormatter(file)
				}
				f.File = file
				f.Line, _ = strconv.Atoi(matches[2])
				if n > 3 {
					f.Entry = matches[3]
				}
			}
		}
	}
	if f != nil {
		stack.Frames = append(stack.Frames, f)
	}

	n := len(stack.Frames)
	if skip > 0 {
		if skip < n {
			stack.Frames = stack.Frames[skip:]
		} else {
			stack.Frames = []*Frame{}
		}
	}
	if len(limit) != 0 && limit[0] > 0 && limit[0] < n {
		stack.Frames = stack.Frames[:limit[0]]
	}

	return stack
}

type StackInfo struct {
	Goroutine int
	Frames    []*Frame
}

type Frame struct {
	Func  string `json:"func,omitempty"`
	Args  string `json:"args,omitempty"`
	File  string `json:"file,omitempty"`
	Line  int    `json:"line,omitempty"`
	Entry string `json:"entry,omitempty"`
}

func (this *StackInfo) Log(logger logs.Logger, level int, msg ...string) {
	if logs.IsLevelEnable(logger, level) {
		var theMsg string
		if len(msg) != 0 {
			theMsg = msg[0]
		}
		str := this.ToString(theMsg)
		logs.Write(logger, level, str)
	}
}

func (this *StackInfo) Logs(logger logs.Logger, level int, msg string, f ...func(f *Frame) (string, bool)) {
	if logs.IsLevelEnable(logger, level) {
		str := this.ToString(msg, f...)
		logs.Write(logger, level, str)
	}
}

func (this *StackInfo) Print(f ...func(f *Frame) (string, bool)) {
	this.FPrints(os.Stdout, "", f...)
}

func (this *StackInfo) Prints(msg string, f ...func(f *Frame) (string, bool)) {
	this.FPrints(os.Stdout, msg, f...)
}

func (this *StackInfo) FPrint(w io.Writer, f ...func(f *Frame) (string, bool)) {
	this.FPrints(w, "", f...)
}

func (this *StackInfo) FPrints(w io.Writer, msg string, f ...func(f *Frame) (string, bool)) {
	str := this.ToString(msg, f...)
	_, _ = w.Write([]byte(str))
}

func (this *StackInfo) ToString(msg string, f ...func(f *Frame) (string, bool)) string {
	lines, cnt := make([]string, len(this.Frames)+1), 0
	if msg != "" {
		lines[cnt], cnt = msg, cnt+1
	}

	if len(f) == 0 || f[0] == nil {
		for _, frame := range this.Frames {
			lines[cnt], cnt = frame.String(), cnt+1
		}
	} else {
		for _, frame := range this.Frames {
			str, ok := f[0](frame)
			if str != "" {
				lines[cnt], cnt = str, cnt+1
			}
			if !ok {
				break
			}
		}
	}
	return strings.Join(lines[:cnt], "\n")
}

func (this *Frame) String() string {
	return fmt.Sprintf("%v(%v)\n\t%v:%v, %v", this.Func, this.Args, this.File, this.Line, this.Entry)
}

func defaultStackFileFormatter(file string) (newFile string) {
	if pos := strings.Index(file, "easygo"); pos != -1 {
		return file[pos+24:]
	} else if pos := strings.Index(file, "/gocommon"); pos != -1 {
		return file[pos+1:]
	}
	if pos := strings.Index(file, "/pkg/mod/"); pos != -1 {
		return file[pos+9:]
	}
	return file
}

func defaultStackFunctionFormatter(funcName string) (newFuncName string) {
	pos := strings.LastIndex(funcName, "/")
	if pos != -1 {
		return funcName[pos+1:]
	}
	return funcName
}
