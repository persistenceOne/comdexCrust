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

type IssueAssetReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	To            string       `json:"to"`
	DocumentHash  string       `json:"documentHash"`
	AssetType     string       `json:"assetType" `
	AssetPrice    int64        `json:"assetPrice" `
	QuantityUnit  string       `json:"quantityUnit" `
	AssetQuantity int64        `json:"assetQuantity" `
	Moderated     bool         `json:"moderated"`
	TakerAddress  string       `json:"takerAddress"`
	Password      string       `json:"password"`
	Mode          string       `json:"mode"`
}

func IssueAssetHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IssueAssetReq
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
		
		if req.Moderated && req.To == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("ToAddress is missing."))
			return
		}
		
		var to, takerAddress cTypes.AccAddress
		if req.To == "" {
			to = fromAddr
		} else {
			to, err = cTypes.AccAddressFromBech32(req.To)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			if !req.Moderated {
				if to.String() != fromAddr.String() {
					rest.WriteErrorResponse(w, http.StatusInternalServerError,
						fmt.Sprintf("Cannot issue an asset. ReceiverAddress should be same as issuerAddress."))
					
					return
				}
			}
		}
		
		if req.TakerAddress != "" {
			takerAddress, err = cTypes.AccAddressFromBech32(req.TakerAddress)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryACLAccount", to), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}
		
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}
		
		var account acl.ACLAccount
		cliCtx.Codec.MustUnmarshalJSON(res, &account)
		
		if !account.GetACL().IssueAsset {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}
		
		if req.Moderated {
			zoneID := account.GetZoneID()
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryZone", zoneID), nil)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError,
					fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
				return
			}
			var zoneAddr cTypes.AccAddress
			cliCtx.Codec.MustUnmarshalJSON(res, &zoneAddr)
			
			if !reflect.DeepEqual(fromAddr, zoneAddr) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("you are not authorized person. Only zones can issue fiats."))
				return
			}
		}
		assetPeg := &types.BaseAssetPeg{
			AssetQuantity: req.AssetQuantity,
			AssetType:     req.AssetType,
			AssetPrice:    req.AssetPrice,
			DocumentHash:  req.DocumentHash,
			QuantityUnit:  req.QuantityUnit,
			Moderated:     req.Moderated,
			TakerAddress:  takerAddress,
		}
		
		msg := client.BuildIssueAssetMsg(fromAddr, to, assetPeg)
		rest2.SignAndBroadcast(w, req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}
}
