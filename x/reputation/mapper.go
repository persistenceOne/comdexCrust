package reputation

import (
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// Mapper : for txns to msg mapping
type Mapper struct {
	key   sdk.StoreKey
	proto func() sdk.AccountReputation
	cdc   *wire.Codec
}

// NewMapper : creates newMapper
func NewMapper(cdc *wire.Codec, key sdk.StoreKey, proto func() sdk.AccountReputation) Mapper {
	return Mapper{
		key:   key,
		proto: proto,
		cdc:   cdc,
	}
}

// AccountStoreKey : create the key for person's account
func AccountStoreKey(addr sdk.AccAddress) []byte {
	return append([]byte("address:"), addr.Bytes()...)
}

func (fm Mapper) encodeAccountReputation(accountReputation sdk.AccountReputation) []byte {
	bz, err := fm.cdc.MarshalBinaryBare(accountReputation)
	if err != nil {
		panic(err)
	}
	return bz
}

func (fm Mapper) decodeAccountReputation(bz []byte) (accountReputation sdk.AccountReputation) {
	err := fm.cdc.UnmarshalBinaryBare(bz, &accountReputation)
	if err != nil {
		panic(err)
	}
	return
}

// GetAccountReputation : gets account from store
func (fm Mapper) GetAccountReputation(ctx sdk.Context, addr sdk.AccAddress) sdk.AccountReputation {
	store := ctx.KVStore(fm.key)
	bz := store.Get(AccountStoreKey(addr))
	if bz == nil {
		accountReputation := sdk.NewAccountReputation()
		accountReputation.SetAddress(addr)
		fm.SetAccountReputation(ctx, accountReputation)
		bz = store.Get(AccountStoreKey(addr))
	}
	accountReputation := fm.decodeAccountReputation(bz)
	return accountReputation
}

// SetAccountReputation : sets account to store
func (fm Mapper) SetAccountReputation(ctx sdk.Context, accountReputation sdk.AccountReputation) {
	addr := accountReputation.GetAddress()
	store := ctx.KVStore(fm.key)
	bz := fm.encodeAccountReputation(accountReputation)
	store.Set(AccountStoreKey(addr), bz)
}
