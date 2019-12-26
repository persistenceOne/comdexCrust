package assetFactory

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, assetPeg := range data.AssetPegs {
		keeper.SetAssetPeg(ctx, assetPeg)
	}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	assetPegs := keeper.GetAssetPegs(ctx)

	return GenesisState{assetPegs}
}
