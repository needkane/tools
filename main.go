package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type HttpResult struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

func ErrorResponse(w http.ResponseWriter, errCode int, err error) {

	var hr HttpResult
	hr.Error = err.Error()
	bytez, err := json.Marshal(hr)
	w.WriteHeader(errCode)
	w.Write(bytez)
}

func ResultResponse(w http.ResponseWriter, result []byte) {
	fmt.Println("resp:   ", string(result))
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

type HashBody struct {
	Method  string `json:"method"`
	Content string `json:"content"`
}

type CodecBody struct {
	Method  string `json:"method"`
	Content string `json:"content"`
}

type AsymmetricBody struct {
	Method    string `json:"method"`
	Operation string `json:"operation"`
	Content   string `json:"content"`
}

func CheckRequest(r *http.Request) (reqBytes []byte, err error) {

	if r.Method != "POST" {
		err = fmt.Errorf("invald http method:%s", r.Method)
		return
	}
	reqBytes, err = ioutil.ReadAll(r.Body)
	fmt.Println("req:   ", string(reqBytes))
	return
}
func CryptoCodecHandler(w http.ResponseWriter, r *http.Request) {

	reqBytes, err := CheckRequest(r)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	var cb CodecBody
	err = json.Unmarshal(reqBytes, &cb)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	var hr HttpResult
	switch cb.Method {
	case "base64_encode":
		uEnc := base64.URLEncoding.EncodeToString([]byte(cb.Content))
		hr.Result = uEnc
	case "base64_decode":
		uDec, err := base64.URLEncoding.DecodeString(cb.Content)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		hr.Result = string(uDec)
	default:
		err = fmt.Errorf("invalid crypto method: %s", cb.Method)
		ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	bytez, _ := json.Marshal(hr)
	ResultResponse(w, bytez)
}

func CryptoAsymmetricHandler(w http.ResponseWriter, r *http.Request) {

	reqBytes, err := CheckRequest(r)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	var ab AsymmetricBody
	err = json.Unmarshal(reqBytes, &ab)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	var hr HttpResult
	switch ab.Method {
	case "secp256k1":
		if ab.Operation == "generate" {
			privkey, err := crypto.GenerateKey()
			if err != nil {
				ErrorResponse(w, http.StatusInternalServerError, err)
				return
			}
			var result struct {
				Privkey string `json:"privkey"`
				Pubkey  string `json:"pubkey"`
				Address string `json:"address"`
			}
			address := crypto.PubkeyToAddress(privkey.PublicKey)
			result.Privkey = common.Bytes2Hex(crypto.FromECDSA(privkey))
			result.Pubkey = common.Bytes2Hex(crypto.FromECDSAPub(&privkey.PublicKey))
			result.Address = address.Hex()
			hr.Result = result
		} else if ab.Operation == "get_address" {
			privBytes := common.Hex2Bytes(ab.Content)
			privkey, err := crypto.ToECDSA(privBytes)
			if err != nil {
				ErrorResponse(w, http.StatusInternalServerError, err)
				return
			}
			address := crypto.PubkeyToAddress(privkey.PublicKey)
			var result struct {
				Privkey string `json:"privkey"`
				Address string `json:"address"`
			}
			result.Privkey = ab.Content
			result.Address = address.Hex()
			hr.Result = result

		}
	default:
		err = fmt.Errorf("invalid crypto method: %s", ab.Method)
		ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	bytez, _ := json.Marshal(hr)
	ResultResponse(w, bytez)
}
func CryptoHashHandler(w http.ResponseWriter, r *http.Request) {

	reqBytes, err := CheckRequest(r)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	var hb HashBody
	err = json.Unmarshal(reqBytes, &hb)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	var hr HttpResult
	switch hb.Method {
	case "md5":
		m := md5.New()
		m.Write([]byte(hb.Content))
		hr.Result = hex.EncodeToString(m.Sum(nil))
	default:
		err = fmt.Errorf("invalid crypto method: %s", hb.Method)
		ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	bytez, _ := json.Marshal(hr)
	ResultResponse(w, bytez)
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/crypto/hash", CryptoHashHandler)
	mux.HandleFunc("/crypto/codec", CryptoCodecHandler)
	mux.HandleFunc("/crypto/asymmetric", CryptoAsymmetricHandler)
	http.ListenAndServe("0.0.0.0:12580", mux)
}
