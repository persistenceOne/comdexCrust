package rest

import (
	"fmt"
	"net/http"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	
	"github.com/commitHub/commitBlockchain/modules/negotiation"
	"github.com/commitHub/commitBlockchain/modules/orders/internal/keeper"
	"github.com/commitHub/commitBlockchain/modules/orders/internal/types"
)

func QueryOrderRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		
		negotiationID, err := negotiation.GetNegotiationIDFromString(vars["negotiation-id"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("cannot decode NegotiationID. Error: %s", err.Error()))
			return
		}
		
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, keeper.QueryOrder, negotiationID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		var order types.Order
		cliCtx.Codec.MustUnmarshalJSON(res, &order)
		
		rest.PostProcessResponse(w, cliCtx, order)
	}
}
