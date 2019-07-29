package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/asaskevich/govalidator"
	cliclient "github.com/commitHub/commitBlockchain/client"
	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	"github.com/commitHub/commitBlockchain/crypto/keys"
	"github.com/commitHub/commitBlockchain/rest"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/acl"
	aclTypes "github.com/commitHub/commitBlockchain/x/acl"
	authctx "github.com/commitHub/commitBlockchain/x/auth/client/context"
	"github.com/commitHub/commitBlockchain/x/ibc"
)

// IssueFiatHandlerFunction - http request handler to Issue fiat
// on a different chain via IBC
// nolint: gocyclo
func IssueFiatHandlerFunction(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext, kafka bool, kafkaState rest.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var m ibc.IssueFiatBody

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = govalidator.ValidateStruct(m)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		adjustment, ok := utils.ParseFloat64OrReturnBadRequest(w, m.GasAdjustment, cliclient.DefaultGasAdjustment)
		if !ok {
			return
		}

		cliCtx = cliCtx.WithGasAdjustment(adjustment)
		cliCtx = cliCtx.WithFromAddressName(m.From)
		cliCtx.JSON = true

		txCtx := authctx.TxContext{
			Codec:         cdc,
			ChainID:       m.SourceChainID,
			AccountNumber: m.AccountNumber,
			Sequence:      m.Sequence,
			Gas:           m.Gas,
		}

		if err := cliCtx.EnsureAccountExists(); err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		from, err := cliCtx.GetFromAddress()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		toStr := m.To

		to, err := sdk.AccAddressFromBech32(toStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(acl.AccountStoreKey(from), "acl")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}

		// decode the value
		decoder := aclTypes.GetACLAccountDecoder(cdc)
		account, err := decoder(res)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't parse query result. Result: %s. Error: %s", res, err.Error()))
			return
		}
		if !account.GetACL().IssueFiat {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}
		pegHashHex, err := sdk.GetFiatPegHashHex(m.PegHash)
		fiatPeg := sdk.BaseFiatPeg{
			PegHash:           pegHashHex,
			TransactionID:     m.TransactionID,
			TransactionAmount: m.TransactionAmount,
		}

		msg := ibc.BuildIssueFiatMsg(from, to, &fiatPeg, m.SourceChainID, m.DestinationChainID)

		if kafka == true {
			ticketID := rest.TicketIDGenerator("IBCIF")
			jsonResponse := rest.SendToKafka(rest.NewKafkaMsgFromRest(msg, ticketID, txCtx, cliCtx, m.Password), kafkaState, cdc)
			w.WriteHeader(http.StatusAccepted)
			w.Write(jsonResponse)
		} else {
			output, err := utils.SendTxWithResponse(txCtx, cliCtx, []sdk.Msg{msg}, m.Password)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(utils.ResponseBytesToJSON(output))
		}
	}
}

// //IssueFiatKafkaHandlerFunction : handles rest request with kafka
// func IssueFiatKafkaHandlerFunction(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext, kafkaState rest.KafkaState) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")

// 		var m ibc.IssueFiatBody

// 		body, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		err = json.Unmarshal(body, &m)
// 		if err != nil {
// 			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		_, err = govalidator.ValidateStruct(m)
// 		if err != nil {
// 			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		ticketID := rest.TicketIDGenerator("IBIF")

// 		ticketID, err := cdc.MarshalJSON(ticketID)
// 		if err != nil {
// 			panic(err)
// 		}

// 		rest.SetTicketIDtoDB(ticketID, kafkaState.KafkaDB, cdc)
// 		rest.KafkaProducerDeliverMessage(&m, ticketID, "IBCIssueFiat", kafkaState.Producer, cdc)

// 		jsonResponse, err := cdc.MarshalJSON(struct {
// 			TicketID rest.Ticket `json:"ticketID"`
// 		}{TicketID: ticketID})
// 		if err != nil {
// 			panic(err)
// 		}

