package rest

import (
	"fmt"
	"net/http"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	
	aclTypes "github.com/commitHub/commitBlockchain/modules/acl/internal/types"
)

// GetACLRequestHandler : query Acl account details
func GetACLRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		cliCtx := cliCtx
		
		addr, err := cTypes.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", aclTypes.QuerierRoute, "queryACLAccount", addr), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}
		
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		var account aclTypes.ACLAccount
		cliCtx.Codec.MustUnmarshalJSON(res, &account)
		
		rest.PostProcessResponse(w, cliCtx, account)
	}
}
