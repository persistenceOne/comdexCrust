package rest

import (
	"github.com/asaskevich/govalidator"
	"github.com/commitHub/commitBlockchain/kafka"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/bank/internal/types"
)

type DefineZoneReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	To       string       `json:"to" valid:"required~Enter the toAddress,matches(^commit[a-z0-9]{39}$)~toAddress is Invalid"`
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
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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

		zoneID, err := acl.GetZoneIDFromString(req.ZoneID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		msg := types.BuildMsgDefineZone(fromAddr, to, zoneID)

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
