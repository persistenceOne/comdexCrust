package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	rest2 "github.com/persistenceOne/comdexCrust/client/rest"
	bankTypes "github.com/persistenceOne/comdexCrust/modules/bank/internal/types"
	"github.com/persistenceOne/comdexCrust/modules/fiatFactory"
	"github.com/persistenceOne/comdexCrust/types"
)

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
			rest2.WriteErrorResponse(w, types.ErrGoValidator(bankTypes.DefaultCodespace))
			return
		}

		res, _, err := cliCtx.QueryStore(fiatFactory.FiatPegHashStoreKey(peghashHex), fiatFactory.ModuleName)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, peghashstr))
			return
		}

		if len(res) == 0 {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, peghashstr))
			return
		}

		var fiatPeg types.FiatPeg
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &fiatPeg)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrUnmarshal(bankTypes.DefaultCodespace, "fiatPeg"))
			return
		}

		output, err := cliCtx.Codec.MarshalJSON(fiatPeg)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrMarshal(bankTypes.DefaultCodespace, "fiatPeg"))
			return
		}

		w.Write(output)
	}

}
