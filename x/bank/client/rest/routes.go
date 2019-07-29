package rest

import (
	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/crypto/keys"
	"github.com/commitHub/commitBlockchain/rest"
	"github.com/commitHub/commitBlockchain/wire"
	assetcmd "github.com/commitHub/commitBlockchain/x/assetFactory/client/cli"
	fiatcmd "github.com/commitHub/commitBlockchain/x/fiatFactory/client/cli"
	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase, kafka bool, kafkaState rest.KafkaState, storeName string, fiatStoreName string) {
	r.HandleFunc("/sendCoin", SendRequestHandlerFn(cdc, kb, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/defineZone", DefineZoneHandler(cdc, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/defineOrganization", DefineOrganizationHandler(cdc, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/defineACL", DefineACLHandler(cdc, cliCtx, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/issueAsset", IssueAssetHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/redeemAsset", RedeemAssetHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/issueFiat", IssueFiatHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/sendAsset", SendAssetRequestHandlerFn(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/sendFiat", SendFiatRequestHandlerFn(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/redeemFiat", RedeemFiatHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/releaseAsset", ReleaseAssetHandlerFunction(cliCtx, cdc, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/buyerExecuteOrder", BuyerExecuteOrderRequestHandlerFn(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/sellerExecuteOrder", SellerExecuteOrderRequestHandlerFn(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/asset/{peghash}", QueryAssetHandlerFn(cliCtx, r, cdc, storeName, assetcmd.GetAssetPegDecoder(cdc))).Methods("GET")
	r.HandleFunc("/fiat/{peghash}", QueryFiatHandlerFn(cliCtx, r, cdc, fiatStoreName, fiatcmd.GetFiatPegDecoder(cdc))).Methods("GET")
}
