package keys

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/spf13/cobra"
)

func listKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all keys",
		Long: `Return a list of all public keys stored by this key manager
along with their associated name and address.`,
		RunE: runListCmd,
	}
	cmd.Flags().Bool(flags.FlagIndentResponse, false, "Add indent to JSON response")
	return cmd
}

func runListCmd(cmd *cobra.Command, args []string) error {
	kb, err := NewKeyBaseFromHomeFlag()
	if err != nil {
		return err
	}

	infos, err := kb.List()
	if err == nil {
		printInfos(infos)
	}
	return err
}

func QueryKeysRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	kb, err := NewKeyBaseFromHomeFlag()
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	info, err := kb.List()
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	kos, err := keys.Bech32KeysOutput(info)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	bz, err := json.MarshalIndent(kos, "", "")
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, _ = w.Write(bz)
}
