package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	reputationTypes "github.com/commitHub/commitBlockchain/modules/reputation/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

func QueryReputationRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		cliCtx := cliCtx

		bech32addr := vars["address"]
		if bech32addr == "" {
			rest2.WriteErrorResponse(w, types.ErrEmptyRequestFields(reputationTypes.DefaultCodeSpace, "address"))
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", reputationTypes.QuerierRoute, "reputationQuery", bech32addr), nil)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(reputationTypes.DefaultCodeSpace, "reputation"))
			return
		}

		if len(res) == 0 {
			rest2.WriteErrorResponse(w, types.ErrQueryResponseLengthZero(reputationTypes.DefaultCodeSpace, "reputation"))
			return
		}

		var reputation reputationTypes.AccountReputation
		cliCtx.Codec.MustUnmarshalJSON(res, &reputation)

		rest.PostProcessResponse(w, cliCtx, reputation)
	}
}
