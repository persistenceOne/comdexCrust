package rest

import (
	"fmt"
	"net/http"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	
	"github.com/commitHub/commitBlockchain/types"
	
	fiatFactoryTypes "github.com/commitHub/commitBlockchain/modules/fiatFactory/internal/types"
)

func QueryFiatRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		vars := mux.Vars(r)
		pegHashStr := vars["pegHash"]
		
		pegHashHex, err := types.GetFiatPegHashHex(pegHashStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", fiatFactoryTypes.QuerierRoute,
			fiatFactoryTypes.PegHashKey), pegHashHex)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError,
				fmt.Sprintf("couldn't query account. Error: %s", err.Error()))
			return
		}
		
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
