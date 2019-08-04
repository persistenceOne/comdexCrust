package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/reputation/{address}", QueryReputationRequestHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc("/submitBuyerFeedback", SubmitBuyerFeedbackRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/submitSellerFeedback", SubmitSellerFeedbackRequestHandler(cliCtx)).Methods("POST")
}
