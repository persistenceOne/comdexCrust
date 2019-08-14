package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	fiatFactoryTypes "github.com/commitHub/commitBlockchain/modules/fiatFactory/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

type issueFiatReq struct {
	BaseReq           rest.BaseReq `json:"base_req"`
	To                string       `json:"to"`
	PegHash           string       `json:"pegHash"`
	TransactionID     string       `json:"transactionID"`
	TransactionAmount int64        `json:"transactionAmount"`
	Password          string       `json:"password"`
	Mode              string       `json:"mode"`
}

func IssueFiatHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req issueFiatReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, name, err := context.GetFromFields(req.BaseReq.From, false)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrFromName(fiatFactoryTypes.DefaultCodeSpace))
			return
		}

		cliCtx = cliCtx.WithFromAddress(fromAddr)
		cliCtx = cliCtx.WithFromName(name)

		toAddr, err := cTypes.AccAddressFromBech32(req.To)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(fiatFactoryTypes.DefaultCodeSpace, req.To))
			return
		}

		pegHashHex, err := types.GetFiatPegHashHex(req.PegHash)
		fiat := types.BaseFiatPeg{
			TransactionAmount: req.TransactionAmount,
			TransactionID:     req.TransactionID,
			PegHash:           pegHashHex,
		}

		msg := fiatFactoryTypes.BuildIssueFiatMsg(fromAddr, toAddr, &fiat)
		rest2.SignAndBroadcast(req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}

}
