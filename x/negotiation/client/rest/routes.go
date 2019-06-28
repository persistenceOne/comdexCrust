package rest

import (
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/crypto/keys"
	"github.com/comdex-blockchain/rest"
	"github.com/comdex-blockchain/wire"
	negotiation "github.com/comdex-blockchain/x/negotiation/client/cli"
	"github.com/gorilla/mux"
)

// RegisterRoutes : register all the routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string, kb keys.Keybase, kafka bool, kafkaState rest.KafkaState) {
	r.HandleFunc("/changeBuyerBid", ChangeBuyerBidRequestHandlerFn(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/changeSellerBid", ChangeSellerBidRequestHandlerFn(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/confirmBuyerBid", ConfirmBuyerBidRequestHandlerFn(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/confirmSellerBid", ConfirmSellerBidRequestHandlerFn(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/negotiation/{negotiationID}", QueryNegotiationRequestHandlerFn(storeName, cdc, negotiation.GetNegotiationDecoder(cdc), cliCtx)).Methods("GET")
	r.HandleFunc("/negotiationID/{from}/{to}/{pegHash}", GetNegotiationIDHandlerFn(cdc, cliCtx)).Methods("GET")
}
