//go:build windows
// +build windows

package ioUtil

func dup(from, to int) error {
	_, err := syscall.DuplicateHandle(syscall.GetCurrentProcess(), syscall.Handle(from), syscall.GetCurrentProcess(), &syscall.Handle(to), 0, true, syscall.DUPLICATE_SAME_ACCESS)
	return err
}
