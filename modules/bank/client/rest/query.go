package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	rest2 "github.com/persistenceOne/persistenceSDK/client/rest"
	bankTypes "github.com/persistenceOne/persistenceSDK/modules/bank/internal/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

func QueryBalancesRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := cTypes.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(bankTypes.DefaultCodespace, bech32addr))
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := bankTypes.NewQueryBalanceParams(addr)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, "balance params"))
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/bank/balances", bz)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(bankTypes.DefaultCodespace, "balance"))
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		if len(res) == 0 {
			rest.PostProcessResponse(w, cliCtx, cTypes.Coins{})
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}