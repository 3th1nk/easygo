package runtimeUtil

// 安全的调用一个方法。
// 返回值：
//   panic: 方法执行期间产生的 panic 对象
func CallFunc(f func(), msg ...string) (panic interface{}) {
	defer func() {
		var theMsg string
		if len(msg) != 0 {
			theMsg = msg[0]
		}
		panic = HandleRecover(theMsg, recover())
	}()
	f()
	return
}

// 安全的调用一个方法，并返回方法执行期间产生的 panic
// 返回值：
//   err: 方法返回的 error 对象
//   panic: 方法执行期间产生的 panic 对象
func CallFuncE(f func() (err error), msg ...string) (err error, panic interface{}) {
	defer func() {
		var theMsg string
		if len(msg) != 0 {
			theMsg = msg[0]
		}
		panic = HandleRecover(theMsg, recover())
	}()
	err = f()
	return
}

// 安全的调用一个方法，并返回方法执行期间产生的 panic
// 返回值：
//   err: 方法返回的 error 对象
//   panic: 方法执行期间产生的 panic 对象
func CallFuncI(f func() (i interface{}), msg ...string) (i interface{}, panic interface{}) {
	defer func() {
		var theMsg string
		if len(msg) != 0 {
			theMsg = msg[0]
		}
		panic = HandleRecover(theMsg, recover())
	}()
	i = f()
	return
}

// 安全的调用一个方法，并返回方法执行期间产生的 panic
// 返回值：
//   err: 方法返回的 error 对象
//   panic: 方法执行期间产生的 panic 对象
func CallFuncIE(f func() (i interface{}, err error), msg ...string) (i interface{}, err error, panic interface{}) {
	defer func() {
		var theMsg string
		if len(msg) != 0 {
			theMsg = msg[0]
		}
		panic = HandleRecover(theMsg, recover())
	}()
	i, err = f()
	return
}
