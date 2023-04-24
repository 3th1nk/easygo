package runtimeUtil

import (
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	isTesting     = -1
	isTestingOnce sync.Once
)

// IsTesting 判断是否是在执行单元测试
func IsTesting() bool {
	isTestingOnce.Do(func() {
		isTesting = 0
		for _, str := range os.Args[1:] {
			if -1 != strings.Index(strings.ToLower(str), "-test.v") ||
				-1 != strings.Index(strings.ToLower(str), "-test.run") {
				isTesting = 1
				break
			}
		}
	})
	return isTesting == 1
}

// ChdirWithTesting 判断是否在执行单元测试，尝试切换目录正则匹配的目录
func ChdirWithTesting(reg string) {
	if IsTesting() {
		if dir, _ := os.Getwd(); dir != "" {
			if loc := regexp.MustCompile(reg).FindStringIndex(dir); len(loc) > 1 {
				rootDir := dir[:loc[1]]
				os.Chdir(rootDir)
			}
		}
	}
}
