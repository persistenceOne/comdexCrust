package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

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

// IssueAssetHandlerFunction - http request handler to IssueAsset
// on a different chain via IBC
// nolint: gocyclo
func IssueAssetHandlerFunction(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext, kafka bool, kafkaState rest.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var m ibc.IssueAssetBody
		var to, takerAddress sdk.AccAddress

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

		txCtx := authctx.TxContext{
			Codec:         cdc,
			ChainID:       m.SourceChainID,
			AccountNumber: m.AccountNumber,
			Sequence:      m.Sequence,
			Gas:           m.Gas,
		}

		adjustment, ok := utils.ParseFloat64OrReturnBadRequest(w, m.GasAdjustment, cliclient.DefaultGasAdjustment)
		if !ok {
			return
		}

		cliCtx = cliCtx.WithGasAdjustment(adjustment)
		cliCtx = cliCtx.WithFromAddressName(m.From)
		cliCtx.JSON = true

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
		if (m.Moderated && toStr == "") || (!m.Moderated && toStr != "") {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("private variable is not valid for toAddress."))
			return
		}

		if toStr == "" {
			to = from
		} else {
			to, err = sdk.AccAddressFromBech32(toStr)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			if !m.Moderated {
				if to.String() != from.String() {
					utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Cannot issue an issue. receiverAddress should same as issuerAddress."))
					return
				}
			}
		}
		if m.TakerAddress != "" {
			takerAddress, err = sdk.AccAddressFromBech32(m.TakerAddress)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		res, err := cliCtx.QueryStore(acl.AccountStoreKey(to), "acl")
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
		if !account.GetACL().IssueAsset {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}
		if m.Moderated {
			zoneID := account.GetZoneID()
			zoneData, err := cliCtx.QueryStore(acl.ZoneStoreKey(zoneID), "acl")
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
				return
			}
			zoneAcc := sdk.AccAddress(string(zoneData))
			if !reflect.DeepEqual(from, zoneAcc) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("you are not authorized person. Only zones can issue assets."))
				return
			}
		}
		assetPeg := &sdk.BaseAssetPeg{
			AssetQuantity: m.AssetQuantity,
			AssetType:     m.AssetType,
			AssetPrice:    m.AssetPrice,
			DocumentHash:  m.DocumentHash,
			QuantityUnit:  m.QuantityUnit,
			Moderated:     m.Moderated,
			TakerAddress:  takerAddress,
		}

		msg := ibc.BuildIssueAssetMsg(from, to, assetPeg, m.SourceChainID, m.DestinationChainID)
		// build message

		if kafka == true {
			ticketID := rest.TicketIDGenerator("IBCIA")
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

// //IssueAssetKafkaHandlerFunction : handles rest request with kafka
// func IssueAssetKafkaHandlerFunction(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext, kafkaState rest.KafkaState) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")

// 		var m ibc.IssueAssetBody

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

// 		ticketID := rest.TicketIDGenerator("IBIA")

// 		rest.SetTicketIDtoDB(ticketID, kafkaState.KafkaDB, cdc)
// 		rest.KafkaProducerDeliverMessage(&m, ticketID, "IBCIssueAsset", kafkaState.Producer, cdc)

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

// //IssueAssetKafkaMsgHandlerFunction : handles message from kafka
// func IssueAssetKafkaMsgHandlerFunction(cliCtx context.CLIContext, cdc *wire.Codec, kafkaState rest.KafkaState) {
// 	msgInterface, ticketID := rest.KafkaTopicConsumer("IBCIssueAsset", kafkaState.Consumers, cdc)
// 	m := msgInterface.(ibc.IssueAssetBody)
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
// 	if !account.GetACL().IBCIssueAssets {
// 		rest.AddResponseToDB(ticketID, []byte("Unathorised Transaction"), kafkaState.KafkaDB, cdc)
// 		return
// 	}

// 	toStr := m.To

// 	to, err := sdk.AccAddressFromBech32(toStr)
// 	if err != nil {
// 		rest.AddResponseToDB(ticketID, []byte(err.Error()), kafkaState.KafkaDB, cdc)
// 		return
// 	}

// 	assetPeg := &sdk.BaseAssetPeg{
// 		AssetQuantity: m.AssetQuantity,
// 		AssetType:     m.AssetType,
// 		AssetPrice:    m.AssetPrice,
// 		DocumentHash:  m.DocumentHash,
// 		QuantityUnit:  m.QuantityUnit,
// 	}

// 	msg := ibc.BuildIssueAssetMsg(from, to, assetPeg, m.SourceChainID, m.DestinationChainID)
// 	// build message

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
