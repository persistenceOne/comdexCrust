package orders

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, order := range data.Orders {
		keeper.SetOrder(ctx, order)
	}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	var orders []Order

	keeper.IterateOrders(ctx, func(order Order) (stop bool) {
		orders = append(orders, order)
		return false
	})

	return GenesisState{Orders: orders}
}
