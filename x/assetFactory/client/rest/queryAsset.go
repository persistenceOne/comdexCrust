package rest

import (
	"fmt"
	"net/http"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/assetFactory"
	"github.com/gorilla/mux"
)

// QueryAssetRequestHandlerFn : query the issued asset
func QueryAssetRequestHandlerFn(storeName string, cdc *wire.Codec, r *mux.Router, decoder sdk.AssetPegDecoder, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		pegHashStr := vars["pegHash"]
		
		pegHashHex, err := sdk.GetAssetPegHashHex(pegHashStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		res, err := cliCtx.QueryStore(assetFactory.AssetPegHashStoreKey(pegHashHex), storeName)
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
