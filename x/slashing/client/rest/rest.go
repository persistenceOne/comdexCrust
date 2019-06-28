package rest

import (
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/crypto/keys"
	"github.com/comdex-blockchain/wire"
	
	"github.com/gorilla/mux"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	registerQueryRoutes(cliCtx, r, cdc)
	registerTxRoutes(cliCtx, r, cdc, kb)
}
