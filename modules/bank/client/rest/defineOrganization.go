package rest

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/asaskevich/govalidator"
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/persistenceOne/persistenceSDK/client/rest"
	"github.com/persistenceOne/persistenceSDK/kafka"
	"github.com/persistenceOne/persistenceSDK/modules/acl"
	bankTypes "github.com/persistenceOne/persistenceSDK/modules/bank/internal/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

type DefineOrganizationReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	To             string       `json:"to" valid:"required~Enter the to Address,matches(^persist[a-z0-9]{39}$)~to Address is Invalid"`
	OrganizationID string       `json:"organizationID" valid:"required~Enter the organizationID, matches(^[A-Fa-f0-9]+$)~Invalid OrganizationID,length(2|40)~OrganizationID length should be 2 to 40"`
	ZoneID         string       `json:"zoneID" valid:"required~Enter the zoneID, matches(^[A-Fa-f0-9]+$)~Invalid zoneID,length(2|40)~ZoneID length should be 2 to 40"`
	Password       string       `json:"password" valid:"required~Enter the password"`
	Mode           string       `json:"mode"`
}

func DefineOrganizationHandler(cliCtx context.CLIContext, kafkaBool bool, kafkaState kafka.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DefineOrganizationReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			rest2.WriteErrorResponse(w, cTypes.NewError(bankTypes.DefaultCodespace, http.StatusBadRequest, err.Error()))
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, name, err := context.GetFromFields(req.BaseReq.From, false)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrFromName(bankTypes.DefaultCodespace))
			return
		}

		cliCtx = cliCtx.WithFromAddress(fromAddr)
		cliCtx = cliCtx.WithFromName(name)

		to, err := cTypes.AccAddressFromBech32(req.To)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(bankTypes.DefaultCodespace, req.To))
			return
		}

		organizationID, err := acl.GetOrganizationIDFromString(req.OrganizationID)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrOrganizationIDFromString(bankTypes.DefaultCodespace, req.OrganizationID))
			return
		}
		zoneID, err := acl.GetZoneIDFromString(req.ZoneID)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrZoneIDFromString(bankTypes.DefaultCodespace, req.ZoneID))
			return
		}
		zoneData, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryZone", zoneID), nil)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, req.ZoneID))
			return
		}

		var zoneAddr cTypes.AccAddress
		cliCtx.Codec.MustUnmarshalJSON(zoneData, &zoneAddr)
		if !reflect.DeepEqual(fromAddr, zoneAddr) {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}

		msg := bankTypes.BuildMsgDefineOrganization(fromAddr, to, organizationID, zoneID)

		if kafkaBool == true {
			ticketID := kafka.TicketIDGenerator("DEOR")
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
