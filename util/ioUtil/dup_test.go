package ioUtil

import (
	"os"
	"testing"
)

func TestDupStderr(t *testing.T) {
	filename := "./testdata/stderr.dat"
	DupStderr(filename, func(err error) {
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("stderr redirect to " + filename)
		}
	})

	os.Stderr.WriteString("test dup stderr\n")
}
