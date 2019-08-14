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

type redeemAssetBody struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Owner    string       `json:"owner"`
	To       string       `json:"to" `
	PegHash  string       `json:"pegHash"`
	Password string       `json:"password"`
	Mode     string       `json:"mode"`
}

func RedeemAssetHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req redeemAssetBody
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

		ownerAddr, err := cTypes.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(assetFactoryTypes.DefaultCodeSpace, req.Owner))
			return
		}

		toAddr, err := cTypes.AccAddressFromBech32(req.To)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(assetFactoryTypes.DefaultCodeSpace, req.To))
			return
		}

		pegHashHex, err := types.GetAssetPegHashHex(req.PegHash)
		msg := assetFactoryTypes.BuildRedeemAssetMsg(fromAddr, ownerAddr, toAddr, pegHashHex)
		_, _ = rest2.SignAndBroadcast(req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}

}
