package negotiation

import sdk "github.com/commitHub/commitBlockchain/types"

// InitNegotiation will initialize the negotiation Store
func InitNegotiation(ctx sdk.Context, keeper Keeper) (err error) {
	from := sdk.AccAddress([]byte("FromAddress"))
	to := sdk.AccAddress([]byte("ToAddress"))
	pegHash := sdk.PegHash([]byte("1"))

	negotiation := keeper.nm.NewNegotiation(from, to, pegHash)
	keeper.nm.SetNegotiation(ctx, negotiation)

	return nil
}
