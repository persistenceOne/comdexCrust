package reputation

import sdk "github.com/comdex-blockchain/types"

// InitReputation will initialize the reputation store
func InitReputation(ctx sdk.Context, keeper Keeper) (err error) {
	
	accountReputation := DefaultAccountReputation()
	keeper.fm.SetAccountReputation(ctx, accountReputation)
	
	return
}

func DefaultAccountReputation() sdk.AccountReputation {
	return &sdk.BaseAccountReputation{
		Address: sdk.AccAddress([]byte("Address")),
	}
}
