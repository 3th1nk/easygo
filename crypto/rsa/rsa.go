package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"github.com/wenzhenxi/gorsa"
)

// NewKeyPair 生成密钥对
func NewKeyPair(length ...int) (publicKey, privateKey string, err error) {
	var theLength int
	if len(length) != 0 {
		theLength = length[0]
	}
	if theLength == 0 {
		theLength = 1024
	} else if theLength < 256 {
		theLength = 256
	}

	buffer := bytes.NewBuffer(make([]byte, 0, 4096))

	rsaKey, err := rsa.GenerateKey(rand.Reader, theLength)
	if err != nil {
		return
	}

	// 生成私钥
	buf, err := asn1.Marshal(struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{
		PrivateKeyAlgorithm: []asn1.ObjectIdentifier{
			{1, 2, 840, 113549, 1, 1, 1},
		},
		PrivateKey: x509.MarshalPKCS1PrivateKey(rsaKey),
	})
	if err == nil {
		err = pem.Encode(buffer, &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: buf,
		})
	}
	if err != nil {
		return
	}
	privateKey = buffer.String()

	// 生成公钥
	buffer.Reset()
	buf, err = x509.MarshalPKIXPublicKey(&rsaKey.PublicKey)
	if err == nil {
		err = pem.Encode(buffer, &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: buf,
		})
	}
	if err != nil {
		return
	}
	publicKey = buffer.String()

	return
}

// Encrypt 公钥加密
func Encrypt(plainText, publicKey string) (cipherText string, err error) {
	enc, err := NewSecurity(publicKey, "")
	if err != nil {
		return "", err
	}
	return enc.Encrypt(plainText)
}

// Decrypt 私钥解密
func Decrypt(cipherText, privateKey string) (plainText string, err error) {
	enc, err := NewSecurity("", privateKey)
	if err != nil {
		return "", err
	}
	return enc.Decrypt(cipherText)
}

// NewSecurity 根据指定的 RSA 密钥对创建 Security。
//   公钥或者私钥可以设置为空，但这种情况下无法使用相应的功能。
func NewSecurity(publicKey, privateKey string) (*Security, error) {
	enc := &Security{rsa: &gorsa.RSASecurity{}}
	if err := enc.SetKey(publicKey, privateKey); err != nil {
		return nil, err
	}
	return enc, nil
}

type Security struct {
	pubStr string
	priStr string
	rsa    *gorsa.RSASecurity
}

func (this *Security) PublicKey() string {
	return this.pubStr
}

func (this *Security) PrivateKey() string {
	return this.priStr
}

func (this *Security) SetKey(publicKey, privateKey string) error {
	if publicKey != "" {
		if err := this.rsa.SetPublicKey(publicKey); err != nil {
			return err
		}
		this.pubStr = publicKey
	}
	if privateKey != "" {
		if err := this.rsa.SetPrivateKey(privateKey); err != nil {
			return err
		}
		this.priStr = privateKey
	}
	return nil
}

// Encrypt 公钥加密
func (this *Security) Encrypt(plainText string) (string, error) {
	buf, err := this.rsa.PubKeyENCTYPT([]byte(plainText))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

// Decrypt 私钥解密
func (this *Security) Decrypt(cipherText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	buf, err := this.rsa.PriKeyDECRYPT(data)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
