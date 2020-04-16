package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	rest2 "github.com/persistenceOne/comdexCrust/client/rest"
	aclTypes "github.com/persistenceOne/comdexCrust/modules/acl/internal/types"
	"github.com/persistenceOne/comdexCrust/types"
)

func GetACLRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		cliCtx := cliCtx

		addr, err := cTypes.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrAccAddressFromBech32(aclTypes.DefaultCodeSpace, bech32addr))
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", aclTypes.QuerierRoute, "queryACLAccount", addr), nil)
		if err != nil {
			rest2.WriteErrorResponse(w, types.ErrQuery(aclTypes.DefaultCodeSpace, "ACL Account"))
			return
		}

		if len(res) == 0 {
			rest2.WriteErrorResponse(w, types.ErrQueryResponseLengthZero(aclTypes.DefaultCodeSpace, "ACL Account"))
			return
		}

		var account aclTypes.ACLAccount
		cliCtx.Codec.MustUnmarshalJSON(res, &account)

		rest.PostProcessResponse(w, cliCtx, account)
	}
}
