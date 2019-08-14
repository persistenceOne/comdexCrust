package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	fiatFactoryTypes "github.com/commitHub/commitBlockchain/modules/fiatFactory/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

func QueryFiatRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		pegHashStr := vars["pegHash"]

		pegHashHex, err := types.GetFiatPegHashHex(pegHashStr)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrPegHashHex(fiatFactoryTypes.DefaultCodeSpace, pegHashStr))
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", fiatFactoryTypes.QuerierRoute,
			fiatFactoryTypes.PegHashKey), pegHashHex)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(fiatFactoryTypes.DefaultCodeSpace, pegHashStr))
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
