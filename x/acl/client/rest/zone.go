package rest

import (
	"net/http"

	"github.com/commitHub/commitBlockchain/x/acl"

	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/gorilla/mux"
)

//GetZoneRequestHandler : query zone account address Handler
func GetZoneRequestHandler(storeName string, r *mux.Router, cdc *wire.Codec, cliContext context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		strZoneID := vars["zoneID"]
		cliCtx := cliContext

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
