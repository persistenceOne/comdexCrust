package assetFactory

import (
	"github.com/comdex-blockchain/types"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// AssetPegMapper : encoder decoder for asset type
type AssetPegMapper struct {
	key   sdk.StoreKey
	proto func() sdk.AssetPeg
	cdc   *wire.Codec
}

// NewAssetPegMapper : returns asset mapper
func NewAssetPegMapper(cdc *wire.Codec, key sdk.StoreKey, proto func() sdk.AssetPeg) AssetPegMapper {
	return AssetPegMapper{
		key:   key,
		proto: proto,
		cdc:   cdc,
	}
}

// AssetPegHashStoreKey : converts peg hash to keystore key
func AssetPegHashStoreKey(assetPegHash sdk.PegHash) []byte {
	return append([]byte("PegHash:"), assetPegHash.Bytes()...)
}

// encode the AssetPeg inteface
func (am AssetPegMapper) encodeAssetPeg(asset sdk.AssetPeg) []byte {
	bz, err := am.cdc.MarshalBinaryBare(asset)
	if err != nil {
		panic(err)
	}
	return bz
}

// decode the AssetPeg interface
func (am AssetPegMapper) decodeAssetPeg(bz []byte) (asset sdk.AssetPeg) {
	err := am.cdc.UnmarshalBinaryBare(bz, &asset)
	if err != nil {
		panic(err)
	}
	return
}

// SetAssetPeg : set asset peg
func (am AssetPegMapper) SetAssetPeg(ctx sdk.Context, asset sdk.AssetPeg) {
	pegHash := asset.GetPegHash()
	store := ctx.KVStore(am.key)
	bz := am.encodeAssetPeg(asset)
	store.Set(AssetPegHashStoreKey(pegHash), bz)
}

// GetAssetPeg : get asset peg
func (am AssetPegMapper) GetAssetPeg(ctx sdk.Context, pegHash sdk.PegHash) types.AssetPeg {
	store := ctx.KVStore(am.key)
	bz := store.Get(AssetPegHashStoreKey(pegHash))
	if bz == nil {
		return nil
	}
	acc := am.decodeAssetPeg(bz)
	return acc
}

// IterateAssets : iterate over assets in kv store and add assets
func (am AssetPegMapper) IterateAssets(ctx sdk.Context, process func(sdk.AssetPeg) (stop bool)) {
	store := ctx.KVStore(am.key)
	iter := sdk.KVStorePrefixIterator(store, []byte("PegHash:"))
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc := am.decodeAssetPeg(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}
