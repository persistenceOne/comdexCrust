package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/negotiation"
	"github.com/gorilla/mux"
)

var msgWireCdc = wire.NewCodec()

func init() {
	negotiation.RegisterWire(msgWireCdc)
}

// QueryNegotiationRequestHandlerFn : handler for query negotiation
func QueryNegotiationRequestHandlerFn(storeName string, cdc *wire.Codec, decoder sdk.NegotiationDecoder, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		vars := mux.Vars(r)
		nego := vars["negotiationID"]
		negotiationID, err := sdk.GetNegotiationIDHex(nego)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("cannot decode NegotiationID. Error: %s", err.Error()))
			return
		}
		
		res, err := cliCtx.QueryStore(negotiation.StoreKey(sdk.NegotiationID([]byte(negotiationID))), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't query Negotiation. Error: %s", err.Error()))
			return
		}
		
		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		// decode the value
		negotiationResponse, err := decoder(res)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't parse query result. Result: %s. Error: %s", res, err.Error()))
			return
		}
		
		// print out whole account
		output, err := cdc.MarshalJSON(negotiationResponse)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't marshall query result. Error: %s", err.Error()))
			return
		}
		
		w.Write(output)
	}
}

func GetNegotiationIDHandlerFn(cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		vars := mux.Vars(r)
		from := vars["from"]
		to := vars["to"]
		pegHash := vars["pegHash"]
		
		cliCtx = cliCtx.WithFromAddressName(from)
		cliCtx.JSON = true
		if err := cliCtx.EnsureAccountExists(); err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		buyerAddress, err := cliCtx.GetFromAddress()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		cliCtx = cliCtx.WithFromAddressName(to)
		cliCtx.JSON = true
		if err := cliCtx.EnsureAccountExists(); err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		sellerAddress, err := cliCtx.GetFromAddress()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		pegHashHex, err := sdk.GetAssetPegHashHex(pegHash)
		negotiation := sdk.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHashHex.Bytes()...))
		
		output, err := json.Marshal(struct {
			NegotiationID string `json:"negotiationID"`
		}{NegotiationID: negotiation.String()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.Write(output)
	}
}
