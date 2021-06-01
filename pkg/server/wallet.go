package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/clcng/bitcoin-wallet/pkg/errors"
	"github.com/clcng/bitcoin-wallet/pkg/util"
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

	//CS = ENT / 32
	//MS = (ENT + CS) / 11
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
