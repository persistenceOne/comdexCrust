package rest

import (
	"encoding/json"
	"fmt"
	"github.com/commitHub/commitBlockchain/types"
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

func GetNegotiationIDHandlerFn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		from := vars["buyerAddress"]
		to := vars["sellerAddress"]
		pegHash := vars["pegHash"]

		buyerAddress, _, err := context.GetFromFields(from, false)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't query Negotiation. Error: %s", err.Error()))
			return
		}
		sellerAddress, _, err := context.GetFromFields(to, false)
		fmt.Println(sellerAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't query Negotiation. Error: %s", err.Error()))
			return
		}

		pegHashHex, err := types.GetAssetPegHashHex(pegHash)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't query Negotiation. Error: %s", err.Error()))
			return
		}
		negotiation := negotiationTypes.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHashHex.Bytes()...))

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
