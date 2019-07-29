package rest

import (
	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/crypto/keys"
	"github.com/commitHub/commitBlockchain/rest"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/reputation/client/cli"
	"github.com/gorilla/mux"
)

//RegisterRoutes : register all the routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string, kb keys.Keybase, kafka bool, kafkaState rest.KafkaState) {
	r.HandleFunc("/submitBuyerFeedback", SubmitBuyerFeedbackRequestHandler(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/submitSellerFeedback", SubmitSellerFeedbackRequestHandler(cliCtx, cdc, kb, kafka, kafkaState)).Methods("POST")
	r.HandleFunc("/reputation/{address}", QueryReputationRequestHandlerFn("reputation", cdc, cli.GetReputationDecoder(cdc), cliCtx)).Methods("GET")
}
