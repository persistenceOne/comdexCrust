package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/negotiation/{negotiation-id}", QueryNegotiationRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/negotiationID/{buyerAddress}/{sellerAddress}/{pegHash}", GetNegotiationIDHandlerFn()).Methods("GET")
	r.HandleFunc("/changeBuyerBid", ChangeBuyerBidRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/changeSellerBid", ChangeSellerBidRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/confirmBuyerBid", ConfirmBuyerBidRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/confirmSellerBid", ConfirmSellerBidRequestHandlerFn(cliCtx)).Methods("POST")
}
