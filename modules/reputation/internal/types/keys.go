package types

const (
	ModuleName   = "reputations"
	StoreKey     = ModuleName
	RouterKey    = StoreKey
	QuerierRoute = RouterKey
)

var (
	ReputationKey = []byte{0x08}
)

func GetReputationKey(storeKey string) []byte {
	return append(ReputationKey, []byte(storeKey)...)
}
