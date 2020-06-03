package rest

import (
	"github.com/commitHub/commitBlockchain/kafka"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, kafkaBool bool, kafkaState kafka.KafkaState) {
	r.HandleFunc("/negotiation/{negotiation-id}", QueryNegotiationRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/negotiationID/{buyerAddress}/{sellerAddress}/{pegHash}", GetNegotiationIDHandlerFn()).Methods("GET")
	r.HandleFunc("/changeBuyerBid", ChangeBuyerBidRequestHandlerFn(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/changeSellerBid", ChangeSellerBidRequestHandlerFn(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/confirmBuyerBid", ConfirmBuyerBidRequestHandlerFn(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/confirmSellerBid", ConfirmSellerBidRequestHandlerFn(cliCtx, kafkaBool, kafkaState)).Methods("POST")
}
