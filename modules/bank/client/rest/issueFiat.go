package rest

import (
	"fmt"
	"net/http"
	"reflect"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	
	"github.com/commitHub/commitBlockchain/types"
	
	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/bank/client"
)

type IssueFiatReq struct {
	BaseReq           rest.BaseReq `json:"base_req"`
	To                string       `json:"to"`
	GasAdjustment     string       `json:"gasAdjustment"`
	TransactionID     string       `json:"transactionID"`
	TransactionAmount int64        `json:"transactionAmount" `
	Password          string       `json:"password"`
	Mode              string       `json:"mode"`
}

func IssueFiatHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IssueFiatReq
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
		
		to, err := cTypes.AccAddressFromBech32(req.To)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		res, _, err := cliCtx.QueryStore(acl.GetACLAccountKey(to), "acl")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}
		
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("Unauthorized transaction. ACL is not defined for buyer."))
			return
		}
		
		var account acl.ACLAccount
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &account)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't parse query result. Result: %s. Error: %s", res, err.Error()))
			return
		}
		
		if !account.GetACL().IssueFiat {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}
		
		zoneID := account.GetZoneID()
		zoneData, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryZone", zoneID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}
		zoneAcc := cTypes.AccAddress(string(zoneData))
		cliCtx.Codec.MustUnmarshalJSON(zoneData, &zoneAcc)
		if !reflect.DeepEqual(fromAddr, zoneAcc) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("you are not authorized person. Only zones can issue fiats."))
			return
		}
		
		fiatPeg := types.BaseFiatPeg{
			
			TransactionID:     req.TransactionID,
			TransactionAmount: req.TransactionAmount,
		}
		
		msg := client.BuildIssueFiatMsg(fromAddr, to, &fiatPeg)
		rest2.SignAndBroadcast(w, req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}
	
}
