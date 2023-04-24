package runtimeUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/util/logs"
	"github.com/modern-go/reflect2"
	"os"
	"reflect"
)

// 创建一个用于捕获并处理 Panic 的对象。
//   该对象总是会将错误信息及堆栈打印到日志中，除非通过修改 PanicLogger 屏蔽日志。
func Recover() *recoverHandler {
	return &recoverHandler{}
}

// 处理 recover：如果 e 不为空，则使用 logs.Default 打印堆栈信息
// 参数：
//   e: recover() 返回的对象。由于 recover 函数特性，无法在 HandleRecover 内部捕获上层函数的 panic 对象，所以必须由调用方传入。
func HandleRecover(msg string, e ...interface{}) interface{} {
	var p interface{}
	var handler reflect.Value
	var handlerIn []reflect.Value
	for _, v := range e {
		if v == nil {
			continue
		}
		rv := reflect.ValueOf(v)
		rt := rv.Type()
		if rt.Kind() == reflect.Func {
			switch rt.NumIn() {
			case 0:
				handler = rv
			case 1:
				if rt.In(0).Name() == "" {
					handler = rv
					handlerIn = make([]reflect.Value, 1)
				}
			}
		} else {
			p = v
		}
	}
	if p == nil {
		p = recover()
	}
	if p == nil {
		return nil
	}

	// 打印日志
	logPanicStack(1, p, msg)

	// 执行回调函数
	if handler.IsValid() {
		if len(handlerIn) == 1 {
			handlerIn[0] = reflect.ValueOf(p)
		}
		handler.Call(handlerIn)
	}

	return p
}

func PanicIfError(err error, msg ...string) {
	if err != nil {
		logPanicStack(1, err, msg...)
		panic(wrapError(err, msg...))
	}
}

func PanicIfNil(a interface{}, msg ...string) {
	if reflect2.IsNil(a) {
		msg = append(msg, "nil pointer")
		logPanicStack(1, msg[0])
		panic(msg[0])
	}
}

// 打印日志
func logPanicStack(skip int, err interface{}, msg ...string) {
	var head string
	if len(msg) != 0 && msg[0] != "" {
		head = fmt.Sprintf("[PANIC] %s: %v", msg[0], err)
	} else {
		head = fmt.Sprintf("[PANIC] %v", err)
	}
	stack := Stack(skip + 1)
	if PanicLogger != nil {
		stack.Log(PanicLogger, logs.LevelFatal, head)
	} else {
		WriteStack(os.Stderr, stack, head)
	}
}

var PanicLogger logs.Logger

type recoverHandler struct {
	msg   string
	err   *error
	f     func(p interface{})
	panic interface{}
}

// 设置打印 Panic 信息时的附加消息。
//
// 包含附加消息时的格式：
//   [PANIC] {$msg}: {$panic}
// 不当包含附加消息时的格式：
//   [PANIC]: {$panic}
func (this *recoverHandler) Msg(format string, a ...interface{}) *recoverHandler {
	if len(a) != 0 {
		this.msg = fmt.Sprintf(format, a...)
	} else {
		this.msg = format
	}
	return this
}

// 当发生 Panic 时，把 Panic 对象转化为 error 对象、并赋值给指定的 error 指针。
func (this *recoverHandler) Err(err *error) *recoverHandler {
	this.err = err
	return this
}

func (this *recoverHandler) OnPanic(f func(p interface{})) *recoverHandler {
	this.f = f
	return this
}

// 主处理函数。
func (this *recoverHandler) Handle(e ...interface{}) interface{} {
	if len(e) != 0 {
		this.panic = e[0]
	} else {
		this.panic = recover()
	}

	if this.panic == nil {
		return nil
	}

	// 处理 *error
	if this.err != nil {
		if e, _ := this.panic.(error); e != nil {
			*this.err = e
		} else {
			*this.err = fmt.Errorf("%v", this.panic)
		}
	}

	// 执行回调函数
	if this.f != nil {
		this.f(this.panic)
	}

	// 打印日志
	logPanicStack(1, this.panic, this.msg)

	return this.panic
}

func wrapError(err error, msg ...string) error {
	if len(msg) != 0 && msg[0] != "" {
		return fmt.Errorf("%v: %v", msg[0], err)
	}
	return err
}
