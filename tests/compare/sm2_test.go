package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"

	zzsm2 "github.com/ZZMarquis/gm/sm2"
	zzsm3 "github.com/ZZMarquis/gm/sm3"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/guanzhi/GmSSL/go/gmssl"
	"github.com/stretchr/testify/assert"
	"github.com/tjfoc/gmsm/sm2"
)

func TestSM2Result(t *testing.T) {
	var msg = []byte("needkane")

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
	sign, err := priv_tj2.Sign(rand.Reader, msg, nil) // 签名
	assert.Nil(t, err)
	ok := priv_tj2.PublicKey.Verify(msg, sign) // 密钥验证
	assert.Equal(t, ok, true)

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
	sign, err = sm2sk.Sign(signMethodName, msg, nil)
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
	priv_zz0, pub_zz0, err := zzsm2.GenerateKey(rand.Reader)
	assert.Nil(t, err)
	priv_zz1, err := zzsm2.RawBytesToPrivateKey(priv_zz0.D.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, priv_zz0, priv_zz1)
	sign, err = zzsm2.Sign(priv_zz1, nil, msg)
	assert.Nil(t, err)
	ok = zzsm2.Verify(pub_zz0, nil, msg, sign)
	assert.Equal(t, ok, true)
	fmt.Println(pub_zz0.GetRawBytes(), len(pub_zz0.GetRawBytes()))

	privStr := "ae08dc67186f140235a36a06e55dc2ccabbc5365525825c382aa36e055de84cd"
	bytez, err = hex.DecodeString(privStr)
	assert.Nil(t, err)
	priv_zz2, err := zzsm2.RawBytesToPrivateKey(bytez)
	assert.Nil(t, err)

	pub_zz2 := zzsm2.CaculatePubKey(priv_zz2)
	d := zzsm3.New()
	d.Write(pub_zz2.GetRawBytes()) //GetRawBytes == raw[1:]
	hash := d.Sum(nil)
	assert.Equal(t, "4a21554fcca7fdd8c183ecaab3a797c7dfce6de5", hex.EncodeToString(hash[12:]))

}
