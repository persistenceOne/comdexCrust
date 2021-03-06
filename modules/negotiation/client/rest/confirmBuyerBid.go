package rest

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/commitHub/commitBlockchain/kafka"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/commitHub/commitBlockchain/types"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/modules/acl"
	negotiationTypes "github.com/commitHub/commitBlockchain/modules/negotiation/internal/types"
)

type confirmBuyerBidReq struct {
	BaseReq           rest.BaseReq `json:"base_req"`
	To                string       `json:"to" valid:"required~Enter the ToAddress,matches(^commit[a-z0-9]{39}$)~ToAddress is Invalid"`
	Bid               int64        `json:"bid" valid:"required~Enter the Bid,matches(^[1-9]{1}[0-9]*$)~Enter valid Bid"`
	Time              int64        `json:"time" valid:"required~Enter the Time,matches(^[1-9]{1}[0-9]*$)~Enter valid Time"`
	PegHash           string       `json:"pegHash" valid:"required~Enter the PegHash,matches(^[0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	BuyerContractHash string       `json:"buyerContractHash" valid:"required~Enter the BuyerContractHash, matches(^.*$)~Invalid BuyerContractHash,length(1|1000)~BuyerContractHash length should be 1 to 1000"`
	Password          string       `json:"password" valid:"required~Enter the Password"`
	Mode              string       `json:"mode"`
}

func ConfirmBuyerBidRequestHandlerFn(cliCtx context.CLIContext, kafkaBool bool, kafkaState kafka.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req confirmBuyerBidReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		_, err := govalidator.ValidateStruct(req)
		if err != nil {
			rest2.WriteErrorResponse(w, cTypes.NewError(negotiationTypes.DefaultCodeSpace, http.StatusBadRequest, err.Error()))
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

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryACLAccount", fromAddr), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}

		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("Unauthorized transaction. ACL is not defined for Buyer"))
			return
		}

		var account acl.ACLAccount
		cliCtx.Codec.MustUnmarshalJSON(res, &account)

		if !account.GetACL().ConfirmBuyerBid {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}

		to, err := cTypes.AccAddressFromBech32(req.To)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		pegHashHex, err := types.GetAssetPegHashHex(req.PegHash)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		negotiationID := negotiationTypes.NegotiationID(append(append(fromAddr.Bytes(), to.Bytes()...), pegHashHex.Bytes()...))

		kb, err := keys.NewKeyBaseFromHomeFlag()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		SignBytes := negotiationTypes.NewSignNegotiationBody(fromAddr, to, pegHashHex, req.Bid, req.Time)
		signature, _, err := kb.Sign(name, req.Password, SignBytes.GetSignBytes())
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposedNegotiation := &negotiationTypes.BaseNegotiation{
			NegotiationID:     negotiationID,
			BuyerAddress:      fromAddr,
			SellerAddress:     to,
			PegHash:           pegHashHex,
			Bid:               req.Bid,
			Time:              req.Time,
			BuyerContractHash: req.BuyerContractHash,
			BuyerSignature:    signature,
			SellerSignature:   nil,
		}

		msg := negotiationTypes.BuildMsgConfirmBuyerBid(proposedNegotiation)

		if kafkaBool == true {
			ticketID := kafka.TicketIDGenerator("COBB")
			jsonResponse := kafka.SendToKafka(kafka.NewKafkaMsgFromRest(msg, ticketID, req.BaseReq, cliCtx, req.Mode, req.Password), kafkaState, cliCtx.Codec)
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write(jsonResponse)
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
