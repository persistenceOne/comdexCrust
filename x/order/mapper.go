package order

import (
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// Mapper : encoder decoder for order type
type Mapper struct {
	key   sdk.StoreKey
	proto func() sdk.Order
	cdc   *wire.Codec
}

// NewMapper : returns order mapper
func NewMapper(cdc *wire.Codec, key sdk.StoreKey, proto func() sdk.Order) Mapper {
	return Mapper{
		key:   key,
		proto: proto,
		cdc:   cdc,
	}
}

// StoreKey : converts negotiationID to keystore key
func StoreKey(negotiationID sdk.NegotiationID) []byte {
	return append([]byte("NegotiationID:"), negotiationID.Bytes()...)
}

func (om Mapper) encodeOrder(order sdk.Order) []byte {
	bz, err := om.cdc.MarshalBinaryBare(order)
	if err != nil {
		panic(err)
	}
	return bz
}

func (om Mapper) decodeOrder(bz []byte) (order sdk.Order) {
	err := om.cdc.UnmarshalBinaryBare(bz, &order)
	if err != nil {
		panic(err)
	}
	return
}

// SetOrder : Set Order
func (om Mapper) SetOrder(ctx sdk.Context, order sdk.Order) {
	negotiationID := order.GetNegotiationID()
	store := ctx.KVStore(om.key)
	bz := om.encodeOrder(order)
	store.Set(StoreKey(negotiationID), bz)
}

// GetOrder : get order
func (om Mapper) GetOrder(ctx sdk.Context, negotiationID sdk.NegotiationID) sdk.Order {
	store := ctx.KVStore(om.key)
	bz := store.Get(StoreKey(negotiationID))
	if bz == nil {
		return nil
	}
	acc := om.decodeOrder(bz)
	return acc
}

// IterateOrders : iterate over orders in kv store and add orders
func (om Mapper) IterateOrders(ctx sdk.Context, process func(sdk.Order) (stop bool)) {
	store := ctx.KVStore(om.key)
	iter := sdk.KVStorePrefixIterator(store, []byte("order:"))
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		order := om.decodeOrder(val)
		if process(order) {
			return
		}
		iter.Next()
	}
}

// NewOrder : Crates a new order in the kv store
func (om Mapper) NewOrder(buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash) sdk.Order {
	order := sdk.BaseOrder{}
	negotiationID := sdk.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	order.SetNegotiationID(negotiationID)
	return &order
}
