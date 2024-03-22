package ioUtil

import (
	"github.com/3th1nk/easygo/util"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testGetDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	baseDir := filepath.Join(dir, "testDir")
	os.RemoveAll(baseDir)

	return baseDir
}

func testCreateDirFile() string {

	dir := testGetDir()
	util.PrintArgsLn("testDir: \t", dir)

	err := MakeDir(dir)
	if err != nil {
		return ""
	}

	defer func() {
		if err != nil {
			util.PrintArgsLn("test dir data err", err.Error())
			os.RemoveAll(dir)
		}
	}()

	for _, path := range []string{
		"emptyDir",
		"subDir",
	} {
		arr := strings.Split(path, "/")
		if err = MakeDir(filepath.Join(append([]string{dir}, arr...)...)); err != nil {
			return ""
		}
	}

	for _, path := range []string{
		"a.txt",
		"subDir/b.txt",
	} {
		arr := strings.Split(path, "/")
		f, e := OpenFile(filepath.Join(append([]string{dir}, arr...)...))
		if e != nil {
			err = e
			return ""
		}
		f.WriteString("\n")
		f.WriteString(path)
		f.Close()
	}

	return dir
}

func TestGZipTarToFile(t *testing.T) {
	baseDir := testCreateDirFile()
	if baseDir == "" {
		return
	}

	targetDir := baseDir + ".tar.gz"
	os.Remove(targetDir)
	if err := GZipTarToFile(baseDir, targetDir, true); err != nil {
		t.Fatal(err)
	}
}

func TestUnGZipTarFileToDir(t *testing.T) {
	baseDir := testGetDir()
	if baseDir == "" {
		return
	}

	gzipFilename := baseDir + ".tar.gz"
	if !IsFileExists(gzipFilename) {
		return
	}

	err := UnGZipTarFileToDir(gzipFilename, baseDir, true)
	if err != nil {
		t.Fatal(err)
	}

}
