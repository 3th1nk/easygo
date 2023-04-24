package arrUtil

import (
	"github.com/3th1nk/easygo/util/mathUtil"
)

func DiffInt(src, dest []int) (matches, added, removed []int) {
	maxLen := mathUtil.MaxInt(len(src), len(dest))
	matches, added, removed = make([]int, 0, maxLen), make([]int, 0, maxLen), make([]int, 0, maxLen)
	for _, str := range src {
		if -1 != IndexOfInt(dest, str) {
			matches = append(matches, str)
		} else {
			removed = append(removed, str)
		}
	}
	for _, str := range dest {
		if -1 == IndexOfInt(src, str) {
			added = append(added, str)
		}
	}
	return
}

func DiffInt64(src, dest []int64) (matches, added, removed []int64) {
	maxLen := mathUtil.MaxInt(len(src), len(dest))
	matches, added, removed = make([]int64, 0, maxLen), make([]int64, 0, maxLen), make([]int64, 0, maxLen)
	for _, str := range src {
		if -1 != IndexOfInt64(dest, str) {
			matches = append(matches, str)
		} else {
			removed = append(removed, str)
		}
	}
	for _, str := range dest {
		if -1 == IndexOfInt64(src, str) {
			added = append(added, str)
		}
	}
	return
}

func DiffString(src, dest []string, ignoreCase ...bool) (matches, added, removed []string) {
	maxLen := mathUtil.MaxInt(len(src), len(dest))
	matches, added, removed = make([]string, 0, maxLen), make([]string, 0, maxLen), make([]string, 0, maxLen)
	for _, str := range src {
		if -1 != IndexOfString(dest, str, ignoreCase...) {
			matches = append(matches, str)
		} else {
			removed = append(removed, str)
		}
	}
	for _, str := range dest {
		if -1 == IndexOfString(src, str, ignoreCase...) {
			added = append(added, str)
		}
	}
	return
}
