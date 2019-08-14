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

func GetZoneRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strZoneID := vars["zoneID"]
		cliCtx := cliCtx

		zoneID, err := aclTypes.GetZoneIDFromString(strZoneID)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrZoneIDFromString(aclTypes.DefaultCodeSpace, strZoneID))
			return
		}
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", aclTypes.QuerierRoute, "queryZone", zoneID), nil)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(aclTypes.DefaultCodeSpace, "zone"))
			return
		}

		if res == nil {
			rest2.WriteErrorResponse(w, types.ErrQueryResponseLengthZero(aclTypes.DefaultCodeSpace, "zone"))
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
