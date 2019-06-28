package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	
	"github.com/asaskevich/govalidator"
	cliclient "github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	"github.com/comdex-blockchain/rest"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/acl"
	authctx "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/bank"
)

// DefineZoneHandler : to add zone
func DefineZoneHandler(cdc *wire.Codec, cliCtx context.CLIContext, kafka bool, kafkaState rest.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var msg acl.DefineZoneBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = json.Unmarshal(body, &msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		_, err = govalidator.ValidateStruct(msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		cliCtx = cliCtx.WithFromAddressName(msg.From)
		from, err := cliCtx.GetFromAddress()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		to, err := sdk.AccAddressFromBech32(msg.To)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		txCtx := authctx.TxContext{
			Codec:         cdc,
			AccountNumber: msg.AccountNumber,
			Sequence:      msg.Sequence,
			Gas:           msg.Gas,
			ChainID:       msg.ChainID,
		}
		zoneID, err := sdk.GetZoneIDFromString(msg.ZoneID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		msgZone := bank.BuildMsgDefineZone(from, to, zoneID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		
		adjustment, ok := utils.ParseFloat64OrReturnBadRequest(w, msg.GasAdjustment, cliclient.DefaultGasAdjustment)
		if !ok {
			return
		}
		cliCtx = cliCtx.WithGasAdjustment(adjustment)
		cliCtx.JSON = true
		
		if err := cliCtx.EnsureAccountExists(); err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		AccNumber, err := cliCtx.GetAccountNumber(from)
		if AccNumber != 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("You are not a genesis account"))
			return
		}
		if kafka == true {
			ticketID := rest.TicketIDGenerator("ACDZ")
			
			jsonResponse := rest.SendToKafka(rest.NewKafkaMsgFromRest(msgZone, ticketID, txCtx, cliCtx, msg.Password), kafkaState, cdc)
			w.WriteHeader(http.StatusAccepted)
			w.Write(jsonResponse)
		} else {
			output, err := utils.SendTxWithResponse(txCtx, cliCtx, []sdk.Msg{msgZone}, msg.Password)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}
			
			w.Write(utils.ResponseBytesToJSON(output))
		}
	}
}
