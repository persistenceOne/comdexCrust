package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	
	"github.com/asaskevich/govalidator"
	
	cliclient "github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	"github.com/comdex-blockchain/rest"
	
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/acl"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	authctx "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/bank"
)

// DefineACLHandler : Rest hanler to define acl
func DefineACLHandler(cdc *wire.Codec, cliCtx context.CLIContext, kafka bool, kafkaState rest.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var msg acl.DefineACLBody
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
		zoneID, err := sdk.GetZoneIDFromString(msg.ZoneID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		cliCtx = cliCtx.WithFromAddressName(msg.From)
		from, err := cliCtx.GetFromAddress()
		if err != nil {
			sdk.ErrInvalidAddress("The given Address is Invalid")
			return
		}
		
		txCtx := authctx.TxContext{
			Codec:         cdc,
			AccountNumber: msg.AccountNumber,
			Sequence:      msg.Sequence,
			Gas:           msg.Gas,
			ChainID:       msg.ChainID,
		}
		to, err := sdk.AccAddressFromBech32(msg.ACLAddress)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Converting Address to bech32 is faild"))
			return
		}
		cliCtx = cliCtx.WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
		ACLReq := GetReqACLTxns(msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		organizationID, err := sdk.GetOrganizationIDFromString(msg.OrganizationID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		
		aclAccount := &sdk.BaseACLAccount{
			Address:        from,
			ZoneID:         zoneID,
			OrganizationID: organizationID,
			ACL:            ACLReq,
		}
		msgACL := bank.BuildMsgDefineACL(from, to, aclAccount)
		
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
		
		if kafka == true {
			ticketID := rest.TicketIDGenerator("ACDA")
			
			jsonResponse := rest.SendToKafka(rest.NewKafkaMsgFromRest(msgACL, ticketID, txCtx, cliCtx, msg.Password), kafkaState, cdc)
			w.WriteHeader(http.StatusAccepted)
			w.Write(jsonResponse)
		} else {
			output, err := utils.SendTxWithResponse(txCtx, cliCtx, []sdk.Msg{msgACL}, msg.Password)
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			
			w.Write(utils.ResponseBytesToJSON(output))
		}
	}
}

// GetReqACLTxns : build and returns the sdk.Request type object
func GetReqACLTxns(DefineACL acl.DefineACLBody) sdk.ACL {
	var Request sdk.ACL
	data, err := strconv.ParseBool(DefineACL.IssueAsset)
	if err == nil {
		Request.IssueAsset = data
	}
	data, err = strconv.ParseBool(DefineACL.IssueFiat)
	if err == nil {
		Request.IssueFiat = data
	}
	data, err = strconv.ParseBool(DefineACL.SendAsset)
	if err == nil {
		Request.SendAsset = data
	}
	data, err = strconv.ParseBool(DefineACL.SendFiat)
	if err == nil {
		Request.SendFiat = data
	}
	data, err = strconv.ParseBool(DefineACL.BuyerExecuteOrder)
	if err == nil {
		Request.BuyerExecuteOrder = data
	}
	data, err = strconv.ParseBool(DefineACL.SellerExecuteOrder)
	if err == nil {
		Request.SellerExecuteOrder = data
	}
	data, err = strconv.ParseBool(DefineACL.ChangeBuyerBid)
	if err == nil {
		Request.ChangeBuyerBid = data
	}
	data, err = strconv.ParseBool(DefineACL.ChangeSellerBid)
	if err == nil {
		Request.ChangeSellerBid = data
	}
	data, err = strconv.ParseBool(DefineACL.ConfirmBuyerBid)
	if err == nil {
		Request.ConfirmBuyerBid = data
	}
	data, err = strconv.ParseBool(DefineACL.ConfirmSellerBid)
	if err == nil {
		Request.ConfirmSellerBid = data
	}
	data, err = strconv.ParseBool(DefineACL.Negotiation)
	if err == nil {
		Request.Negotiation = data
	}
	data, err = strconv.ParseBool(DefineACL.RedeemFiat)
	if err == nil {
		Request.RedeemFiat = data
	}
	data, err = strconv.ParseBool(DefineACL.RedeemAsset)
	if err == nil {
		Request.RedeemAsset = data
	}
	data, err = strconv.ParseBool(DefineACL.ReleaseAsset)
	if err == nil {
		Request.ReleaseAsset = data
	}
	return Request
}
