package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	assetFactoryTypes "github.com/commitHub/commitBlockchain/modules/assetFactory/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

func QueryAssetRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		pegHashStr := vars["peg-hash"]

		pegHashHex, err := types.GetAssetPegHashHex(pegHashStr)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrPegHashHex(assetFactoryTypes.DefaultCodeSpace, pegHashStr))
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", assetFactoryTypes.QuerierRoute,
			assetFactoryTypes.PegHashKey), pegHashHex)

		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(assetFactoryTypes.DefaultCodeSpace, "ACL Account"))
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
