package rest

import (
	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/crypto/keys"
	"github.com/commitHub/commitBlockchain/rest"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, kafka bool, kafkaState rest.KafkaState) {
	r.HandleFunc("/ibcIssueAsset", IssueAssetHandlerFunction(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/ibcRedeemAsset", RedeemAssetHandlerFunction(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/ibcIssueFiat", IssueFiatHandlerFunction(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/ibcRedeemFiat", RedeemFiatHandlerFunction(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/ibcSendAsset", SendAssetHandlerFunction(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/ibcSendFiat", SendFiatHandlerFunction(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/ibc/{destchain}/{address}/send", TransferRequestHandlerFn(cdc, kb, cliCtx)).Methods("POST")
	r.HandleFunc("/ibcBuyerExecuteOrder", BuyerExecuteOrderHandlerFunction(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/ibcSellerExecuteOrder", SellerExecuteOrderHandlerFuncion(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")

}

// RegisterRoutesKafka : Central function to define routes that get registered by the main application
// func RegisterRoutesKafka(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, kafkaState rest.KafkaState) {
// 	r.HandleFunc("/ibcIssueAsset", IssueAssetKafkaHandlerFunction(cdc, kb, cliCtx, kafkaState)).Methods("POST")
// 	r.HandleFunc("/ibcIssueFiat", IssueFiatKafkaHandlerFunction(cdc, kb, cliCtx, kafkaState)).Methods("POST")
// 	r.HandleFunc("/ibc/{destchain}/{address}/send", TransferRequestHandlerFn(cdc, kb, cliCtx)).Methods("POST")

// }
