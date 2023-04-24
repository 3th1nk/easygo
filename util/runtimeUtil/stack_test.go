package runtimeUtil

import (
	"github.com/3th1nk/easygo/util"
	"runtime/debug"
	"strings"
	"testing"
)

func TestStackString_1(t *testing.T) {
	a()
}

type tmp struct {
}

func a() {
	(&tmp{}).b()
}

func (this *tmp) b() {
	c()
}

func c() {
	func() {
		util.Println("=========================================== Stack:\n" + StackStr(0))

		b := debug.Stack()
		util.Println("=========================================== debug.Stack:\n" + string(b))
	}()
}

func TestStackString_2(t *testing.T) {
	stackString := strings.TrimSpace(`
runtime/debug.Stack(0x44da18f, 0x4a00600, 0xc00008e9f8)
        /usr/local/go/src/runtime/debug/stack.go:24 +0x9f
runtime/debug.PrintStack()
        /usr/local/go/src/runtime/debug/stack.go:16 +0x25
ops/center/alarm/handler.(*MetricCollectFailedAlarm).alarm.func1()
        /Users/hsfish/go/src/ops/center/alarm/handler/alarm_collect_failed.go:210 +0x25
panic(0x4a56940, 0x528c300)
        /usr/local/go/src/runtime/panic.go:969 +0x1b9
ops/center/alarm/handler.(*MetricCollectFailedAlarm).alarm(0xc0001be180, 0xd1d, 0x6141d94f, 0xc000654d80, 0x1, 0x1)
        /Users/hsfish/go/src/ops/center/alarm/handler/alarm_collect_failed.go:244 +0x7c3
ops/center/alarm/handler.(*MetricCollectFailedAlarm).execSelect.func1(0xc0001be180)
        /Users/hsfish/go/src/ops/center/alarm/handler/alarm_collect_failed.go:201 +0x174
created by ops/center/alarm/handler.(*MetricCollectFailedAlarm).execSelect
        /Users/hsfish/go/src/ops/center/alarm/handler/alarm_collect_failed.go:198 +0x5e
`)
	info := debugStringToStack(stackString, 0)
	for _, frame := range info.Frames {
		util.Println(frame.String())
	}
}
