package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/zone/{zoneID}", GetZoneRequestHandler(ctx)).Methods("GET")
	r.HandleFunc("/organization/{organizationID}", GetOrganizationRequestHandler(ctx)).Methods("GET")
	r.HandleFunc("/acl/{address}", GetACLRequestHandler(ctx)).Methods("GET")
}
