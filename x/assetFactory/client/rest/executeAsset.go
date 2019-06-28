package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	
	"github.com/asaskevich/govalidator"
	cliclient "github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	"github.com/comdex-blockchain/crypto/keys"
	"github.com/comdex-blockchain/rest"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/assetFactory"
	authctx "github.com/comdex-blockchain/x/auth/client/context"
)

// ExecuteAssetHandlerFunction : handles execute asset rest message
func ExecuteAssetHandlerFunction(cliCtx context.CLIContext, cdc *wire.Codec, kb keys.Keybase, kafka bool, kafkaState rest.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var msg assetFactory.ExecuteAssetBody
		
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		err = json.Unmarshal(body, &msg)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		_, err = govalidator.ValidateStruct(msg)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		txCtx := authctx.TxContext{
			Codec:         cdc,
			AccountNumber: msg.AccountNumber,
			Sequence:      msg.Sequence,
			Gas:           msg.Gas,
			ChainID:       msg.ChainID,
		}
		
		adjustment, ok := utils.ParseFloat64OrReturnBadRequest(w, msg.GasAdjustment, cliclient.DefaultGasAdjustment)
		if !ok {
			return
		}
		
		cliCtx = cliCtx.WithGasAdjustment(adjustment)
		cliCtx = cliCtx.WithFromAddressName(msg.From)
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
		
		toStr := msg.To
		to, err := sdk.AccAddressFromBech32(toStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		ownerStr := msg.Owner
		owner, err := sdk.AccAddressFromBech32(ownerStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		pegHash, err := sdk.GetAssetPegHashHex(msg.PegHash)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		buildmsg := assetFactory.BuildExecuteAssetMsg(from, owner, to, pegHash)
		
		if kafka == true {
			ticketID := rest.TicketIDGenerator("AFEA")
			
			jsonResponse := rest.SendToKafka(rest.NewKafkaMsgFromRest(buildmsg, ticketID, txCtx, cliCtx, msg.Password), kafkaState, cdc)
			w.WriteHeader(http.StatusAccepted)
			w.Write(jsonResponse)
		} else {
			
			output, err := utils.SendTxWithResponse(txCtx, cliCtx, []sdk.Msg{buildmsg}, msg.Password)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			w.Write(utils.ResponseBytesToJSON(output))
		}
	}
	
}
