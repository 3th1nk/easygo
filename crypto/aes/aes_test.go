package aes

import (
	"encoding/base64"
	"github.com/3th1nk/easygo/util"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var (
	str = base64.StdEncoding.EncodeToString([]byte("abcdefg"))
	key = "01234567890123456789012345678901" // size == 32
)

func TestEncryptAndDecrypt(t *testing.T) {
	util.Println("str: %v", str)

	// 加密结果与 V1 相同，只是多了个前缀
	str2_0 := Encrypt(str, key)
	assert.True(t, strings.HasPrefix(str2_0, EncryptPrefix))
	util.Println("加密 1: %v", str2_0)

	// 已加密的字符串不会重复加密
	str2_1 := Encrypt(str2_0, key)
	assert.Equal(t, str2_0, str2_1)

	// 如果强制加密，还是可以加密的
	str2_2 := Encrypt(str2_1, key, true)
	assert.NotEqual(t, str2_1, str2_2)
	util.Println("加密 2: %v", str2_2)

	// 解密成功
	str3_1, err := Decrypt(str2_2, key)
	assert.NoError(t, err)
	assert.Equal(t, str2_1, str3_1)
	util.Println("解密 1: %v", str3_1)

	// 解密成功
	str3_2, err := Decrypt(str3_1, key)
	assert.NoError(t, err)
	assert.Equal(t, str, str3_2)
	util.Println("解密 2: %v", str3_2)

	// 解密成功，不指定强制解密
	str3_3, err := Decrypt(str3_2, key)
	assert.NoError(t, err)
	assert.Equal(t, str, str3_3)
	util.Println("解密 3: %v", str3_3)

	// 强制解密会失败，因为此时已经不是密文了
	str3_4, err := Decrypt(str3_2, key, true)
	assert.Error(t, err)
	assert.Equal(t, "", str3_4)
	util.Println("强制解密 4: %v", err)
}
