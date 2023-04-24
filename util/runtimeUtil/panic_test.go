package runtimeUtil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecoverHandle(t *testing.T) {
	assert.NotPanics(t, func() {
		func(a *testing.T) {
			defer Recover().Handle()
			a.Log("trigger panic")
		}(nil)
	})

	assert.Panics(t, func() {
		func(a *testing.T) {
			defer func() {
				if e := Recover().Handle(); e != nil {
					// defer 的匿名函数中再调用 testHandler 是捕获不到的
					// 程序一定会 panic，不会到这里
					t.Errorf("不会到这里")
				}
			}()
			a.Log("trigger panic")
		}(nil)
	})

	assert.NotPanics(t, func() {
		func(a *testing.T) {
			defer func() {
				if e := Recover().Handle(recover()); e == nil {
					// recover 能够捕获到 panic，然后传入 Handle。故 Handle 结果非空。
					t.Errorf("不会到这里")
				}
			}()
			a.Log("trigger panic")
		}(nil)
	})
}

func TestRecoverErr(t *testing.T) {
	err := func(a *testing.T) (err error) {
		defer Recover().Err(&err).Handle()
		a.Log("trigger panic")
		return nil
	}(nil)
	assert.Error(t, err, "函数返回的 error 对象不应该为空")
}

func TestRecoverFunc(t *testing.T) {
	var thePanic interface{}
	err := func(a *testing.T) (err error) {
		defer Recover().Err(&err).OnPanic(func(p interface{}) {
			thePanic = p
		}).Handle()
		a.Log("trigger panic")
		return nil
	}(nil)
	assert.Error(t, err)
	assert.NotNil(t, thePanic)
}

func ExampleHandleRecover() {
	func() {
		defer HandleRecover("出错啦")
		// do something
	}()

	func() {
		defer func() {
			// 由于此处是在 defer 匿名函数内部，必须在参数中传递 recover 捕获的结果。
			// 否则在 HandleRecover 内部将由于嵌套层级过深而无法捕获 panic。
			if panicErr := HandleRecover("出错啦", recover()); panicErr != nil {
				// ...
			}
		}()
		// do something
	}()
}

func ExampleRecover() {
	func(a *testing.T) {
		// 成功，能捕获到 recover
		defer Recover().Handle()
		a.Log("trigger panic")
	}(nil)

	func(a *testing.T) {
		defer func() {
			// 成功，在匿名函数中，先 recover 然后传递给 Handle。
			// Handle 结束后会将 recover 捕获到的对象重新返回。
			p := Recover().Handle(recover())
			print(p) // not nil
		}()
		a.Log("trigger panic")
	}(nil)

	func(a *testing.T) {
		defer func() {
			// 失败:
			// 由于在 defer 后的匿名函数内部，层次太深 Handle 内无法捕获 recover。
			// 如果要在匿名函数中使用，则必须使用 Handle(recover()) 代替。
			p := Recover().Handle()
			print(p) // nil
		}()
		a.Log("trigger panic")
	}(nil)
}
