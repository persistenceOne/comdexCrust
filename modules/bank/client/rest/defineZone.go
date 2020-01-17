package rest

import (
	"net/http"

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

type DefineZoneReq struct {
	BaseReq  rest.BaseReq `json:"base_req" `
	To       string       `json:"to" valid:"required~Enter the toAddress,matches(^persist[a-z0-9]{39}$)~toAddress is Invalid"`
	ZoneID   string       `json:"zoneID" valid:"required~Enter the zoneID, matches(^[A-Fa-f0-9]+$)~Invalid zoneID,length(2|40)~ZoneID length should be 2 to 40"`
	Password string       `json:"password" valid:"required~Enter the password"`
	Mode     string       `json:"mode"`
}

func DefineZoneHandler(cliCtx context.CLIContext, kafkaBool bool, kafkaState kafka.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DefineZoneReq
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

		zoneID, err := acl.GetZoneIDFromString(req.ZoneID)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrZoneIDFromString(bankTypes.DefaultCodespace, req.ZoneID))
			return
		}

		msg := bankTypes.BuildMsgDefineZone(fromAddr, to, zoneID)

		if kafkaBool == true {
			ticketID := kafka.TicketIDGenerator("DEZO")
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
