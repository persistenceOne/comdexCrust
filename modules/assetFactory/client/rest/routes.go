package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/fiat/{peg-hash}", QueryAssetRequestHandlerFn(cliCtx)).Methods("GET")
	
	r.HandleFunc("/asset/issue", IssueAssetHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc("/asset/send", SendAssetHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc("/asset/execute", ExecuteAssetHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc("/asset/redeem", RedeemAssetHandlerFunction(cliCtx)).Methods("POST")
}
