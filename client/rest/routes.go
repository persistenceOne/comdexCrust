package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

// RegisterRoutes : resgister REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/keys", QueryKeysRequestHandler).Methods("GET")
	r.HandleFunc("/keys", AddNewKeyRequestHandler).Methods("POST")
	r.HandleFunc("/keys/mnemonic", QueryMnemonicRequestHandler).Methods("GET")
	r.HandleFunc("/updatePassword/{name}", UpdateKeyRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/forgotPassword/{name}", ForgotPasswordRequestHandler(cliCtx)).Methods("POST")
}
