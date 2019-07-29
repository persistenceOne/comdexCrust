package fiatFactory

import (
	"github.com/commitHub/commitBlockchain/types"
	sdk "github.com/commitHub/commitBlockchain/types"
	wire "github.com/commitHub/commitBlockchain/wire"
)

//FiatPegMapper : encoder decoder for fiat type
type FiatPegMapper struct {
	key   sdk.StoreKey
	proto func() sdk.FiatPeg
	cdc   *wire.Codec
}

//NewFiatPegMapper : returns fiat mapper
func NewFiatPegMapper(cdc *wire.Codec, key sdk.StoreKey, proto func() sdk.FiatPeg) FiatPegMapper {
	return FiatPegMapper{
		key:   key,
		proto: proto,
		cdc:   cdc,
	}
}

//FiatPegHashStoreKey : converts peg hash to keystore key
func FiatPegHashStoreKey(fiatPegHash sdk.PegHash) []byte {
	return append([]byte("PegHash:"), fiatPegHash.Bytes()...)
}

func (fm FiatPegMapper) encodeFiatPeg(fiat sdk.FiatPeg) []byte {
	bz, err := fm.cdc.MarshalBinaryBare(fiat)
	if err != nil {
		panic(err)
	}
	return bz
}

func (fm FiatPegMapper) decodeFiatPeg(bz []byte) (fiat sdk.FiatPeg) {
	err := fm.cdc.UnmarshalBinaryBare(bz, &fiat)
	if err != nil {
		panic(err)
	}
	return
}

//SetFiatPeg : set fiat peg
func (fm FiatPegMapper) SetFiatPeg(ctx sdk.Context, fiat sdk.FiatPeg) {
	pegHash := fiat.GetPegHash()
	store := ctx.KVStore(fm.key)
	bz := fm.encodeFiatPeg(fiat)
	store.Set(FiatPegHashStoreKey(pegHash), bz)
}

//GetFiatPeg : get fiat peg
func (fm FiatPegMapper) GetFiatPeg(ctx sdk.Context, pegHash sdk.PegHash) types.FiatPeg {
	store := ctx.KVStore(fm.key)
	bz := store.Get(FiatPegHashStoreKey(pegHash))
	if bz == nil {
		return nil
	}
	acc := fm.decodeFiatPeg(bz)
	return acc
}

//IterateFiats : iterate over fiats in kv store and add fiats
func (fm FiatPegMapper) IterateFiats(ctx sdk.Context, process func(sdk.FiatPeg) (stop bool)) {
	store := ctx.KVStore(fm.key)
	iter := sdk.KVStorePrefixIterator(store, []byte("PegHash:"))
	for {
		if !iter.Valid() {
			return
		}

		val := iter.Value()
		acc := fm.decodeFiatPeg(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}
