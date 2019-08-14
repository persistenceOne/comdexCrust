package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	rest2 "github.com/commitHub/commitBlockchain/client/rest"
	aclTypes "github.com/commitHub/commitBlockchain/modules/acl/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

// GetOrganizationRequestHandler query organization account address Handler
func GetOrganizationRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		strOrganizationID := vars["organizationID"]
		cliCtx := cliCtx

		if strOrganizationID == "" {
			rest2.WriteErrorResponse(w, types.ErrEmptyRequestFields(aclTypes.DefaultCodeSpace, strOrganizationID))
			return
		}

		organizationID, err := aclTypes.GetOrganizationIDFromString(vars["organizationID"])
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrOrganizationIDFromString(aclTypes.DefaultCodeSpace, strOrganizationID))
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", aclTypes.QuerierRoute, "queryOrganization",
			organizationID), nil)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(aclTypes.DefaultCodeSpace, "organization"))
			return
		}

		if res == nil {
			rest2.WriteErrorResponse(w, types.ErrQueryResponseLengthZero(aclTypes.DefaultCodeSpace, "organization"))
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
