package rest

import (
	"bytes"
	"fmt"
	"github.com/persistenceOne/comdexCrust/modules/acl"
	bankTypes "github.com/persistenceOne/comdexCrust/modules/bank/internal/types"

	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/persistenceOne/comdexCrust/client/rest"
	"github.com/persistenceOne/comdexCrust/kafka"
	"github.com/persistenceOne/comdexCrust/types"
)

type DefineACLReq struct {
	BaseReq            rest.BaseReq `json:"base_req"`
	ACLAddress         string       `json:"aclAddress" valid:"required~Enter the aclAddress,matches(^persist[a-z0-9]{39}$)~aclAddress is Invalid"`
	OrganizationID     string       `json:"organizationID" valid:"required~Enter the organizationID, matches(^[A-Fa-f0-9]+$)~Invalid organizationID,length(2|40)~OrganizationID length should be 2 to 40"`
	ZoneID             string       `json:"zoneID" valid:"required~Enter the zoneID, matches(^[A-Fa-f0-9]+$)~Invalid zoneID,length(2|40)~ZoneID length should be 2 to 40"`
	IssueAsset         string       `json:"issueAsset" valid:"required~Enter the issueAsset, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid issueAsset"`
	IssueFiat          string       `json:"issueFiat" valid:"required~Enter the issueFiat, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid issueFiat"`
	SendAsset          string       `json:"sendAsset" valid:"required~Enter the sendAsset, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid sendAsset"`
	SendFiat           string       `json:"sendFiat" valid:"required~Enter the sendFiat, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid issueFiat"`
	BuyerExecuteOrder  string       `json:"buyerExecuteOrder" valid:"required~Enter the buyerExecuteOrder, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid buyerExecuteOrder"`
	SellerExecuteOrder string       `json:"sellerExecuteOrder" valid:"required~Enter the issueAsset, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid issueAsset"`
	ChangeBuyerBid     string       `json:"changeBuyerBid" valid:"required~Enter the changeBuyerBid, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid changeBuyerBid"`
	ChangeSellerBid    string       `json:"changeSellerBid" valid:"required~Enter the changeSellerBid, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid changeSellerBid"`
	ConfirmBuyerBid    string       `json:"confirmBuyerBid" valid:"required~Enter the confirmBuyerBid, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid issueAsset"`
	ConfirmSellerBid   string       `json:"confirmSellerBid" valid:"required~Enter the confirmSellerBid, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid confirmSellerBid"`
	Negotiation        string       `json:"negotiation" valid:"required~Enter the negotiation, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid negotiation"`
	RedeemAsset        string       `json:"redeemAsset" valid:"required~Enter the redeemAsset, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid redeemAsset"`
	RedeemFiat         string       `json:"redeemFiat" valid:"required~Enter the redeemFiat, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid redeemFiat"`
	ReleaseAsset       string       `json:"releaseAsset" valid:"required~Enter the releaseAsset, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid releaseAsset"`
	Password           string       `json:"password" valid:"required~Enter the password"`
	Mode               string       `json:"mode"`
}

func DefineACLHandler(cliCtx context.CLIContext, kafkaBool bool, kafkaState kafka.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req DefineACLReq
		cliCtx := cliCtx

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			rest2.WriteErrorResponse(w, cTypes.NewError(acl.DefaultCodeSpace, http.StatusBadRequest, err.Error()))
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, name, err := context.GetFromFields(req.BaseReq.From, false)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrFromName(acl.DefaultCodeSpace))
			return
		}

		cliCtx = cliCtx.WithFromAddress(fromAddr)
		cliCtx = cliCtx.WithFromName(name)

		zoneID, err := acl.GetZoneIDFromString(req.ZoneID)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrZoneIDFromString(acl.DefaultCodeSpace, req.ZoneID))
			return
		}

		to, err := cTypes.AccAddressFromBech32(req.ACLAddress)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(acl.DefaultCodeSpace, req.ACLAddress))
			return
		}

		ACLReq := GetReqACLTxns(req)

		organizationID, err := acl.GetOrganizationIDFromString(req.OrganizationID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		Bytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryOrganization",
			organizationID), nil)
		if Bytes == nil {
			rest2.WriteErrorResponse(w, types.ErrQueryResponseLengthZero(acl.DefaultCodeSpace, "organization"))
			return
		}

		var rawMsg acl.Organization
		cliCtx.Codec.MustUnmarshalJSON(Bytes, &rawMsg)

		if bytes.Compare(rawMsg.ZoneID, zoneID) != 0 {
			rest2.WriteErrorResponse(w, types.ErrInvalidOrganizationWithZone(acl.DefaultCodeSpace))
			return
		}

		aclAccount := &acl.BaseACLAccount{
			Address:        to,
			ZoneID:         zoneID,
			OrganizationID: organizationID,
			ACL:            ACLReq,
		}
		msg := bankTypes.BuildMsgDefineACL(fromAddr, to, aclAccount)

		if kafkaBool == true {
			ticketID := kafka.TicketIDGenerator("DEAC")
			jsonResponse := kafka.SendToKafka(kafka.NewKafkaMsgFromRest(msg, ticketID, req.BaseReq, cliCtx, req.Mode, req.Password), kafkaState, cliCtx.Codec)
			w.WriteHeader(http.StatusAccepted)
			w.Write(jsonResponse)
		} else {
			output, err := rest2.SignAndBroadcast(req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
			if err != nil {
				rest2.WriteErrorResponse(w, err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(output)
		}
	}
}

func GetReqACLTxns(DefineACL DefineACLReq) acl.ACL {
	var Request acl.ACL
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
