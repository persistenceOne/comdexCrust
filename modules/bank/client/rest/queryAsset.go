package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	rest2 "github.com/persistenceOne/comdexCrust/client/rest"
	"github.com/persistenceOne/comdexCrust/modules/assetFactory"
	bankTypes "github.com/persistenceOne/comdexCrust/modules/bank/internal/types"
	"github.com/persistenceOne/comdexCrust/types"
)

func QueryAssetHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		peghashstr := vars["peghash"]
		nodeURI := "tcp://0.0.0.0:46657"

		cliCtx := cliCtx
		cliCtx = cliCtx.WithNodeURI(nodeURI)
		cliCtx = cliCtx.WithTrustNode(true)

		peghashHex, err := types.GetAssetPegHashHex(peghashstr)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrPegHashHex(bankTypes.DefaultCodespace, peghashstr))
			return
		}

		res, _, err := cliCtx.QueryStore(assetFactory.AssetPegHashStoreKey(peghashHex), "asset")
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, "asset"))
			return
		}

		if len(res) == 0 {
			rest2.WriteErrorResponse(w, types.ErrQueryResponseLengthZero(bankTypes.DefaultCodespace, "asset"))
			return
		}

		var assetPeg types.AssetPeg
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &assetPeg)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrUnmarshal(bankTypes.DefaultCodespace, "assetPeg"))
			return
		}

		output, err := cliCtx.Codec.MarshalJSON(assetPeg)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrMarshal(bankTypes.DefaultCodespace, "asset"))
			return
		}

		w.Write(output)

	}
}
