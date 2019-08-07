package rest

import (
	"fmt"
	"net/http"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	
	"github.com/commitHub/commitBlockchain/types"
	
	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/modules/acl"
	negotiationTypes "github.com/commitHub/commitBlockchain/modules/negotiation/internal/types"
)

type confirmSellerBidReq struct {
	BaseReq            rest.BaseReq `json:"base_req"`
	To                 string       `json:"to" `
	Bid                int64        `json:"bid" `
	Time               int64        `json:"time"`
	PegHash            string       `json:"pegHash"`
	SellerContractHash string       `json:"sellerContractHash"`
	Password           string       `json:"password"`
	Mode               string       `json:"mode"`
}

func ConfirmSellerBidRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		var req confirmSellerBidReq
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
		
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryACLAccount", fromAddr), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}
		
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("Unauthorized transaction. ACL is not defined for buyer"))
			return
		}
		
		var account acl.ACLAccount
		cliCtx.Codec.MustUnmarshalJSON(res, &account)
		
		if !account.GetACL().ConfirmSellerBid {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}
		
		to, err := cTypes.AccAddressFromBech32(req.To)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		pegHashHex, err := types.GetAssetPegHashHex(req.PegHash)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		negotiationID := negotiationTypes.NegotiationID(append(append(to.Bytes(), fromAddr.Bytes()...), pegHashHex.Bytes()...))
		kb, err := keys.NewKeyBaseFromHomeFlag()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		SignBytes := negotiationTypes.NewSignNegotiationBody(to, fromAddr, pegHashHex, req.Bid, req.Time)
		signature, _, err := kb.Sign(name, req.Password, SignBytes.GetSignBytes())
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		proposedNegotiation := &negotiationTypes.BaseNegotiation{
			NegotiationID:      negotiationID,
			BuyerAddress:       to,
			SellerAddress:      fromAddr,
			PegHash:            pegHashHex,
			Bid:                req.Bid,
			Time:               req.Time,
			SellerContractHash: req.SellerContractHash,
			BuyerSignature:     nil,
			SellerSignature:    signature,
		}
		
		msg := negotiationTypes.BuildMsgConfirmSellerBid(proposedNegotiation)
		rest2.SignAndBroadcast(w, req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}
}
