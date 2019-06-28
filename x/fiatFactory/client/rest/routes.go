package rest

import (
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/crypto/keys"
	"github.com/comdex-blockchain/rest"
	"github.com/comdex-blockchain/wire"
	fiatcmd "github.com/comdex-blockchain/x/fiatFactory/client/cli"
	"github.com/gorilla/mux"
)

// RegisterRoutes : routes to register
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string, kb keys.Keybase, kafka bool, kafkaState rest.KafkaState) {
	r.HandleFunc("/issue/fiat", IssueFiatHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/fiat/{pegHash}", QueryFiatRequestHandlerFn(storeName, cdc, r, fiatcmd.GetFiatPegDecoder(cdc), cliCtx)).Methods("GET")
	r.HandleFunc("/redeem/fiat", RedeemFiatHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/send/fiat", SendFiatHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/execute/fiat", ExecuteFiatHandlerFunction(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
}
