package negotiation

import (
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// Mapper : encoder decoder for negotiation type
type Mapper struct {
	key   sdk.StoreKey
	proto func() sdk.Negotiation
	cdc   *wire.Codec
}

// NewMapper : returns negotiation mapper
func NewMapper(cdc *wire.Codec, key sdk.StoreKey, proto func() sdk.Negotiation) Mapper {
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

func (nm Mapper) encodeNegotiation(negotiation sdk.Negotiation) []byte {
	bz, err := nm.cdc.MarshalBinaryBare(negotiation)
	if err != nil {
		panic(err)
	}
	return bz
}

func (nm Mapper) decodeNegotiation(bz []byte) (negotiation sdk.Negotiation) {
	err := nm.cdc.UnmarshalBinaryBare(bz, &negotiation)
	if err != nil {
		panic(err)
	}
	return
}

// SetNegotiation : Set Negotiation
func (nm Mapper) SetNegotiation(ctx sdk.Context, negotiation sdk.Negotiation) {
	negotiationID := negotiation.GetNegotiationID()
	store := ctx.KVStore(nm.key)
	bz := nm.encodeNegotiation(negotiation)
	store.Set(StoreKey(negotiationID), bz)
}

// GetNegotiation : get negotiation
func (nm Mapper) GetNegotiation(ctx sdk.Context, negotiationID sdk.NegotiationID) sdk.Negotiation {
	store := ctx.KVStore(nm.key)
	bz := store.Get(StoreKey(negotiationID))
	if bz == nil {
		return nil
	}
	acc := nm.decodeNegotiation(bz)
	return acc
}

// IterateNegotiations : iterate over negotiations in kv store and add negotiation
func (nm Mapper) IterateNegotiations(ctx sdk.Context, process func(sdk.Negotiation) (stop bool)) {
	store := ctx.KVStore(nm.key)
	iter := sdk.KVStorePrefixIterator(store, []byte("negotiation:"))
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		negotiation := nm.decodeNegotiation(val)
		if process(negotiation) {
			return
		}
		iter.Next()
	}
}

// NewNegotiation : Crates a new negotiation in the kv store
func (nm Mapper) NewNegotiation(buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash) sdk.Negotiation {
	negotiation := sdk.ProtoBaseNegotiation()
	negotiationID := sdk.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	negotiation.SetNegotiationID(negotiationID)
	negotiation.SetBuyerAddress(buyerAddress)
	negotiation.SetSellerAddress(sellerAddress)
	negotiation.SetPegHash(pegHash)
	return negotiation
}
