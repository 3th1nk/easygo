//go:build linux && (arm || arm64)
// +build linux
// +build arm arm64

package ioUtil

func dup(from, to int) error {
	return syscall.Dup3(from, to, 0)
}
