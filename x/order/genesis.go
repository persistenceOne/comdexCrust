package order

import (
	sdk "github.com/commitHub/commitBlockchain/types"
)

// InitOrder will initialize the order store
func InitOrder(ctx sdk.Context, keeper Keeper) (err error) {
	from := sdk.AccAddress([]byte("FromAddress"))
	to := sdk.AccAddress([]byte("ToAddress"))
	pegHash := sdk.PegHash([]byte("1"))

	order := keeper.om.NewOrder(from, to, pegHash)
	keeper.om.SetOrder(ctx, order)
	return
}
