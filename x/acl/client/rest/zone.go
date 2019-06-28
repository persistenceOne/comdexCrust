package rest

import (
	"net/http"
	
	"github.com/comdex-blockchain/x/acl"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/gorilla/mux"
)

// GetZoneRequestHandler : query zone account address Handler
func GetZoneRequestHandler(storeName string, r *mux.Router, cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		strZoneID := vars["zoneID"]
		
		zoneID, err := types.GetZoneIDFromString(strZoneID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		res, err := cliCtx.QueryStore(acl.ZoneStoreKey(zoneID), storeName)
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
		output, err := wire.MarshalJSONIndent(cdc, types.AccAddress(res).String())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(output)
		return
	}
}
