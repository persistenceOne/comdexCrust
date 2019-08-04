package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/commitHub/commitBlockchain/modules/fiatFactory"
	"github.com/commitHub/commitBlockchain/types"
)

// QueryFiatHandlerFn :
func QueryFiatHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		peghashstr := vars["peghash"]
		nodeURI := "tcp://0.0.0.0:56657"

		cliCtx := cliCtx
		cliCtx = cliCtx.WithNodeURI(nodeURI)
		cliCtx = cliCtx.WithTrustNode(true)

		peghashHex, err := types.GetFiatPegHashHex(peghashstr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, _, err := cliCtx.QueryStore(fiatFactory.FiatPegHashStoreKey(peghashHex), fiatFactory.ModuleName)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("Couldn't query account. Error: %s", err.Error()))
			return
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var fiatPeg types.FiatPeg
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &fiatPeg)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't Unmarshall fiatPeg. Error: %s", err.Error()))
			return
		}

		output, err := cliCtx.Codec.MarshalJSON(fiatPeg)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("Couldn't marshall query result. Error: %s", err.Error()))
			return
		}

		w.Write(output)
	}

}
