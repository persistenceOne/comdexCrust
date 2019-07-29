package rest

import (
	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/crypto/keys"
	"github.com/commitHub/commitBlockchain/rest"
	"github.com/commitHub/commitBlockchain/wire"
	fiatcmd "github.com/commitHub/commitBlockchain/x/fiatFactory/client/cli"
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
