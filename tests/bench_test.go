/*
Copyright Suzhou Tongji Fintech Research Institute 2017 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tests

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/tjfoc/gmsm/sm2"
)

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

func BenchmarkSM2(t *testing.B) {
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
			fmt.Printf("VerifySignature error:%v\n", err)
		}
	}
}
