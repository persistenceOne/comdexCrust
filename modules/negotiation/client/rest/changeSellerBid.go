package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/modules/acl"
	negotiationTypes "github.com/commitHub/commitBlockchain/modules/negotiation/internal/types"
	types2 "github.com/commitHub/commitBlockchain/types"
)

type changeSellerBidBody struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	To       string       `json:"to" `
	Bid      int64        `json:"bid" `
	Time     int64        `json:"time" `
	PegHash  string       `json:"pegHash"`
	Password string       `json:"password"`
	Mode     string       `json:"mode"`
}

func ChangeSellerBidRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req changeSellerBidBody
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
				fmt.Sprintf("Unauthorized transaction. ACL is not defined for seller"))
			return
		}

		var account acl.ACLAccount
		cliCtx.Codec.MustUnmarshalJSON(res, &account)
		if !account.GetACL().ChangeSellerBid {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}

		to, err := cTypes.AccAddressFromBech32(req.To)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		pegHashHex, err := types2.GetAssetPegHashHex(req.PegHash)
		negotiationID := negotiationTypes.NegotiationID(append(append(to.Bytes(), fromAddr.Bytes()...), pegHashHex...))

		proposedNegotiation := &negotiationTypes.BaseNegotiation{
			NegotiationID: negotiationID,
			BuyerAddress:  to,
			SellerAddress: fromAddr,
			PegHash:       pegHashHex,
			Bid:           req.Bid,
			Time:          req.Time,
		}

		msg := negotiationTypes.BuildMsgChangeSellerBid(proposedNegotiation)
		rest2.SignAndBroadcast(w, req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}
}
