package runtimeUtil

import (
	"time"
)

// 重试
// 参数：
//    loops: 最大重试次数
//    interval: 重试间隔
// 返回值：
//    looped: 一共重试了多少次，从 0 开始。
func Retry(loops int, interval time.Duration, f func(loop int) (breaking bool, err error)) (looped int, err error) {
	var breaking bool
	for {
		breaking, err = f(looped)
		if err == nil || breaking {
			return
		} else if looped >= loops {
			return
		}
		time.Sleep(interval)
	}
}
