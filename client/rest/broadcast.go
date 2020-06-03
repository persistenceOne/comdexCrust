package rest

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/commitHub/commitBlockchain/codec"

	"github.com/commitHub/commitBlockchain/modules/auth"
)

func BroadcastRest(cliCtx context.CLIContext, cdc *codec.Codec, stdTx auth.StdTx, mode string) ([]byte, cTypes.Error) {

	txBytes, err := cdc.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		return nil, cTypes.NewError(DefaultCodeSpace, http.StatusInternalServerError, err.Error())
	}
	cliCtx = cliCtx.WithBroadcastMode(mode)

	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, cTypes.NewError(DefaultCodeSpace, http.StatusInternalServerError, err.Error())
	}

	return PostProcessResponse(cliCtx, res)
}
