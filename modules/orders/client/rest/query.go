package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	rest2 "github.com/persistenceOne/comdexCrust/client/rest"
	"github.com/persistenceOne/comdexCrust/modules/orders/internal/keeper"
	orderTypes "github.com/persistenceOne/comdexCrust/modules/orders/internal/types"
	"github.com/persistenceOne/comdexCrust/types"
)

func QueryOrderRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		negotiationID, err := types.GetNegotiationIDFromString(vars["negotiation-id"])
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrNegotiationIDFromString(orderTypes.DefaultCodeSpace, vars["negotiation-id"]))
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", orderTypes.QuerierRoute, keeper.QueryOrder, negotiationID), nil)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(orderTypes.DefaultCodeSpace, "order"))
			return
		}

		if len(res) == 0 {
			rest2.WriteErrorResponse(w, types.ErrQueryResponseLengthZero(orderTypes.DefaultCodeSpace, "order"))
			return
		}

		var order types.Order
		cliCtx.Codec.MustUnmarshalJSON(res, &order)

		rest.PostProcessResponse(w, cliCtx, order)
	}
}
