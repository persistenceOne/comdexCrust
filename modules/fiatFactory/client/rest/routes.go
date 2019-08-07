package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/fiat/{pegHash}", QueryFiatRequestHandlerFn(cliCtx)).Methods("GET")
	
	r.HandleFunc("/fiat/issue", IssueFiatHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc("/fiat/send", SendFiatHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc("/fiat/execute", ExecuteFiatHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc("/fiat/redeem", RedeemFiatHandlerFunction(cliCtx)).Methods("POST")
}
