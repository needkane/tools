package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpResult struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

func ErrorResponse(w http.ResponseWriter, errCode int, err error) {

	var hr HttpResult
	hr.Error = err
	bytez, _ := json.Marshal(hr)
	w.WriteHeader(errCode)
	w.Write(bytez)
}

func ResultResponse(w http.ResponseWriter, result []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

type InputBody struct {
	Type  string `json:"type"`
	Input string `json:"input"`
}

func CryptoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		err := fmt.Errorf("invald http method:%s", r.Method)
		ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	var ib InputBody
	err = json.Unmarshal(body, &ib)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	var hr HttpResult
	switch ib.Type {
	case "md5":
		w := md5.New()
		w.Write([]byte(ib.Input))
		hr.Result = string(hex.EncodeToString(w.Sum(nil)))
	default:
		hr.Error = errors.New("invalid crypto type")
		ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	bytez, _ := json.Marshal(hr)
	ResultResponse(w, bytez)
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/category/crypto", CryptoHandler)
	http.ListenAndServe("0.0.0.0:12580", mux)
}
