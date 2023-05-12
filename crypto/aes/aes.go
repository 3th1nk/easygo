package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"strings"
)

var (
	EncryptPrefix  = "encrypt:"
	ErrInvalidKey  = errors.New("encrypt: invalid key")
	ErrInvalidData = errors.New("encrypt: invalid data")
)

// Encrypt 字符串加密
//  加密后，会在密文前面追加一个字符串用来标记这是一个加密字符串。
//  加密前，可以根据这个标记判断这是否已经是一个密文，从而决定还要不要继续加密。
// 	鉴于有可能会存在部分场景确实需要对密文进行二次加密，故增加 force 参数用来指定强制加密。
func Encrypt(plainText, key string, force ...bool) (string, error) {
	if n := len(key); n != 16 && n != 24 && n != 32 {
		// aes.NewCipher 产生错误只有一个原因：密钥长度不是 16|24|32。
		// 密钥应当是在系统设计阶段就协商确定的，必须符合该要求。故出现错误说明存在严重的设计缺陷，应当立即以最严重的方式暴露、以便尽快修改。
		return plainText, ErrInvalidKey
	}

	if strings.HasPrefix(plainText, EncryptPrefix) && (len(force) == 0 || !force[0]) {
		// 已经加密过、且不强制加密，直接返回。
		return plainText, nil
	}

	keyBytes := []byte(key)
	block, _ := aes.NewCipher(keyBytes)

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, keyBytes[:blockSize])

	data := []byte(plainText)
	encryptBytes := paddingPKCS7(data, blockSize)
	data = make([]byte, len(encryptBytes))
	blockMode.CryptBlocks(data, encryptBytes)
	return EncryptPrefix + base64.StdEncoding.EncodeToString(data), nil
}

// Decrypt 字符串解密
//  与加密方法对应，先判断开头是否存在标记，如果不存在则表示这不是一个加密字符串。
//  解密前，可以根据这个标记判断这是否已经是一个密文，从而决定还要不要继续解密。
//  兼容不带加密标识的密文，故增加 force 参数用来指定强制解密。
func Decrypt(cipherText string, key string, force ...bool) (result string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = ErrInvalidData
		}
	}()

	if n := len(key); n != 16 && n != 24 && n != 32 {
		// aes.NewCipher 产生错误只有一个原因：密钥长度不是 16|24|32。
		// 密钥应当是在系统设计阶段就协商确定的，必须符合该要求。故出现错误说明存在严重的设计缺陷，应当立即以最严重的方式暴露、以便尽快修改。
		return "", ErrInvalidKey
	}

	if strings.HasPrefix(cipherText, EncryptPrefix) {
		cipherText = cipherText[len(EncryptPrefix):]
	} else if len(force) == 0 || !force[0] {
		// 没有加密标识且非强制解密，直接返回。
		return cipherText, nil
	}

	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", ErrInvalidData
	}

	keyBytes := []byte(key)
	block, _ := aes.NewCipher(keyBytes)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyBytes[:blockSize])

	decryptStr := make([]byte, len(data))
	// 当密文数据异常时会panic
	blockMode.CryptBlocks(decryptStr, data)
	decryptStr = unPaddingPKCS7(decryptStr)
	return string(decryptStr), nil
}
