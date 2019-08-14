package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	assetFactoryTypes "github.com/commitHub/commitBlockchain/modules/assetFactory/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

type issueAssetReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	To            string       `json:"to"`
	DocumentHash  string       `json:"documentHash"`
	PegHash       string       `json:"pegHash"`
	AssetType     string       `json:"assetType"`
	AssetPrice    int64        `json:"assetPrice"`
	QuantityUnit  string       `json:"quantityUnit"`
	AssetQuantity int64        `json:"assetQuantity"`
	Moderated     bool         `json:"moderated"`
	TakerAddress  string       `json:"takerAddress"`
	Password      string       `json:"password"`
	Mode          string       `json:"mode"`
}

func IssueAssetHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueAssetReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, name, err := context.GetFromFields(req.BaseReq.From, false)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrFromName(assetFactoryTypes.DefaultCodeSpace))
			return
		}

		cliCtx = cliCtx.WithFromAddress(fromAddr)
		cliCtx = cliCtx.WithFromName(name)

		toAddr, err := cTypes.AccAddressFromBech32(req.To)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(assetFactoryTypes.DefaultCodeSpace, req.To))
			return
		}

		pegHashHex, err := types.GetAssetPegHashHex(req.PegHash)
		asset := types.BaseAssetPeg{
			AssetQuantity: req.AssetQuantity,
			AssetType:     req.AssetType,
			AssetPrice:    req.AssetPrice,
			DocumentHash:  req.DocumentHash,
			QuantityUnit:  req.QuantityUnit,
			PegHash:       pegHashHex,
		}

		msg := assetFactoryTypes.BuildIssueAssetMsg(fromAddr, toAddr, &asset)
		rest2.SignAndBroadcast(req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}
}
