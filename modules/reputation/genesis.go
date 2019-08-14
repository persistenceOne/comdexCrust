package reputation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, reputation := range data.Reputations {
		keeper.SetAccountReputation(ctx, reputation)
	}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {

	reputations := keeper.GetReputations(ctx)

	return GenesisState{Reputations: reputations}
}
