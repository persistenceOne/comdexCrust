package rest

import (
	"fmt"
	"net/http"
	
	"github.com/comdex-blockchain/x/acl"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	"github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/gorilla/mux"
)

// GetOrganizationRequestHandler query organization account address Handler
func GetOrganizationRequestHandler(storeName string, r *mux.Router, cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		strOrganizationID := vars["organizationID"]
		if strOrganizationID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		organizationID, err := types.GetOrganizationIDFromString(strOrganizationID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		
		res, err := cliCtx.QueryStore(acl.OrganizationStoreKey(organizationID), storeName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if res == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		decoder := acl.GetOrganizationDecoder(cdc)
		bytes, err := decoder(res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		output, err := cdc.MarshalJSON(bytes)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("couldn't marshall query result. Error: %s", err.Error()))
			return
		}
		w.Write(output)
		return
	}
}
