package keys

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func updateKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <name>",
		Short: "Change the password used to protect private key",
		RunE:  runUpdateCmd,
		Args:  cobra.ExactArgs(1),
	}
	return cmd
}

func runUpdateCmd(cmd *cobra.Command, args []string) error {
	name := args[0]

	buf := bufio.NewReader(cmd.InOrStdin())
	kb, err := NewKeyBaseFromHomeFlag()
	if err != nil {
		return err
	}
	oldpass, err := input.GetPassword("Enter the current passphrase:", buf)
	if err != nil {
		return err
	}

	getNewpass := func() (string, error) {
		return input.GetCheckPassword(
			"Enter the new passphrase:",
			"Repeat the new passphrase:", buf)
	}
	if err := kb.Update(name, oldpass, getNewpass); err != nil {
		return err
	}

	cmd.PrintErrln("Password successfully updated!")
	return nil
}

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
		kb, err := NewKeyBaseFromHomeFlag()
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
