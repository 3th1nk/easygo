package osUtil

import (
	"github.com/3th1nk/easygo/util/logs"
	"github.com/3th1nk/easygo/util/runtimeUtil"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	NoLogOnSignalExit bool

	exitSig  syscall.Signal
	once     sync.Once
	onSig    = make([]*sigFuncWrap, 0, 16)
	teardown []func()
)

func TearDown(f func()) {
	teardown = append(teardown, f)
}

// 获取已经收到的退出信号。
// 如果尚未收到退出信号，则返回 0
func GetSignal() syscall.Signal {
	return exitSig
}

// 当收到操作系统的退出信号时触发。
// 信号包括：SIGHUP, SIGINT, SIGQUIT, SIGILL, SIGABRT, SIGKILL, SIGTERM
func OnSignalExit(f func(sig syscall.Signal) (ok bool), msg ...string) {
	addSigWrap(&sigFuncWrap{
		f:       f,
		logging: true,
		msg:     append(msg, "")[0],
	})
}

// 当收到操作系统的退出信号时触发。
// 信号包括：SIGHUP, SIGINT, SIGQUIT, SIGILL, SIGABRT, SIGKILL, SIGTERM
func OnSignalExitNoLog(f func(sig syscall.Signal), msg ...string) {
	addSigWrap(&sigFuncWrap{
		f:       func(s syscall.Signal) (ok bool) { f(s); return true },
		logging: false,
		msg:     append(msg, "")[0],
	})
}

func addSigWrap(w *sigFuncWrap) {
	if w.logging {
		_, w.file, w.line, _ = runtimeUtil.Caller(2)
		if w.file == "" {
			w.file = "???"
		}
	}
	onSig = append(onSig, w)
}

func init() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		sig := <-ch
		exitSig = sig.(syscall.Signal)

		var logger logs.Logger
		if !NoLogOnSignalExit {
			if logger = logs.Default; logger == nil {
				logger = logs.Stdout()
			}
		}
		if logger != nil {
			logger.Info("recv signal: %v[%v]", exitSig.String(), int(exitSig))
		}

		for i := len(onSig) - 1; i >= 0; i-- {
			start := time.Now()
			v := onSig[i]
			ok := v.f(exitSig)
			if v.logging && logger != nil {
				took := time.Now().Sub(start).Round(time.Millisecond)
				if v.msg != "" {
					logger.Info("on signal exit: %v, ok=%v, took=%v, %v:%v", v.msg, ok, took, v.file, v.line)
				} else {
					logger.Info("on signal exit: ok=%v, took=%v, %v:%v", ok, took, v.file, v.line)
				}
			}
		}

		// TearDown
		if len(teardown) != 0 {
			for _, f := range teardown {
				f()
			}
			// 留一小段时间
			time.Sleep(20 * time.Millisecond)
		}

		os.Exit(int(exitSig))
	}()
}

type sigFuncWrap struct {
	f       func(syscall.Signal) (ok bool)
	logging bool
	msg     string
	file    string
	line    int
}
