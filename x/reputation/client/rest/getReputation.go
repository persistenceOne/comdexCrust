package rest

import (
	"fmt"
	"net/http"

	"github.com/commitHub/commitBlockchain/x/reputation"
	"github.com/gorilla/mux"

	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
)

var msgWireCdc = wire.NewCodec()

func init() {
	reputation.RegisterWire(msgWireCdc)
}

// QueryReputationRequestHandlerFn : handler for query reputation
func QueryReputationRequestHandlerFn(storeName string, cdc *wire.Codec, decoder sdk.ReputationDecoder, cliContext context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		cliCtx := cliContext

		bech32addr := vars["address"]
		if bech32addr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(reputation.AccountStoreKey(addr), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// decode the value
		account, err := decoder(res)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't parse query result. Result: %s. Error: %s", res, err.Error()))
			return
		}

		// print out whole account
		output, err := cdc.MarshalJSON(account)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't marshall query result. Error: %s", err.Error()))
			return
		}

		w.Write(output)
	}
}
