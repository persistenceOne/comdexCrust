package rest

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/persistenceOne/comdexCrust/client/rest"
	"github.com/persistenceOne/comdexCrust/kafka"
	"github.com/persistenceOne/comdexCrust/modules/acl"
	"github.com/persistenceOne/comdexCrust/modules/bank/client"
	bankTypes "github.com/persistenceOne/comdexCrust/modules/bank/internal/types"
	"github.com/persistenceOne/comdexCrust/types"
)

type RedeemAssetReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	To       string       `json:"to" valid:"required~Enter the ToAddress,matches(^persist[a-z0-9]{39}$)~ToAddress is Invalid"`
	PegHash  string       `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	Password string       `json:"password" valid:"required~Enter the Password"`
	Mode     string       `json:"mode"`
}

func RedeemAssetHandlerFunction(cliCtx context.CLIContext, kafkaBool bool, kafkaState kafka.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RedeemAssetReq
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

		issuerAddress, err := cTypes.AccAddressFromBech32(req.To)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(bankTypes.DefaultCodespace, req.To))
			return
		}

		res, _, err := cliCtx.QueryStore(acl.GetACLAccountKey(fromAddr), "acl")
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, "ACL Account"))
			return
		}

		if len(res) == 0 {
			rest2.WriteErrorResponse(w, types.ErrQueryResponseLengthZero(bankTypes.DefaultCodespace, "ACL Account"))
			return
		}

		var account acl.ACLAccount
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &account)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrUnmarshal(bankTypes.DefaultCodespace, "ACL Account"))
			return
		}

		if !account.GetACL().RedeemAsset {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}

		pegHash, err := types.GetAssetPegHashHex(req.PegHash)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrPegHashHex(bankTypes.DefaultCodespace, req.PegHash))
			return
		}

		msg := client.BuildRedeemAssetMsg(issuerAddress, fromAddr, pegHash)

		if kafkaBool == true {
			ticketID := kafka.TicketIDGenerator("RDAS")
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
