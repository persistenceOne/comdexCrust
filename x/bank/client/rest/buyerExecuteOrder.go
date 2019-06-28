package rest

import (
	"encoding/json"
	"fmt"
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
	"github.com/comdex-blockchain/x/acl"
	aclTypes "github.com/comdex-blockchain/x/acl"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/bank"
	"github.com/comdex-blockchain/x/bank/client"
)

// BuyerExecuteOrderRequestHandlerFn : handles rest request
func BuyerExecuteOrderRequestHandlerFn(cliCtx context.CLIContext, cdc *wire.Codec, kb keys.Keybase, kafka bool, kafkaState rest.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var msg bank.BuyerExecuteOrderBody
		
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
		
		_, err = govalidator.ValidateStruct(msg)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
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
		buyerAddressStr := msg.BuyerAddress
		
		buyerAddress, err := sdk.AccAddressFromBech32(buyerAddressStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		res, err := cliCtx.QueryStore(acl.AccountStoreKey(buyerAddress), "acl")
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
		
		zoneID := account.GetZoneID()
		if zoneID == nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}
		zoneData, err := cliCtx.QueryStore(acl.ZoneStoreKey(zoneID), "acl")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't parse query result. Result: %s. Error: %s", res, err.Error()))
			return
		}
		zoneAddress := sdk.AccAddress(string(zoneData))
		if err != nil || zoneAddress.String() != from.String() && from.String() != buyerAddress.String() {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}
		
		if !account.GetACL().BuyerExecuteOrder {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}
		
		sellerAddressStr := msg.SellerAddress
		
		sellerAddress, err := sdk.AccAddressFromBech32(sellerAddressStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		pegHashHex, err := sdk.GetAssetPegHashHex(msg.PegHash)
		sendmsg := client.BuildBuyerExecuteOrderMsg(from, buyerAddress, sellerAddress, pegHashHex, msg.FiatProofHash)
		
		if kafka == true {
			ticketID := rest.TicketIDGenerator("BBEO")
			jsonResponse := rest.SendToKafka(rest.NewKafkaMsgFromRest(sendmsg, ticketID, txCtx, cliCtx, msg.Password), kafkaState, cdc)
			w.WriteHeader(http.StatusAccepted)
			w.Write(jsonResponse)
		} else {
			output, err := utils.SendTxWithResponse(txCtx, cliCtx, []sdk.Msg{sendmsg}, msg.Password)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			w.Write(utils.ResponseBytesToJSON(output))
		}
	}
}
