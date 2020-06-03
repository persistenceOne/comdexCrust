package rest

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type UpdateKeyBody struct {
	NewPassword        string `json:"newPassword"`
	OldPassword        string `json:"oldPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}

func UpdateKeyRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		name := vars["name"]

		var m UpdateKeyBody
		kb, err := keys.NewKeyBaseFromHomeFlag()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &m)

		if m.NewPassword != m.ConfirmNewPassword {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "Password do not match")
			return
		}
		getNewPass := func() (string, error) { return m.NewPassword, nil }

		info, err := kb.List()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		for i, in := range info {
			if in.GetName() != name && i == len(info) {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Account with name %s does not exits", name))
				return
			} else {
				continue
			}
		}
		_, err = kb.ExportPrivateKeyObject(name, m.OldPassword)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"error":%t ,"message":"%v"}`, true, err.Error())))
			return
		}
		err = kb.Update(name, m.OldPassword, getNewPass)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error":%t ,"message":"%v"}`, false, "Password updated successfully")))
		w.WriteHeader(200)
	}
}
