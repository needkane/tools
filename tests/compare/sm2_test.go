package main

import (
	"crypto/rand"
	"encoding/pem"
	"strings"
	"testing"

	zzsm2 "github.com/ZZMarquis/gm/sm2"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/guanzhi/GmSSL/go/gmssl"
	"github.com/stretchr/testify/assert"
	"github.com/tjfoc/gmsm/sm2"
)

func TestSM2Result(t *testing.T) {
	//secp256k1 btc
	priv0, err := btcec.NewPrivateKey(btcec.S256())
	assert.Nil(t, err)
	bytez := priv0.Serialize()
	priv1, _ := btcec.PrivKeyFromBytes(btcec.S256(), bytez)
	assert.Equal(t, priv0, priv1)

	//secp256k1 ethereum
	priv_eth0, err := crypto.GenerateKey()
	assert.Nil(t, err)
	bytez = crypto.FromECDSA(priv_eth0)
	priv_eth1, err := crypto.ToECDSA(bytez)
	assert.Nil(t, err)
	assert.Equal(t, priv_eth0, priv_eth1)

	//tjfoc
	priv_tj0, err := sm2.GenerateKey()
	assert.Nil(t, err)
	bytez, err = sm2.MarshalSm2UnecryptedPrivateKey(priv_tj0)
	assert.Nil(t, err)
	priv_tj1, err := sm2.ParsePKCS8UnecryptedPrivateKey(bytez)
	assert.Nil(t, err)
	assert.Equal(t, priv_tj0, priv_tj1)
	priv_tj3, err := sm2.ParseSm2PrivateKeyWithoutAsn(priv_tj1.D.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, priv_tj0, priv_tj3)

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: bytez,
	}
	pemBytes := pem.EncodeToMemory(block)
	priv_tj2, err := sm2.ReadPrivateKeyFromMem(pemBytes, nil)
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
	privBase64 := strings.Split(sm2skpem, "-----")[2]
	//blk, rest := pem.Decode([]byte(sm2skpem))
	var str string
	str += "-----BEGIN ENCRYPTED PRIVATE KEY-----"
	str += privBase64
	str += "-----END ENCRYPTED PRIVATE KEY-----"
	//	fmt.Println("=====", sm2skpem, "========", str)
	priv_ssl0, err := gmssl.NewPrivateKeyFromPEM(str, "pass")
	assert.Nil(t, err)
	assert.Equal(t, sm2sk, priv_ssl0)
	var signMethodName = "sm2sign"
	var msg = []byte("needkane")
	sign, err := sm2sk.Sign(signMethodName, msg, nil)
	assert.Nil(t, err)
	sm2pkpem, err := priv_ssl0.GetPublicKeyPEM()
	assert.Nil(t, err)
	sm2pk, err := gmssl.NewPublicKeyFromPEM(sm2pkpem)
	assert.Nil(t, err)
	//text, err := sm2pk.GetText()
	//assert.Nil(t, err)
	err = sm2pk.Verify(signMethodName, msg, sign, nil)
	assert.Nil(t, err)

	//ZZMarquis
	priv_zz0, _, err := zzsm2.GenerateKey(rand.Reader)
	assert.Nil(t, err)
	priv_zz1, err := zzsm2.RawBytesToPrivateKey(priv_zz0.D.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, priv_zz0, priv_zz1)
}
