package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/commitHub/commitBlockchain/modules/assetFactory"
	"github.com/commitHub/commitBlockchain/types"
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
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, _, err := cliCtx.QueryStore(assetFactory.AssetPegHashStoreKey(peghashHex), "asset")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var assetPeg types.AssetPeg
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &assetPeg)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't Unmarshall assetPeg. Error: %s", err.Error()))
			return
		}

		output, err := cliCtx.Codec.MarshalJSON(assetPeg)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't marshall query result. Error: %s", err.Error()))
			return
		}

		w.Write(output)

	}
}
