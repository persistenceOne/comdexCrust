package rest

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/commitHub/commitBlockchain/kafka"
	bankTypes "github.com/commitHub/commitBlockchain/modules/bank/internal/types"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/commitHub/commitBlockchain/types"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/bank/client"
)

type ReleaseAssetReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	To       string       `json:"to" valid:"required~Enter the ToAddress,matches(^commit[a-z0-9]{39}$)~ToAddress is Invalid"`
	PegHash  string       `json:"pegHash" valid:"required~Enter the PegHash,matches(^[0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	Password string       `json:"password" valid:"required~Enter the Password"`
	Mode     string       `json:"mode"`
}

func ReleaseAssetHandlerFunction(cliCtx context.CLIContext, kafkaBool bool, kafkaState kafka.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req ReleaseAssetReq
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

		zoneID := account.GetZoneID()
		zoneData, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryZone", zoneID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}

		var zoneAddress cTypes.AccAddress
		cliCtx.Codec.MustUnmarshalJSON(zoneData, &zoneAddress)

		if zoneID == nil || zoneAddress.String() != fromAddr.String() {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}

		pegHashHex, err := types.GetAssetPegHashHex(req.PegHash)
		msg := client.BuildReleaseAssetMsg(fromAddr, to, pegHashHex)

		if kafkaBool == true {
			ticketID := kafka.TicketIDGenerator("RLAS")
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
