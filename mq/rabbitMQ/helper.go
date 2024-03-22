package rabbitMQ

import (
	"fmt"
	netUrl "net/url"
	"regexp"
	"strings"
)

var (
	urlRegex = regexp.MustCompile(`(amqp[s]?)://([^:]+)(?::(.*))?@(.*)`)
)

func AutoEncodeUrl(url string) string {
	scheme, username, password, host, err := ParseUrl(url, false)
	if err != nil {
		return url
	}

	username = netUrl.QueryEscape(username)
	password = netUrl.QueryEscape(password)
	if len(password) > 0 {
		return fmt.Sprintf("%s://%s:%s@%s", scheme, username, password, host)
	} else {
		return fmt.Sprintf("%s://%s@%s", scheme, username, host)
	}
}

// 解析形如 ‘amqp://user:password@host’ 的 rabbitMQ URL
//   keepRaw: 是否保持原始字符串、不进行 QueryUnescape 操作。默认为 false。
func ParseUrl(url string, keepRaw ...bool) (scheme, username, password, host string, err error) {
	matches := urlRegex.FindStringSubmatch(url)
	if len(matches) != 5 {
		return "", "", "", "", fmt.Errorf("unknown mq addr format")
	}

	scheme, username, password, host = matches[1], matches[2], matches[3], strings.TrimRight(matches[4], "/")
	if len(keepRaw) == 0 || !keepRaw[0] {
		if str, err := netUrl.QueryUnescape(username); err == nil {
			username = str
		}
		if str, err := netUrl.QueryUnescape(password); err == nil {
			password = str
		}
	}
	return
}
