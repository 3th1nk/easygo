// 文件监控器，用来实现 “监控某个文件的变化，如果文件被修改则触发回调通知”
package fileWatcher

import (
	"crypto/md5"
	"fmt"
	"github.com/3th1nk/easygo/util/timeUtil"
	"github.com/modern-go/reflect2"
	"os"
	"time"
)

func LoadAndWatch(path string, d time.Duration, onChange func(data []byte, err error)) (*timeUtil.Ticker, error) {
	if onChange == nil {
		return nil, fmt.Errorf("missing param: 'onChange'")
	}

	lastMod, lastToken, data, err := loadFile(path, 0, "")
	if data != nil || err != nil {
		onChange(data, err)
	}
	if err != nil {
		return nil, err
	}

	if d > 0 {
		ticker := timeUtil.NewTicker(d, d, func(t time.Time) {
			lastMod, lastToken, data, err = loadFile(path, lastMod, lastToken)
			if data != nil || err != nil {
				onChange(data, err)
			}
		})
		return ticker, nil
	}

	return nil, nil
}

func loadFile(path string, lastMod int64, lastToken string) (int64, string, []byte, error) {
	// 检查文件有没有修改过
	osFile, err := os.Stat(path)
	if !reflect2.IsNil(err) {
		return lastMod, lastToken, nil, fmt.Errorf("无法访问文件 %v", path)
	}

	// 检查文件最后修改时间有没有发生变化
	if osFile.ModTime().UnixNano() <= lastMod {
		return lastMod, lastToken, nil, nil
	} else {
		lastMod = osFile.ModTime().UnixNano()
	}

	// 检查文件内容有没有发生变化
	b, err := os.ReadFile(path)
	if !reflect2.IsNil(err) {
		return lastMod, lastToken, nil, err
	}
	tokenStr := fmt.Sprintf("%x", md5.Sum(b))
	if lastToken == tokenStr {
		return lastMod, lastToken, nil, nil
	}

	return lastMod, tokenStr, b, nil
}
