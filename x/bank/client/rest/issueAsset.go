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
	"github.com/commitHub/commitBlockchain/x/bank"
	"github.com/commitHub/commitBlockchain/x/bank/client"
)

var msgWireCdc = wire.NewCodec()

func init() {
	bank.RegisterWire(msgWireCdc)
}

//IssueAssetHandlerFunction : handles issue asset rest message
func IssueAssetHandlerFunction(cliContext context.CLIContext, cdc *wire.Codec, kb keys.Keybase, kafka bool, kafkaState rest.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var msg bank.IssueAssetBody
		cliCtx := cliContext

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
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
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
		var to, takerAddress sdk.AccAddress
		toStr := msg.To
		if msg.Moderated && toStr == "" {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("ToAddress is missing."))
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
			if !msg.Moderated {
				if to.String() != from.String() {
					utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Cannot issue an asset. ReceiverAddress should be same as issuerAddress."))
					return
				}
			}
		}

		if msg.TakerAddress != "" {
			takerAddress, err = sdk.AccAddressFromBech32(msg.TakerAddress)
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

		if msg.Moderated {
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
			AssetQuantity: msg.AssetQuantity,
			AssetType:     msg.AssetType,
			AssetPrice:    msg.AssetPrice,
			DocumentHash:  msg.DocumentHash,
			QuantityUnit:  msg.QuantityUnit,
			Moderated:     msg.Moderated,
			TakerAddress:  takerAddress,
		}
		issueAssetMsg := client.BuildIssueAssetMsg(from, to, assetPeg)

		if kafka == true {
			ticketID := rest.TicketIDGenerator("BKIA")

			jsonResponse := rest.SendToKafka(rest.NewKafkaMsgFromRest(issueAssetMsg, ticketID, txCtx, cliCtx, msg.Password), kafkaState, cdc)
			w.WriteHeader(http.StatusAccepted)
			w.Write(jsonResponse)
		} else {
			output, err := utils.SendTxWithResponse(txCtx, cliCtx, []sdk.Msg{issueAssetMsg}, msg.Password)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			w.Write(utils.ResponseBytesToJSON(output))
		}
	}
}
