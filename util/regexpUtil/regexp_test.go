package regexpUtil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPattern_IPv4(t *testing.T) {
	for n1 := 0; n1 <= 300; n1 += 30 {
		for n2 := 0; n2 <= 300; n2 += 30 {
			for n3 := 0; n3 <= 300; n3 += 30 {
				for n4 := 0; n4 <= 300; n4 += 30 {
					ip := fmt.Sprintf("%v.%v.%v.%v", n1, n2, n3, n4)
					isValid := n1 <= 255 && n2 <= 255 && n3 <= 255 && n4 <= 255
					assert.Equal(t, isValid, IsIPv4(ip), ip)

					arr1 := FindIPv4(ip)
					arr2 := FindIPv4("aaa" + ip + "bbb")
					if isValid {
						if assert.Equal(t, 1, len(arr1)) {
							assert.Equal(t, ip, arr1[0])
						}
						if assert.Equal(t, 1, len(arr2)) {
							assert.Equal(t, ip, arr2[0])
						}
					} else {
						assert.Equal(t, 0, len(arr1))
						assert.Equal(t, 0, len(arr2))
					}
				}
			}
		}
	}
}

func TestIsInt(t *testing.T) {
	arr := strings.Split("1,123,0,-0,-1,-2,-3", ",")
	for _, s := range arr {
		if !IsInt(s) {
			t.Errorf("assert faild: %v expect true", s)
		}
	}

	arr = strings.Split("a,12a,a12,12a12,a12a", ",")
	for _, s := range arr {
		if IsInt(s) {
			t.Errorf("assert faild: %v expect not true", s)
		}
	}
}

func TestFindInt(t *testing.T) {
	for _, arr := range [][]interface{}{
		{"abc123", "123"},
		{"abc", ""},
		{"-1day", "-1"},
		{"-1 day", "-1"},
		{"1day", "1"},
		{"1 day", "1"},
		{"1 day, 2 sec", "1,2"},
		{"123+456", "123,456"},
	} {
		str, expect := arr[0].(string), arr[1].(string)
		actual := strings.Join(FindInt(str), ",")
		assert.Equal(t, expect, actual)
	}
}

func TestIsFloat(t *testing.T) {
	assert.False(t, IsFloat("300M"))

	arr := strings.Split("1,123,0,-0,-1,-2,-3,.1,.123,123.,123.0,0.123", ",")
	for _, s := range arr {
		if !IsFloat(s) {
			t.Errorf("assert faild: %v expect true", s)
		}
	}

	arr = strings.Split("a,12a,a12,12a12,a12a,..,12..12", ",")
	for _, s := range arr {
		if IsInt(s) {
			t.Errorf("assert faild: %v expect not true", s)
		}
	}
}

func TestFindIPv4(t *testing.T) {
	for _, arr := range [][]interface{}{
		// {"abc123", ""},
		{"192.168.1.0/24, 192.168.1.0/28", "192.168.1.0,192.168.1.0"},
		{`
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
2: ens33: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 00:0c:29:94:01:e5 brd ff:ff:ff:ff:ff:ff
    inet 172.16.20.184/24 brd 172.16.20.255 scope global noprefixroute dynamic ens33
       valid_lft 18637sec preferred_lft 18637sec
`, "127.0.0.1,172.16.20.184,172.16.20.255"},
	} {
		str, expect := arr[0].(string), arr[1].(string)
		actual := strings.Join(FindIPv4(str), ",")
		assert.Equal(t, expect, actual)
	}
}

func TestFindIPv4Addr(t *testing.T) {
	arr := FindIPv4Addr("http://192.168.1.213:20002/api/abc, 192.168.1.213:20002, 127.0.0.1")
	assert.Equal(t, "192.168.1.213:20002", arr[0])
	assert.Equal(t, "192.168.1.213:20002", arr[1])
	assert.Equal(t, "127.0.0.1", arr[2])
	assert.Equal(t, "192.168.1.213:20002", FindFirstIPv4Addr("amqp://user:Pwd@123@192.168.1.213:20002"))
}
