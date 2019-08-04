package types

import (
	"github.com/commitHub/commitBlockchain/modules/negotiation"
)

const (
	ModuleName   = "orders"
	StoreKey     = ModuleName
	RouterKey    = StoreKey
	QuerierRoute = RouterKey
)

var (
	OrdersKey = []byte{0x07}
)

func GetOrderKey(negotiationID negotiation.NegotiationID) []byte {
	return append(OrdersKey, negotiationID.Bytes()...)
}
