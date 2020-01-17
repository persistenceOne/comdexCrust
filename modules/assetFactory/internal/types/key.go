package types

import (
	"github.com/persistenceOne/persistenceSDK/types"
)

const (
	ModuleName   = "asset"
	StoreKey     = ModuleName
	RouterKey    = StoreKey
	QuerierRoute = RouterKey
)

var (
	PegHashKey = []byte{0x01}
)

// AssetPegHashStoreKey : converts peg hash to keystore key
func AssetPegHashStoreKey(assetPegHash types.PegHash) []byte {
	return append(PegHashKey, assetPegHash.Bytes()...)
}
