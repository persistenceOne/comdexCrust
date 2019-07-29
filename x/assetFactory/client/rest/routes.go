package rest

import (
	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/crypto/keys"
	"github.com/commitHub/commitBlockchain/rest"
	"github.com/commitHub/commitBlockchain/wire"
	assetcmd "github.com/commitHub/commitBlockchain/x/assetFactory/client/cli"
	"github.com/gorilla/mux"
)

//RegisterRoutes : registers routes in root
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string, kb keys.Keybase, kafka bool, kafkaState rest.KafkaState) {
	r.HandleFunc("/asset/{pegHash}", QueryAssetRequestHandlerFn(storeName, cdc, r, assetcmd.GetAssetPegDecoder(cdc), cliCtx)).Methods("GET")
	r.HandleFunc("/issue/asset", IssueAssetRestHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/redeem/asset", RedeemAssetHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/send/asset", SendAssetHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/execute/asset", ExecuteAssetHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
}
