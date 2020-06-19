package rest

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/commitHub/commitBlockchain/kafka"
	bankTypes "github.com/commitHub/commitBlockchain/modules/bank/internal/types"
	"net/http"
	"reflect"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/commitHub/commitBlockchain/types"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/bank/client"
)

type IssueAssetReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	To            string       `json:"to" valid:"matches(^commit[a-z0-9]{39}$)~ToAddress is Invalid"`
	DocumentHash  string       `json:"documentHash" valid:"required~Enter the DocumentHash,matches(^.*$)~Invalid DocumentHash,length(1|1000)~DocumentHash length should be 1 to 1000"`
	AssetType     string       `json:"assetType" valid:"required~Enter the assetType,matches(^[A-Za-z1-9 _-]*$)~Invalid AssetType"`
	AssetPrice    int64        `json:"assetPrice" valid:"required~Enter the assetPrice,matches(^[1-9]{1}[0-9]*$)~Invalid assetPrice"`
	QuantityUnit  string       `json:"quantityUnit" valid:"required~Enter the QuantityUnit,matches(^[A-Za-z]*$)~Invalid QuantityUnit"`
	AssetQuantity int64        `json:"assetQuantity" valid:"required~Enter the AssetQuantity,matches(^[1-9]{1}[0-9]*$)~Invalid AssetQuantity"`
	Moderated     bool         `json:"moderated"`
	TakerAddress  string       `json:"takerAddress"`
	Password      string       `json:"password" valid:"required~Enter the Password"`
	Mode          string       `json:"mode"`
}

func IssueAssetHandlerFunction(cliCtx context.CLIContext, kafkaBool bool, kafkaState kafka.KafkaState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IssueAssetReq
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

		if req.Moderated && req.To == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("ToAddress is missing."))
			return
		}

		var to, takerAddress cTypes.AccAddress
		if req.To == "" {
			to = fromAddr
		} else {
			to, err = cTypes.AccAddressFromBech32(req.To)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			if !req.Moderated {
				if to.String() != fromAddr.String() {
					rest.WriteErrorResponse(w, http.StatusInternalServerError,
						fmt.Sprintf("Cannot issue an asset. ReceiverAddress should be same as issuerAddress."))

					return
				}
			}
		}

		if req.TakerAddress != "" {
			takerAddress, err = cTypes.AccAddressFromBech32(req.TakerAddress)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
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

		if !account.GetACL().IssueAsset {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Unauthorized transaction"))
			return
		}

		if req.Moderated {
			zoneID := account.GetZoneID()
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", acl.QuerierRoute, "queryZone", zoneID), nil)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError,
					fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
				return
			}
			var zoneAddr cTypes.AccAddress
			cliCtx.Codec.MustUnmarshalJSON(res, &zoneAddr)

			if !reflect.DeepEqual(fromAddr, zoneAddr) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("you are not authorized person. Only zones can issue fiats."))
				return
			}
		}
		assetPeg := &types.BaseAssetPeg{
			AssetQuantity: req.AssetQuantity,
			AssetType:     req.AssetType,
			AssetPrice:    req.AssetPrice,
			DocumentHash:  req.DocumentHash,
			QuantityUnit:  req.QuantityUnit,
			Moderated:     req.Moderated,
			TakerAddress:  takerAddress,
		}

		msg := client.BuildIssueAssetMsg(fromAddr, to, assetPeg)
		if kafkaBool == true {
			ticketID := kafka.TicketIDGenerator("ISAS")
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