// 		w.WriteHeader(http.StatusAccepted)
// 		w.Write(jsonResponse)
// 	}
// }

// //IssueFiatKafkaMsgHandlerFunction : handles message from kafka
// func IssueFiatKafkaMsgHandlerFunction(cliCtx context.CLIContext, cdc *wire.Codec, kafkaState rest.KafkaState) {
// 	msgInterface, ticketID := rest.KafkaTopicConsumer("IBCIssueFiat", kafkaState.Consumers, cdc)
// 	m := msgInterface.(ibc.IssueFiatBody)
// 	ticketID, err := cdc.MarshalJSON(ticketID)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var adjustment float64
// 	if len(m.GasAdjustment) == 0 {
// 		adjustment = cliclient.DefaultGasAdjustment
// 	} else {

// 		adjustment, err = strconv.ParseFloat(m.GasAdjustment, 64)
// 		if err != nil {
// 			rest.AddResponseToDB(ticketID, []byte(err.Error()), kafkaState.KafkaDB, cdc)
// 		}
// 	}

// 	cliCtx = cliCtx.WithGasAdjustment(adjustment)
// 	cliCtx = cliCtx.WithFromAddressName(m.From)
// 	cliCtx.JSON = true

// 	if err := cliCtx.EnsureAccountExists(); err != nil {
// 		rest.AddResponseToDB(ticketID, []byte(err.Error()), kafkaState.KafkaDB, cdc)
// 		return
// 	}

// 	from, err := cliCtx.GetFromAddress()
// 	if err != nil {
// 		rest.AddResponseToDB(ticketID, []byte(err.Error()), kafkaState.KafkaDB, cdc)
// 		return
// 	}

// 	res, err := cliCtx.QueryStore(acl.AccountStoreKey(from), "acl")
// 	if err != nil {
// 		rest.AddResponseToDB(ticketID, []byte(err.Error()), kafkaState.KafkaDB, cdc)
// 		return
// 	}

// 	// the query will return empty if there is no data for this account
// 	if len(res) == 0 {
// 		rest.AddResponseToDB(ticketID, []byte("Unauthorised Transaction"), kafkaState.KafkaDB, cdc)
// 		return
// 	}

// 	// decode the value
// 	decoder := aclTypes.GetACLAccountDecoder(cdc)
// 	account, err := decoder(res)
// 	if err != nil {
// 		rest.AddResponseToDB(ticketID, []byte(err.Error()), kafkaState.KafkaDB, cdc)
// 		return
// 	}
// 	if !account.GetACL().IBCIssueFiats {
// 		rest.AddResponseToDB(ticketID, []byte("Unauthorised Transaction"), kafkaState.KafkaDB, cdc)
// 		return
// 	}

// 	toStr := m.To

// 	to, err := sdk.AccAddressFromBech32(toStr)
// 	if err != nil {
// 		rest.AddResponseToDB(ticketID, []byte(err.Error()), kafkaState.KafkaDB, cdc)
// 		return
// 	}
// 	fiatPeg := sdk.BaseFiatPeg{
// 		TransactionID:     m.TransactionID,
// 		TransactionAmount: m.TransactionAmount,
// 	}

// 	msg := ibc.BuildIssueFiatMsg(from, to, &fiatPeg, m.SourceChainID, m.DestinationChainID)

// 	txCtx := authctx.TxContext{
// 		Codec:         cdc,
// 		ChainID:       m.SourceChainID,
// 		AccountNumber: m.AccountNumber,
// 		Sequence:      m.Sequence,
// 		Gas:           m.Gas,
// 	}
// 	output, err := utils.SendTxWithResponse(txCtx, cliCtx, []sdk.Msg{msg}, m.Password)
// 	if err != nil {
// 		rest.AddResponseToDB(ticketID, []byte(err.Error()), kafkaState.KafkaDB, cdc)
// 		return
// 	}
// 	rest.AddResponseToDB(ticketID, output, kafkaState.KafkaDB, cdc)
// }
