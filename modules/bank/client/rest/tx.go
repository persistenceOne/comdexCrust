package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	"github.com/commitHub/commitBlockchain/modules/bank/internal/types"
)

type SendReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount   cTypes.Coins `json:"amount" yaml:"amount"`
	Password string       `json:"password"`
	Mode     string       `json:"mode"`
}

func SendRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := cTypes.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req SendReq
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

		msg := types.NewMsgSend(fromAddr, toAddr, req.Amount)

		rest2.SignAndBroadcast(w, req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}
}
