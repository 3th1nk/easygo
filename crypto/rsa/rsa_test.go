package rsa

import (
	"github.com/3th1nk/easygo/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	publicKey, privateKey, _ = NewKeyPair()
	rsaEnc, _                = NewSecurity(publicKey, privateKey)
)

func Test_NewKeyPair(t *testing.T) {
	for _, length := range []int{256, 512, 1024} {
		for i := 0; i < 5; i++ {
			publicKey, privateKey, err := NewKeyPair(length)
			assert.NoError(t, err)
			assert.NotEqual(t, "", publicKey)
			assert.NotEqual(t, "", privateKey)
			util.Println("============================================= %v, %v, %v\n%v\n%v", length, len(publicKey), len(privateKey), publicKey, privateKey)
		}
	}
}

func Test_Encrypt(t *testing.T) {
	for i := 0; i < 10; i++ {
		cryptText, err := Encrypt("hello, world!", publicKey)
		assert.NoError(t, err)
		util.Println("cryptText: [%v]%v", len(cryptText), cryptText)

		plantText, err := Decrypt(cryptText, privateKey)
		assert.NoError(t, err)
		assert.Equal(t, "hello, world!", plantText)
	}
}

func Benchmark_Encrypt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rsaEnc.Encrypt("hello, world!")
	}
}

func Benchmark_Decrypt(b *testing.B) {
	cryptText, _ := rsaEnc.Encrypt("hello, world!")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		plantText, err := rsaEnc.Decrypt(cryptText)
		if err != nil {
			b.Errorf("PrivateDecrypt Error: %v", err)
		} else if plantText != "hello, world!" {
			b.Errorf("assert faild: %v", plantText)
		}
	}
}

func Benchmark_NewKeyPair_256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewKeyPair(256)
	}
}

func Benchmark_NewKeyPair_512(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewKeyPair(512)
	}
}

func Benchmark_NewKeyPair_1024(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewKeyPair(1024)
	}
}
