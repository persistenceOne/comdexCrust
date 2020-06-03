package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/go-bip39"
	"github.com/gorilla/mux"
)

type ForgotPasswordBody struct {
	Seed               string `json:"seed"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}

func ForgotPasswordRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		name := vars["name"]

		var m ForgotPasswordBody
		var index uint32
		var account uint32
		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &m)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		kb, err := keys.NewKeyBaseFromHomeFlag()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = kb.Get(name)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if !bip39.IsMnemonicValid(m.Seed) {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "invalid mnemonic")
			return
		}
		if m.NewPassword != m.ConfirmNewPassword {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "Password do not match")
			return
		}
		_, err = kb.CreateAccount(name, m.Seed, "", m.NewPassword, account, index)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, _ = w.Write([]byte(fmt.Sprintf(`{"error":%t ,"message":"%v"}`, false, "Password updated successfully")))
		w.WriteHeader(200)
	}
}
