package rest

import (
	"fmt"
	"net/http"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/auth"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	
	"github.com/gorilla/mux"
)

// RegisterRoutes : register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string) {
	r.HandleFunc(
		"/accounts/{address}",
		QueryAccountRequestHandlerFn(storeName, cdc, authcmd.GetAccountDecoder(cdc), cliCtx),
	).Methods("GET")
}

// QueryAccountRequestHandlerFn : query accountREST Handler
func QueryAccountRequestHandlerFn(storeName string, cdc *wire.Codec, decoder auth.AccountDecoder, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		
		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		res, err := cliCtx.QueryStore(auth.AddressStoreKey(addr), storeName)
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
