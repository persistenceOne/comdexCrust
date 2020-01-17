package types

import (
	"github.com/persistenceOne/persistenceSDK/types"
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

func GetOrderKey(negotiationID types.NegotiationID) []byte {
	return append(OrdersKey, negotiationID.Bytes()...)
}
