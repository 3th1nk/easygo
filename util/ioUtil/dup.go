package ioUtil

import (
	"os"
	"path/filepath"
)

// DupStderr 重定向 stderr 到指定文件。主要用于捕获进程异常退出的Panic日志信息
func DupStderr(filename string, onResult func(err error)) {
	baseDir := filepath.Dir(filename)
	if baseDir != "" && baseDir != "." {
		if err := MakeDir(baseDir); err != nil {
			onResult(err)
			return
		}
	}

	f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0666)
	err := dup(int(f.Fd()), int(os.Stderr.Fd()))
	onResult(err)
}
