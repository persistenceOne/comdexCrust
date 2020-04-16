package rest

import (
	"fmt"
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

type SellerExecuteOrderReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	BuyerAddress  string       `json:"buyerAddress" valid:"required~Enter the BuyerAddress,matches(^persist[a-z0-9]{39}$)~BuyerAddress is Invalid"`
	SellerAddress string       `json:"sellerAddress" valid:"required~Enter the SellerAddress,matches(^persist[a-z0-9]{39}$)~SellerAddress is Invalid"`
	AWBProofHash  string       `json:"awbProofHash" valid:"required~Mandatory parameter awbProofHash missing,matches(^.*$)~Invalid awbProofHash,length(1|1000)~awbProofHash length should be 1 to 1000"`
	PegHash       string       `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	Password      string       `json:"password" valid:"required~Enter the Password"`
	Mode          string       `json:"mode"`
}

func SellerExecuteOrderRequestHandlerFn(cliCtx context.CLIContext, kafkaBool bool, kafkaState kafka.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SellerExecuteOrderReq
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

		sellerAddress, err := cTypes.AccAddressFromBech32(req.SellerAddress)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(bankTypes.DefaultCodespace, req.SellerAddress))
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryACLAccount",
			sellerAddress), nil)

		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, "ACL Account"))
			return
		}

		if len(res) == 0 {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}

		var account acl.ACLAccount
		cliCtx.Codec.MustUnmarshalJSON(res, &account)

		zoneID := account.GetZoneID()
		if zoneID == nil {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}
		zoneData, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryZone",
			zoneID), nil)

		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, "zone"))
			return
		}

		var zoneAddress cTypes.AccAddress
		cliCtx.Codec.MustUnmarshalJSON(zoneData, &zoneAddress)

		if zoneAddress.String() != fromAddr.String() && fromAddr.String() != sellerAddress.String() {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}

		if !account.GetACL().SellerExecuteOrder {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}

		buyerAddress, err := cTypes.AccAddressFromBech32(req.BuyerAddress)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(bankTypes.DefaultCodespace, req.BuyerAddress))
			return
		}

		pegHashHex, err := types.GetAssetPegHashHex(req.PegHash)
		msg := client.BuildSellerExecuteOrderMsg(fromAddr, buyerAddress, sellerAddress, pegHashHex, req.AWBProofHash)

		if kafkaBool == true {
			ticketID := kafka.TicketIDGenerator("SEEO")
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
