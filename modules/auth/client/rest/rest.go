package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers the auth module REST routes.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(
		"/auth/accounts/{address}", QueryAccountRequestHandlerFn(storeName, cliCtx),
	).Methods("GET")
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", QueryTxsByEventsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/txs", BroadcastTxRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/txs/encode", EncodeTxRequestHandlerFn(cliCtx)).Methods("POST")
}
