package rest

import (
	"fmt"
	"net/http"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	
	"github.com/commitHub/commitBlockchain/modules/reputation/internal/types"
)

func QueryReputationRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		cliCtx := cliCtx
		
		bech32addr := vars["address"]
		if bech32addr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, "reputationQuery", bech32addr), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}
		
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		var reputation types.AccountReputation
		cliCtx.Codec.MustUnmarshalJSON(res, &reputation)
		
		rest.PostProcessResponse(w, cliCtx, reputation)
	}
}
