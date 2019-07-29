package order

import (
	"testing"

	sdk "github.com/commitHub/commitBlockchain/types"
)

func Test_IterateOrders(t *testing.T) {
	// var orders []sdk.Order
	order := om.NewOrder(fromAddress, toAddress, peghash)
	order.SetAssetPegWallet(sdk.AssetPegWallet{assetPeg[0]})
	order.SetFiatPegWallet(fiatPegWallet)
	om.SetOrder(ctx, order)
	om.IterateOrders(ctx, func(order sdk.Order) bool { return false })
}
