package rest

import (
	"net/http"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	
	"github.com/commitHub/commitBlockchain/types"
	
	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	fiatFactoryTypes "github.com/commitHub/commitBlockchain/modules/fiatFactory/internal/types"
)

type sendFiatReq struct {
	BaseReq      rest.BaseReq
	OwnerAddress string `json:"ownerAddress"`
	To           string `json:"to"`
	AssetPegHash string `json:"assetPegHash"`
	FiatPegHash  string `json:"fiatPegHash"`
	Amount       int64  `json:"amount"`
	Password     string `json:"password"`
	Mode         string `json:"mode"`
}

func SendFiatHandlerFunction(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		var req sendFiatReq
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
		
		ownerAddr, err := cTypes.AccAddressFromBech32(req.OwnerAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		assetPegHashHex, err := types.GetAssetPegHashHex(req.AssetPegHash)
		fiatPegHashHex, err := types.GetFiatPegHashHex(req.FiatPegHash)
		fiatPeg := types.BaseFiatPeg{
			PegHash:           fiatPegHashHex,
			TransactionAmount: req.Amount,
		}
		
		msg := fiatFactoryTypes.BuildSendFiatMsg(fromAddr, ownerAddr, toAddr, assetPegHashHex, types.FiatPegWallet{fiatPeg})
		rest2.SignAndBroadcast(w, req.BaseReq, cliCtx, req.Mode, req.Password, []cTypes.Msg{msg})
	}
}
