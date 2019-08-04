package negotiation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, negotiation := range data.Negotiations {
		keeper.SetNegotiation(ctx, negotiation)
	}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	negotiations := keeper.GetNegotiations(ctx)

	return GenesisState{negotiations}
}
