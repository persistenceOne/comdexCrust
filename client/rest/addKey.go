package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bartekn/go-bip39"
	"github.com/cosmos/cosmos-sdk/client/keys"
	cryptoKeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

const (
	mnemonicEntropySize = 256
)

type NewKeyBody struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Seed     string `json:"seed"`
}

func AddNewKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m NewKeyBody
	var account uint32
	var index uint32

	kb, err := keys.NewKeyBaseFromHomeFlag()
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &m)

	seed := m.Seed

	if len(seed) == 0 {
		// read entropy seed straight from crypto.Rand and convert to mnemonic
		entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		seed, err = bip39.NewMnemonic(entropySeed[:])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	infos, err := kb.List()
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, i := range infos {
		if i.GetName() == m.Name {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Account with name %s already exists.", m.Name))
			return
		}
	}

	info, err := kb.CreateAccount(m.Name, seed, "", m.Password, account, index)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	keyOutput, err := cryptoKeys.Bech32KeyOutput(info)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	keyOutput.Mnemonic = seed

	bz, err := json.Marshal(keyOutput)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, _ = w.Write(bz)
}
