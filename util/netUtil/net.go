package netUtil

import (
	"net"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var (
	localIp     []string
	localIpOnce sync.Once
)

/**
 * Local ip address list
 */
func LocalIp() []string {
	ensureLocalIp()
	return localIp
}

func FirstLocalIp() string {
	ensureLocalIp()
	if len(localIp) == 0 {
		return ""
	} else {
		return localIp[0]
	}
}

func ensureLocalIp() {
	localIpOnce.Do(func() {
		if runtime.GOOS == "linux" {
			faces, _ := net.Interfaces()
			for _, face := range faces {
				var name string
				if len(face.Name) > 3 {
					name = strings.ToLower(face.Name[:3])
				} else {
					name = strings.ToLower(face.Name)
				}
				if name == "ens" || name == "eth" {
					addr, _ := face.Addrs()
					for _, a := range addr {
						ipnet, ok := a.(*net.IPNet)
						if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
							localIp = append(localIp, ipnet.IP.String())
						}
					}
				}
			}
		} else {
			addr, _ := net.InterfaceAddrs()
			for _, a := range addr {
				ipnet, ok := a.(*net.IPNet)
				if ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsLinkLocalMulticast() && !ipnet.IP.IsLinkLocalUnicast() && ipnet.IP.To4() != nil {
					ip := ipnet.IP.String()
					if !strings.HasPrefix(ip, "192.168.") && !strings.HasSuffix(ip, ".1") {
						localIp = append(localIp, ip)
					}
				}
			}
		}
	})
}

// 获取 Http 客户端 IP
func GetClientIp(r *http.Request) string {
	// X-Forwarded-For
	if str := r.Header.Get("X-Forwarded-For"); str != "" {
		if pos := strings.Index(str, ","); pos != -1 {
			return str[:pos]
		} else {
			return str
		}
	}

	// X-Real-Ip
	if str := r.Header.Get("X-Real-Ip"); str != "" {
		return str
	}

	// RemoteAddr
	var ip string
	if pos := strings.Index(r.RemoteAddr, ":"); pos != -1 {
		ip = r.RemoteAddr[:pos]
	} else {
		ip = r.RemoteAddr
	}
	if ok, _ := regexp.MatchString("^(?:\\d+\\.){3}\\d+$", ip); ok {
		return ip
	} else {
		return ""
	}
}
