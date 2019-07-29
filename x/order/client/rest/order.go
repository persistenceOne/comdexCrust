package rest

import (
	"fmt"
	"net/http"

	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	"github.com/commitHub/commitBlockchain/crypto/keys"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/order"
	ordercmd "github.com/commitHub/commitBlockchain/x/order/client/cli"
	"github.com/gorilla/mux"
)

//RegisterRoutes :
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string, kb keys.Keybase) {
	r.HandleFunc("/order/{negotiationID}", QueryOrderRequestHandlerFn(storeName, cdc, ordercmd.GetOrderDecoder(cdc), cliCtx, kb)).Methods("GET")
}

//QueryOrderRequestHandlerFn :
func QueryOrderRequestHandlerFn(storeName string, cdc *wire.Codec, decoder sdk.OrderDecoder, cliContext context.CLIContext, kb keys.Keybase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		nego := vars["negotiationID"]
		cliCtx := cliContext
		negotiationID, err := sdk.GetNegotiationIDHex(nego)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("cannot decode NegotiationID. Error: %s", err.Error()))
			return
		}

		res, err := cliCtx.QueryStore(order.StoreKey([]byte(negotiationID)), storeName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		order, err := decoder(res)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't parse query result. Result: %s. Error: %s", res, err.Error()))
			return
		}
		output, err := cdc.MarshalJSON(order)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't marshall query result. Error: %s", err.Error()))
			return
		}

		w.Write(output)
	}
}
