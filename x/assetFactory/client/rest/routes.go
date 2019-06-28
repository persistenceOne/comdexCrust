package rest

import (
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/crypto/keys"
	"github.com/comdex-blockchain/rest"
	"github.com/comdex-blockchain/wire"
	assetcmd "github.com/comdex-blockchain/x/assetFactory/client/cli"
	"github.com/gorilla/mux"
)

// RegisterRoutes : registers routes in root
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string, kb keys.Keybase, kafka bool, kafkaState rest.KafkaState) {
	r.HandleFunc("/asset/{pegHash}", QueryAssetRequestHandlerFn(storeName, cdc, r, assetcmd.GetAssetPegDecoder(cdc), cliCtx)).Methods("GET")
	r.HandleFunc("/issue/asset", IssueAssetRestHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/redeem/asset", RedeemAssetHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/send/asset", SendAssetHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/execute/asset", ExecuteAssetHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
}
