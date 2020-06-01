package rest

import (
	"github.com/gorilla/mux"
)

// RegisterRoutes : resgister REST routes
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/keys", QueryKeysRequestHandler).Methods("GET")
	r.HandleFunc("/keys", AddNewKeyRequestHandler).Methods("POST")
}
