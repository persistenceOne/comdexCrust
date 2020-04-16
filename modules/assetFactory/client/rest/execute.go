package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/persistenceOne/comdexCrust/client/rest"
	assetFactoryTypes "github.com/persistenceOne/comdexCrust/modules/assetFactory/internal/types"
	"github.com/persistenceOne/comdexCrust/types"
)

type executeAssetReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Owner    string       `json:"owner"`
	To       string       `json:"to" `
	PegHash  string       `json:"pegHash"`
	Password string       `json:"passPhrase"`
	Mode     string       `json:"mode"`
}

func ExecuteAssetHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req executeAssetReq
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

		ownerAddr, err := cTypes.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(assetFactoryTypes.DefaultCodeSpace, req.Owner))
			return
		}

		assetPegHashHex, err := types.GetAssetPegHashHex(req.PegHash)

		msg := assetFactoryTypes.BuildExecuteAssetMsg(fromAddr, ownerAddr, toAddr, assetPegHashHex)
		rest2.SignAndBroadcast(req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}

}
