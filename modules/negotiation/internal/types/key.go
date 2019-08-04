package types

const (
	ModuleName   = "negotiation"
	StoreKey     = ModuleName
	RouterKey    = StoreKey
	QuerierRoute = RouterKey
)

var (
	NegotiationKey = []byte{0x01}
)

func GetNegotiationKey(id NegotiationID) NegotiationID {
	return append(NegotiationKey, id.Bytes()...)
}
