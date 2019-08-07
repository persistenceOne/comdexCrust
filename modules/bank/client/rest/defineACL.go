package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	
	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/modules/acl"
	bankTypes "github.com/commitHub/commitBlockchain/modules/bank/internal/types"
)

type DefineACLReq struct {
	BaseReq            rest.BaseReq `json:"base_req"`
	ACLAddress         string       `json:"aclAddress"`
	OrganizationID     string       `json:"organizationID" `
	ZoneID             string       `json:"zoneID"`
	IssueAsset         string       `json:"issueAsset" `
	IssueFiat          string       `json:"issueFiat" `
	SendAsset          string       `json:"sendAsset" `
	SendFiat           string       `json:"sendFiat"  `
	BuyerExecuteOrder  string       `json:"buyerExecuteOrder" `
	SellerExecuteOrder string       `json:"sellerExecuteOrder" `
	ChangeBuyerBid     string       `json:"changeBuyerBid" `
	ChangeSellerBid    string       `json:"changeSellerBid" `
	ConfirmBuyerBid    string       `json:"confirmBuyerBid" `
	ConfirmSellerBid   string       `json:"confirmSellerBid" `
	Negotiation        string       `json:"negotiation" `
	RedeemAsset        string       `json:"redeemAsset" `
	RedeemFiat         string       `json:"redeemFiat" `
	ReleaseAsset       string       `json:"releaseAsset"`
	Password           string       `json:"password"`
	Mode               string       `json:"mode"`
}

func DefineACLHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DefineACLReq
		cliCtx := cliCtx
		
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		
		fromAddr, name, err := context.GetFromFields(req.BaseReq.From, false)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		cliCtx = cliCtx.WithFromAddress(fromAddr)
		cliCtx = cliCtx.WithFromName(name)
		
		zoneID, err := acl.GetZoneIDFromString(req.ZoneID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		
		to, err := cTypes.AccAddressFromBech32(req.ACLAddress)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Converting Address to bech32 is faild"))
			return
		}
		
		ACLReq := GetReqACLTxns(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		organizationID, err := acl.GetOrganizationIDFromString(req.OrganizationID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		
		Bytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryOrganization",
			organizationID), nil)
		if Bytes == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("organization is not defined"))
			return
		}
		
		var rawMsg acl.Organization
		cliCtx.Codec.MustUnmarshalJSON(Bytes, &rawMsg)
		
		if bytes.Compare(rawMsg.ZoneID, zoneID) != 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("organization does not belongs to respected zone"))
			return
		}
		
		aclAccount := &acl.BaseACLAccount{
			Address:        to,
			ZoneID:         zoneID,
			OrganizationID: organizationID,
			ACL:            ACLReq,
		}
		
		msg := bankTypes.BuildMsgDefineACL(fromAddr, to, aclAccount)
		rest2.SignAndBroadcast(w, req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
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
