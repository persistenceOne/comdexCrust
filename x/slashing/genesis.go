package slashing

import (
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/x/stake/types"
)

// InitGenesis initializes the keeper's address to pubkey map.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	for _, validator := range data.Validators {
		keeper.addPubkey(ctx, validator.GetPubKey())
	}
	return
}
