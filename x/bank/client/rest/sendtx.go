package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	
	cliclient "github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	"github.com/comdex-blockchain/crypto/keys"
	"github.com/comdex-blockchain/rest"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/bank"
	"github.com/comdex-blockchain/x/bank/client"
)

var msgCdc = wire.NewCodec()

func init() {
	bank.RegisterWire(msgCdc)
}

// SendRequestHandlerFn - http request handler to send coins to a address
// nolint: gocyclo
func SendRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext, kafka bool, kafkaState rest.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var msg bank.SendTxBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		err = json.Unmarshal(body, &msg)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		adjustment, ok := utils.ParseFloat64OrReturnBadRequest(w, msg.GasAdjustment, cliclient.DefaultGasAdjustment)
		if !ok {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		cliCtx = cliCtx.WithGasAdjustment(adjustment)
		cliCtx = cliCtx.WithFromAddressName(msg.From)
		cliCtx.JSON = true
		
		txCtx := context2.TxContext{
			Codec:         cdc,
			ChainID:       msg.ChainID,
			AccountNumber: msg.AccountNumber,
			Sequence:      msg.Sequence,
			Gas:           msg.Gas,
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
		
		toStr := msg.To
		to, err := sdk.AccAddressFromBech32(toStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		sendMsg := client.BuildMsg(from, to, msg.Amount)
		
		if kafka == true {
			ticketID := rest.TicketIDGenerator("BKST")
			
			jsonResponse := rest.SendToKafka(rest.NewKafkaMsgFromRest(sendMsg, ticketID, txCtx, cliCtx, msg.Password), kafkaState, cdc)
			w.WriteHeader(http.StatusAccepted)
			w.Write(jsonResponse)
		} else {
			output, err := utils.SendTxWithResponse(txCtx, cliCtx, []sdk.Msg{sendMsg}, msg.Password)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			
			w.Write(utils.ResponseBytesToJSON(output))
		}
	}
}
