package types

import (
	"github.com/persistenceOne/comdexCrust/types"
)

const (
	ModuleName   = "negotiation"
	StoreKey     = ModuleName
	RouterKey    = StoreKey
	QuerierRoute = RouterKey
)

var (
	NegotiationKey = []byte{0x01}
)

func GetNegotiationKey(id types.NegotiationID) types.NegotiationID {
	return append(NegotiationKey, id.Bytes()...)
}
