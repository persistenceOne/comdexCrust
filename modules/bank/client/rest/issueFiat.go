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
	"github.com/persistenceOne/persistenceSDK/modules/bank/client"
	bankTypes "github.com/persistenceOne/persistenceSDK/modules/bank/internal/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

type IssueFiatReq struct {
	BaseReq           rest.BaseReq `json:"base_req"`
	To                string       `json:"to" valid:"required~Enter the ToAddress,matches(^commit[a-z0-9]{39}$)~ToAddress is Invalid"`
	GasAdjustment     string       `json:"gasAdjustment"`
	TransactionID     string       `json:"transactionID" valid:"required~Enter the TransactionID,  matches(^[A-Z0-9]+$)~transactionID is Invalid,length(2|40)~TransactionID length should be 2 to 40"`
	TransactionAmount int64        `json:"transactionAmount" valid:"required~Enter the TransactionAmount,matches(^[1-9]{1}[0-9]*$)~Invalid TransactionAmount"`
	Password          string       `json:"password" valid:"required~Enter the Password"`
	Mode              string       `json:"mode"`
}

func IssueFiatHandlerFunction(cliCtx context.CLIContext, kafkaBool bool, kafkaState kafka.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IssueFiatReq
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

		res, _, err := cliCtx.QueryStore(acl.GetACLAccountKey(to), "acl")
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, "ACL Account"))
			return
		}

		if len(res) == 0 {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}

		var account acl.ACLAccount
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &account)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, "ACL Account"))
			return
		}

		if !account.GetACL().IssueFiat {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}

		zoneID := account.GetZoneID()
		zoneData, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryZone", zoneID), nil)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}
		zoneAcc := cTypes.AccAddress(string(zoneData))
		cliCtx.Codec.MustUnmarshalJSON(zoneData, &zoneAcc)
		if !reflect.DeepEqual(fromAddr, zoneAcc) {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}

		fiatPeg := types.BaseFiatPeg{

			TransactionID:     req.TransactionID,
			TransactionAmount: req.TransactionAmount,
		}

		msg := client.BuildIssueFiatMsg(fromAddr, to, &fiatPeg)

		if kafkaBool == true {
			ticketID := kafka.TicketIDGenerator("ISFI")
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
