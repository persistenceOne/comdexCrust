package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	rest2 "github.com/persistenceOne/comdexCrust/client/rest"
	negotiationTypes "github.com/persistenceOne/comdexCrust/modules/negotiation/internal/types"
	"github.com/persistenceOne/comdexCrust/types"
)

func QueryNegotiationRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cliCtx := cliCtx
		vars := mux.Vars(r)

		negotiationID, err := types.GetNegotiationIDFromString(vars["negotiation-id"])
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrNegotiationIDFromString(negotiationTypes.DefaultCodeSpace, vars["negotiation-id"]))
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", negotiationTypes.QuerierRoute, "queryNegotiation", negotiationID), nil)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(negotiationTypes.DefaultCodeSpace, "negotiation"))
			return
		}

		if len(res) == 0 {
			rest2.WriteErrorResponse(w, types.ErrQueryResponseLengthZero(negotiationTypes.DefaultCodeSpace, "negotiation"))
			return
		}

		var _negotiation types.Negotiation
		cliCtx.Codec.MustUnmarshalJSON(res, &_negotiation)

		rest.PostProcessResponse(w, cliCtx, _negotiation)
	}
}

func GetNegotiationIDHandlerFn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		from := vars["buyerAddress"]
		to := vars["sellerAddress"]
		pegHash := vars["pegHash"]

		buyerAddress, _, err := context.GetFromFields(from, false)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrFromName(negotiationTypes.DefaultCodeSpace))
			return
		}
		sellerAddress, _, err := context.GetFromFields(to, false)
		fmt.Println(sellerAddress)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrFromName(negotiationTypes.DefaultCodeSpace))
			return
		}

		pegHashHex, err := types.GetAssetPegHashHex(pegHash)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrPegHashHex(negotiationTypes.DefaultCodeSpace, pegHash))
			return
		}
		negotiation := types.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHashHex.Bytes()...))

		output, err := json.Marshal(struct {
			NegotiationID string `json:"negotiationID"`
		}{NegotiationID: negotiation.String()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(output)
	}
}
