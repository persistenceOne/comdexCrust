package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/commitHub/commitBlockchain/codec"

	"github.com/commitHub/commitBlockchain/modules/auth"
)

func BroadcastRest(w http.ResponseWriter, cliCtx context.CLIContext, cdc *codec.Codec, stdTx auth.StdTx, mode string) {

	txBytes, err := cdc.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	cliCtx = cliCtx.WithBroadcastMode(mode)

	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	rest.PostProcessResponse(w, cliCtx, res)
}
