package gmssltest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/guanzhi/GmSSL/go/gmssl"
	"github.com/tjfoc/gmsm/sm2"
)

//use global variable for fix benchmark error
//(elliptic curve routines:eckey_param2type:missing parameters:crypto/ec/ec_ameth.c:84)
var sm2sk *gmssl.PrivateKey

func init() {
	newSm2sk()
}

func BenchmarkTjSM2(t *testing.B) {
	t.ReportAllocs()
	msg := []byte("abcdefghijklmnopqrstuvwxyz_abcde")
	priv, err := sm2.GenerateKey() // 生成密钥对
	if err != nil {
		log.Fatal(err)
	}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		sign, err := priv.Sign(rand.Reader, msg, nil) // 签名
		if err != nil {
			log.Fatal(err)
		}
		ok := priv.PublicKey.Verify(msg, sign) // 密钥验证
		if ok != true {
			fmt.Printf("Verify error\n")
		}
	}
}

func BenchmarkSecp256(t *testing.B) {

	t.ReportAllocs()
	msg := []byte("abcdefghijklmnopqrstuvwxyz_abcde")
	pubBytes, privBytes := generateKeyPair()
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		sign, err := secp256k1.Sign(msg, privBytes)
		if err != nil {
			log.Fatal(err)
		}
		err = verifySecp256(pubBytes, msg, sign) // 密钥验证
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkGmsslSM2(t *testing.B) {
	t.ReportAllocs()
	sm2pkpem, err := sm2sk.GetPublicKeyPEM()
	if err != nil {
		log.Fatal(err)
	}
	sm2pk, err := gmssl.NewPublicKeyFromPEM(sm2pkpem)
	if err != nil {
		log.Fatal(err)
	}
	dgst, err := getDgst(sm2pk)
	signMethodName := "sm2sign"
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		sign, err := sm2sk.Sign(signMethodName, dgst, nil) // 签名
		if err != nil {
			log.Fatal(err)
		}
		err = sm2pk.Verify(signMethodName, dgst, sign, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func verifySecp256(pubkey, msg, sign []byte) error {
	pubBytez, err := secp256k1.RecoverPubkey(msg, sign)
	if err != nil {
		return err
	}
	if !sliceEqual(pubBytez, pubkey) {
		return errors.New("Verify failed")
	}
	return nil
}

func generateKeyPair() (pubkey, privkey []byte) {

	factor := secp256k1.S256()
	key, err := ecdsa.GenerateKey(factor, rand.Reader)
	if err != nil {
		panic(err)
	}
	pubkey = elliptic.Marshal(factor, key.X, key.Y)

	privkey = make([]byte, 32)
	blob := key.D.Bytes()
	copy(privkey[32-len(blob):], blob)

	return pubkey, privkey
}

func sliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func newSm2sk() {
	sm2keygenargs := map[string]string{
		"ec_paramgen_curve": "sm2p256v1",
		"ec_param_enc":      "named_curve",
	}
	var err error
	sm2sk, err = gmssl.GeneratePrivateKey("EC", sm2keygenargs, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getDgst(sm2pk *gmssl.PublicKey) ([]byte, error) {
	sm3ctx, _ := gmssl.NewDigestContext("SM3")
	sm2zid, _ := sm2pk.ComputeSM2IDDigest("1234567812345678")
	sm3ctx.Reset()
	sm3ctx.Update(sm2zid)
	msg := []byte("abcdefghijklmnopqrstuvwxyz_abcde")
	sm3ctx.Update(msg)
	return sm3ctx.Final()
}
