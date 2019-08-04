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
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithFromAddress(fromAddr)
		cliCtx = cliCtx.WithFromName(name)

		toAddr, err := cTypes.AccAddressFromBech32(req.To)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ownerAddr, err := cTypes.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		assetPegHashHex, err := types.GetAssetPegHashHex(req.PegHash)

		msg := assetFactoryTypes.BuildExecuteAssetMsg(fromAddr, ownerAddr, toAddr, assetPegHashHex)
		rest2.SignAndBroadcast(w, req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}

}
