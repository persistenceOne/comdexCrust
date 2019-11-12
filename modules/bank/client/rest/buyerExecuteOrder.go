package rest

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/kafka"
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/bank/client"
	bankTypes "github.com/commitHub/commitBlockchain/modules/bank/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

type BuyerExecuteOrderReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	BuyerAddress  string       `json:"buyerAddress" valid:"required~Enter the BuyerAddress,matches(^commit[a-z0-9]{39}$)~BuyerAddress is Invalid"`
	SellerAddress string       `json:"sellerAddress" valid:"required~Enter the BuyerAddress,matches(^commit[a-z0-9]{39}$)~BuyerAddress is Invalid"`
	PegHash       string       `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	FiatProofHash string       `json:"fiatProofHash" valid:"required~Mandatory parameter FiatProofHash missing,matches(^.*$)~Invalid FiatProofHash,length(1|1000)~FiatProofHash length should be 1 to 1000"`
	Password      string       `json:"password" valid:"required~Enter the Password"`
	Mode          string       `json:"mode"`
}

func BuyerExecuteOrderRequestHandlerFn(cliCtx context.CLIContext, kafkaBool bool, kafkaState kafka.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BuyerExecuteOrderReq

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

		buyerAddress, err := cTypes.AccAddressFromBech32(req.BuyerAddress)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(bankTypes.DefaultCodespace, req.BuyerAddress))
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryACLAccount", buyerAddress), nil)
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
		zoneData, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryZone", zoneID), nil)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, "zone"))
			return
		}

		var zoneAddress cTypes.AccAddress
		cliCtx.Codec.MustUnmarshalJSON(zoneData, &zoneAddress)

		if zoneAddress.String() != fromAddr.String() && fromAddr.String() != buyerAddress.String() {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}

		if !account.GetACL().BuyerExecuteOrder {
			rest2.WriteErrorResponse(w, types.ErrUnAuthorizedTransaction(bankTypes.DefaultCodespace))
			return
		}

		sellerAddress, err := cTypes.AccAddressFromBech32(req.SellerAddress)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(bankTypes.DefaultCodespace, req.SellerAddress))
			return
		}

		pegHashHex, err := types.GetAssetPegHashHex(req.PegHash)
		msg := client.BuildBuyerExecuteOrderMsg(fromAddr, buyerAddress, sellerAddress, pegHashHex, req.FiatProofHash)

		if kafkaBool == true {
			ticketID := kafka.TicketIDGenerator("BUEO")
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
