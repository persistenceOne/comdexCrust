package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	"github.com/persistenceOne/persistenceSDK/kafka"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, kafkaBool bool, kafkaState kafka.KafkaState) {
	r.HandleFunc("/bank/balances/{address}", QueryBalancesRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/asset/{peghash}", QueryAssetHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/fiat/{peghash}", QueryFiatHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc("/bank/accounts/{address}/transfers", SendRequestHandlerFn(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/defineZone", DefineZoneHandler(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/defineOrganization", DefineOrganizationHandler(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/issueAsset", IssueAssetHandlerFunction(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/redeemAsset", RedeemAssetHandlerFunction(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/issueFiat", IssueFiatHandlerFunction(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/sendAsset", SendAssetRequestHandlerFn(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/sendFiat", SendFiatRequestHandlerFn(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/redeemFiat", RedeemFiatHandlerFunction(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/releaseAsset", ReleaseAssetHandlerFunction(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/buyerExecuteOrder", BuyerExecuteOrderRequestHandlerFn(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/sellerExecuteOrder", SellerExecuteOrderRequestHandlerFn(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/defineACL", DefineACLHandler(cliCtx, kafkaBool, kafkaState)).Methods("POST")
}
