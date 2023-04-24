package ioUtil

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(filename), perm); err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, perm)
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
