package rest

import (
	"fmt"
	"net/http"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	
	negotiationTypes "github.com/commitHub/commitBlockchain/modules/negotiation/internal/types"
)

func QueryNegotiationRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		cliCtx := cliCtx
		vars := mux.Vars(r)
		
		negotiationID, err := negotiationTypes.GetNegotiationIDFromString(vars["negotiation-id"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("cannot decode NegotiationID. Error: %s", err.Error()))
			return
		}
		
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", negotiationTypes.QuerierRoute, "queryNegotiation", negotiationID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't query Negotiation. Error: %s", err.Error()))
			return
		}
		
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		var _negotiation negotiationTypes.Negotiation
		cliCtx.Codec.MustUnmarshalJSON(res, &_negotiation)
		
		rest.PostProcessResponse(w, cliCtx, _negotiation)
	}
}
