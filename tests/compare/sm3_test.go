package main

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/guanzhi/GmSSL/go/gmssl"
	ssl_sm3 "github.com/guanzhi/GmSSL/go/gmssl/sm3"
	"github.com/stretchr/testify/assert"
	"github.com/tjfoc/gmsm/sm3"
)

func TestSM3Result(t *testing.T) {
	bytez := []byte("needkane")

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

	assert.Equal(t, r0, r1)
	assert.Equal(t, r0, r2)

}
