package rest

import (
	"net/http"

	"github.com/bartekn/go-bip39"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

const (
	mnemonicEntropySize = 256
)

func QueryMnemonicRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed[:])
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, _ = w.Write([]byte(mnemonic))
}
