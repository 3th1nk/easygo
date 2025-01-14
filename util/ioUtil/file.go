package ioUtil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(filename), perm); err != nil {
		return err
	}

	return os.WriteFile(filename, data, perm)
}

func WriteFileF(filename string, perm os.FileMode, f func(f *os.File) error) error {
	if err := os.MkdirAll(filepath.Dir(filename), perm); err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	err = f(file)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	return err
}

func IsFileExists(path string) bool {
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		return true
	}
	return false
}

func IsDirExists(dir string) bool {
	if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
		return true
	}
	return false
}

func MakeDir(dir string) error {
	if !IsDirExists(dir) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return fmt.Errorf("dir create failed %s", err.Error())
		}
	}
	return nil
}

// OpenFile 打开文件，如果文件不存在 或者 上层文件夹不存在，都进行创建
//		filename 文件名称
//		flags 打开文件模式，默认 追加读写
func OpenFile(filename string, flags ...int) (*os.File, error) {
	flag := os.O_CREATE | os.O_RDWR | os.O_APPEND
	if len(flags) > 0 {
		tmpFlag := os.O_RDONLY
		for _, f := range flags {
			tmpFlag = tmpFlag | f
		}
		flag = tmpFlag
	}

	if info, err := os.Stat(filename); err == nil {
		if info.IsDir() {
			return nil, fmt.Errorf("create file failed, %s is dir", filename)
		} else {
			return os.OpenFile(filename, flag, info.Mode())
		}
	}

	if err := MakeDir(filepath.Dir(filename)); err != nil {
		return nil, err
	}

	return os.OpenFile(filename, flag, os.ModePerm)
}

// Move 移动文件 首先通过os.Rename移动，移动失败，再尝试通过读文件内容Copy的方式移动
func Move(oldPath, newPath string) error {

	if err := os.Rename(oldPath, newPath); err != nil {
		le := err.(*os.LinkError)
		if le.Unwrap() == syscall.EEXIST {
			return err
		}

		if err = Copy(oldPath, newPath); err != nil {
			return err
		}

		return os.Remove(oldPath)
	}

	return nil
}

func Copy(srcPath, destPath string) error {
	fi, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("not support copy dir '%s'", srcPath)
	}

	r, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return err
}
