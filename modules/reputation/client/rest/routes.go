package rest

import (
	"github.com/commitHub/commitBlockchain/kafka"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, kafkaBool bool, kafkaState kafka.KafkaState) {
	r.HandleFunc("/reputation/{address}", QueryReputationRequestHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc("/submitBuyerFeedback", SubmitBuyerFeedbackRequestHandler(cliCtx, kafkaBool, kafkaState)).Methods("POST")
	r.HandleFunc("/submitSellerFeedback", SubmitSellerFeedbackRequestHandler(cliCtx, kafkaBool, kafkaState)).Methods("POST")
}
