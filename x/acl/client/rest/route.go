package rest

import (
	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/gorilla/mux"
)

//RegisterRoutes : ServeCommand will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string) {
	r.HandleFunc("/zone/{zoneID}", GetZoneRequestHandler(storeName, r, cdc, ctx)).Methods("GET")
	r.HandleFunc("/organization/{organizationID}", GetOrganizationRequestHandler(storeName, r, cdc, ctx)).Methods("GET")
	r.HandleFunc("/acl/{address}", GetACLRequestHandler(storeName, r, cdc, ctx)).Methods("GET")
}
