package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"io/ioutil"

	"github.com/clcng/bitcoin-wallet/pkg/errors"
	"github.com/clcng/bitcoin-wallet/pkg/util"
	"github.com/clcng/bitcoin-wallet/pkg/keymanager"
	"github.com/tyler-smith/go-bip39"
)

func (s *Server) handleGenerateBIP39Mnemonic(w http.ResponseWriter, r *http.Request) {
	possibleNoOfWords := []string{"12", "15", "18", "21", "24"}
	noOfWords := r.URL.Query().Get("words")
	if !util.ContainString(possibleNoOfWords, noOfWords) {
		renderErr(w, errors.Coded(1001, errors.ErrorMap[1001]+", possible value of words: 12,15,18,21,24"))
		return
	}

	noOfWordsFloat, err := strconv.ParseFloat(noOfWords, 64)
	if err != nil {
		renderErr(w, err)
		return
	}

	noOfBits := noOfWordsFloat * 11.0 / (1.0 + 1.0/32.0)
	entropy, err := bip39.NewEntropy(int(noOfBits))
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		renderErr(w, err)
		return
	}

	result := map[string]interface{}{
		"mnemonic": mnemonic,
	}
	renderJSON(w, result)
}

func (s *Server) handleGenerateHDAddress(w http.ResponseWriter, r *http.Request) {
	in := &GenerateHDAddressInput{}
	_, err := decode(r, in)
	if err != nil {
		renderErr(w, errors.Coded(1001, errors.ErrorMap[1001]+", fail to decode body"))
		return
	}

	km, err := keymanager.NewKeyManager(in.Mnemonic)
	if err != nil {
		renderErr(w, errors.Coded(2001, errors.ErrorMap[2001]))
		return
	}

	key, err := km.GetKey(in.Path)
	if err != nil {
		renderErr(w, errors.Coded(2002, errors.ErrorMap[2002]))
		return
	}

	wif, address, _, _, err := key.Encode()
	if err != nil {
		renderErr(w, errors.Coded(2001, errors.ErrorMap[2001]+", fail to init key manager"))
		return
	}

	result := map[string]interface{}{
		"address": address,
		"wif": wif,
	}
	renderJSON(w, result)
}

func renderErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	ba, err2 := json.Marshal(err)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(ba)
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	var ba []byte
	if bv, ok := v.([]byte); ok {
		ba = bv
	} else {
		ba, _ = json.Marshal(v)
	}
	if ba != nil {
		w.Write(ba)
	}
}

func decode(r *http.Request, v interface{}) ([]byte, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, v)
	return body, err
}