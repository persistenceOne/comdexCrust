package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/commitHub/commitBlockchain/modules/acl/internal/types"
)

// GetOrganizationRequestHandler query organization account address Handler
func GetOrganizationRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		strOrganizationID := vars["organizationID"]
		cliCtx := cliCtx

		if strOrganizationID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		organizationID, err := types.GetOrganizationIDFromString(vars["organizationID"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, "queryOrganization",
			organizationID), nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		fmt.Println(string(res))
		if res == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// var org types.Organization
		// cliCtx.Codec.MustUnmarshalJSON(res, &org)

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
