package main

import (
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ZZMarquis/gm/sm4"
	"github.com/ZZMarquis/gm/util"
	"github.com/guanzhi/GmSSL/go/gmssl"
	"github.com/stretchr/testify/assert"
)

func TestSM4Result(t *testing.T) {
	/* Generate random key and IV */
	keylen, _ := gmssl.GetCipherKeyLength("SMS4")
	key, _ := gmssl.GenerateRandom(keylen)
	ivlen, _ := gmssl.GetCipherIVLength("SMS4")
	iv, _ := gmssl.GenerateRandom(ivlen)
	bytez := []byte("needkane1234567890")

	/* SMS4-CBC Encrypt/Decrypt */
	ciphertext := SMS4Crypto(bytez, key, iv, true)

	plaintext := SMS4Crypto(ciphertext, key, iv, false)
	assert.Equal(t, bytez, plaintext)
	fmt.Printf("sms4(\"%s\") = %x\n", plaintext, ciphertext)

	c, err := sm4.NewCipher(key)
	assert.Nil(t, err)
	encrypter := cipher.NewCBCEncrypter(c, iv)
	result := make([]byte, 16)
	encrypter.CryptBlocks(result, util.PKCS5Padding(bytez, 16))
	fmt.Printf("encrypt result:%s\n", hex.EncodeToString(result))

}

/* SMS4-CBC Encrypt/Decrypt */
func SMS4Crypto(data, key, iv []byte, isEncrypt bool) []byte {
	cipCtx, _ := gmssl.NewCipherContext("SMS4", key, iv, isEncrypt)
	plaintext1, _ := cipCtx.Update(data)
	plaintext2, _ := cipCtx.Final()
	plaintext := make([]byte, 0, len(plaintext1)+len(plaintext2))
	plaintext = append(plaintext, plaintext1...)
	plaintext = append(plaintext, plaintext2...)
	return plaintext
}
