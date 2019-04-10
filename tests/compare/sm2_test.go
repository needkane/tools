package main

import (
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/guanzhi/GmSSL/go/gmssl"
	"github.com/stretchr/testify/assert"
	"github.com/tjfoc/gmsm/sm2"
)

func TestSM2Result(t *testing.T) {
	//secp256k1
	priv0, err := btcec.NewPrivateKey(btcec.S256())
	assert.Nil(t, err)
	bytez := priv0.Serialize()
	priv1, _ := btcec.PrivKeyFromBytes(btcec.S256(), bytez)
	assert.Equal(t, priv0, priv1)

	//tjfoc
	priv_tj0, err := sm2.GenerateKey()
	assert.Nil(t, err)
	bytez, err = sm2.MarshalSm2UnecryptedPrivateKey(priv_tj0)
	assert.Nil(t, err)
	priv_tj1, err := sm2.ParsePKCS8UnecryptedPrivateKey(bytez)
	assert.Nil(t, err)
	assert.Equal(t, priv_tj0, priv_tj1)

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: bytez,
	}
	pemBytes := pem.EncodeToMemory(block)
	priv_tj2, err := sm2.ReadPrivateKeyFromMem(pemBytes, nil)
	fmt.Println(len(pemBytes))
	assert.Equal(t, priv_tj0, priv_tj2)

	//gmssl
	sm2keygenargs := map[string]string{
		"ec_paramgen_curve": "sm2p256v1",
		"ec_param_enc":      "named_curve",
	}
	sm2sk, err := gmssl.GeneratePrivateKey("EC", sm2keygenargs, nil)
	assert.Nil(t, err)
	sm2skpem, err := sm2sk.GetPEM("SMS4", "pass")
	assert.Nil(t, err)
	priv_ssl0, err := gmssl.NewPrivateKeyFromPEM(sm2skpem, "pass")
	assert.Nil(t, err)
	assert.Equal(t, sm2sk, priv_ssl0)
	var signMethodName = "sm2sign"
	var msg = []byte("needkane2")
	sign, err := sm2sk.Sign(signMethodName, msg, nil)
	assert.Nil(t, err)
	sm2pkpem, err := priv_ssl0.GetPublicKeyPEM()
	assert.Nil(t, err)
	sm2pk, err := gmssl.NewPublicKeyFromPEM(sm2pkpem)
	assert.Nil(t, err)
	text, err := sm2pk.GetText()
	assert.Nil(t, err)
	fmt.Println(len(sm2pkpem), len(text), "=====", len(sm2skpem), "========")

	err = sm2pk.Verify(signMethodName, msg, sign, nil)
	assert.Nil(t, err)
}
