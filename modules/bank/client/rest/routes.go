package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/bank/balances/{address}", QueryBalancesRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/asset/{peghash}", QueryAssetHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/fiat/{peghash}", QueryFiatHandlerFn(cliCtx)).Methods("GET")
	
	r.HandleFunc("/bank/accounts/{address}/transfers", SendRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/defineZone", DefineZoneHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/defineOrganization", DefineOrganizationHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/defineACL", DefineACLHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/issueAsset", IssueAssetHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc("/redeemAsset", RedeemAssetHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc("/issueFiat", IssueFiatHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc("/sendAsset", SendAssetRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/sendFiat", SendFiatRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/redeemFiat", RedeemFiatHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc("/releaseAsset", ReleaseAssetHandlerFunction(cliCtx)).Methods("POST")
	r.HandleFunc("/buyerExecuteOrder", BuyerExecuteOrderRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/sellerExecuteOrder", SellerExecuteOrderRequestHandlerFn(cliCtx)).Methods("POST")
}
