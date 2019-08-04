package rest

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/bank/internal/types"
)

type DefineOrganizationReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	To             string       `json:"to" `
	OrganizationID string       `json:"organizationID" `
	ZoneID         string       `json:"zoneID"`
	Password       string       `json:"password"`
	Mode           string       `json:"mode"`
}

func DefineOrganizationHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DefineOrganizationReq
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
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		organizationID, err := acl.GetOrganizationIDFromString(req.OrganizationID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		zoneID, err := acl.GetZoneIDFromString(req.ZoneID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		zoneData, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryZone", zoneID), nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var zoneAddr cTypes.AccAddress
		cliCtx.Codec.MustUnmarshalJSON(zoneData, &zoneAddr)
		if !reflect.DeepEqual(fromAddr, zoneAddr) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("You are not authorized person. Only zone can define an Organization"))
			return
		}

		msg := types.BuildMsgDefineOrganization(fromAddr, to, organizationID, zoneID)
		rest2.SignAndBroadcast(w, req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}
}
