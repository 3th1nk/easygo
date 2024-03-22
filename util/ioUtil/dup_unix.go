//go:build (linux && !arm && !arm64) || darwin || freebsd || openbsd || netbsd
// +build linux,!arm,!arm64 darwin freebsd openbsd netbsd

package ioUtil

import (
	"syscall"
)

func dup(from, to int) error {
	return syscall.Dup2(from, to)
}
