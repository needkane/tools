package main

import (
	"encoding/hex"
	"fmt"
	"testing"

	zzsm3 "github.com/ZZMarquis/gm/sm3"
	"github.com/guanzhi/GmSSL/go/gmssl"
	ssl_sm3 "github.com/guanzhi/GmSSL/go/gmssl/sm3"
	"github.com/stretchr/testify/assert"
	"github.com/tjfoc/gmsm/sm3"
)

func TestSM3Result(t *testing.T) {
	bytez := []byte("abc")

	hasher := sm3.New()
	hasher.Write(bytez)
	r0 := hasher.Sum(nil)
	fmt.Println(hex.EncodeToString(r0))
	sm3ctx, _ := gmssl.NewDigestContext("SM3")
	sm3ctx.Update(bytez)
	r1, _ := sm3ctx.Final()

	sslHasher := ssl_sm3.New()
	sslHasher.Write(bytez)
	r2 := sslHasher.Sum(nil)

	d := zzsm3.New()
	d.Write(bytez)
	hash := d.Sum(nil)
	hashHex := hex.EncodeToString(hash[:])
	fmt.Println(string(bytez), hashHex)
	assert.Equal(t, r0, r1)
	assert.Equal(t, r0, r2)

}
