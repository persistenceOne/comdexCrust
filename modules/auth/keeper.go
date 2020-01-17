package auth

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/modules/auth/exported"
	"github.com/persistenceOne/persistenceSDK/modules/auth/types"
	"github.com/persistenceOne/persistenceSDK/modules/params/subspace"
)

// AccountKeeper encodes/decodes accounts using the go-amino (binary)
// encoding/decoding library.
type AccountKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key cTypes.StoreKey

	// The prototypical Account constructor.
	proto func() exported.Account

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	paramSubspace subspace.Subspace
}

// NewAccountKeeper returns a new cTypes.AccountKeeper that uses go-amino to
// (binary) encode and decode concrete cTypes.Accounts.
// nolint
func NewAccountKeeper(
	cdc *codec.Codec, key cTypes.StoreKey, paramstore subspace.Subspace, proto func() exported.Account,
) AccountKeeper {

	return AccountKeeper{
		key:           key,
		proto:         proto,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
	}
}

// Logger returns a module-specific logger.
func (ak AccountKeeper) Logger(ctx cTypes.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// NewAccountWithAddress implements cTypes.AccountKeeper.
func (ak AccountKeeper) NewAccountWithAddress(ctx cTypes.Context, addr cTypes.AccAddress) exported.Account {
	acc := ak.proto()
	err := acc.SetAddress(addr)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	err = acc.SetAccountNumber(ak.GetNextAccountNumber(ctx))
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	return acc
}

// NewAccount creates a new account
func (ak AccountKeeper) NewAccount(ctx cTypes.Context, acc exported.Account) exported.Account {
	if err := acc.SetAccountNumber(ak.GetNextAccountNumber(ctx)); err != nil {
		panic(err)
	}
	return acc
}

// GetAccount implements cTypes.AccountKeeper.
func (ak AccountKeeper) GetAccount(ctx cTypes.Context, addr cTypes.AccAddress) exported.Account {
	store := ctx.KVStore(ak.key)
	bz := store.Get(types.AddressStoreKey(addr))
	if bz == nil {
		return nil
	}
	acc := ak.decodeAccount(bz)
	return acc
}

// GetAllAccounts returns all accounts in the accountKeeper.
func (ak AccountKeeper) GetAllAccounts(ctx cTypes.Context) []exported.Account {
	accounts := []exported.Account{}
	appendAccount := func(acc exported.Account) (stop bool) {
		accounts = append(accounts, acc)
		return false
	}
	ak.IterateAccounts(ctx, appendAccount)
	return accounts
}

// SetAccount implements cTypes.AccountKeeper.
func (ak AccountKeeper) SetAccount(ctx cTypes.Context, acc exported.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	bz, err := ak.cdc.MarshalBinaryBare(acc)
	if err != nil {
		panic(err)
	}
	store.Set(types.AddressStoreKey(addr), bz)
}

// RemoveAccount removes an account for the account mapper store.
// NOTE: this will cause supply invariant violation if called
func (ak AccountKeeper) RemoveAccount(ctx cTypes.Context, acc exported.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	store.Delete(types.AddressStoreKey(addr))
}

// IterateAccounts implements cTypes.AccountKeeper.
func (ak AccountKeeper) IterateAccounts(ctx cTypes.Context, process func(exported.Account) (stop bool)) {
	store := ctx.KVStore(ak.key)
	iter := cTypes.KVStorePrefixIterator(store, types.AddressStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc := ak.decodeAccount(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}

// GetPubKey Returns the PubKey of the account at address
func (ak AccountKeeper) GetPubKey(ctx cTypes.Context, addr cTypes.AccAddress) (crypto.PubKey, cTypes.Error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return nil, cTypes.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}
	return acc.GetPubKey(), nil
}

// GetSequence Returns the Sequence of the account at address
func (ak AccountKeeper) GetSequence(ctx cTypes.Context, addr cTypes.AccAddress) (uint64, cTypes.Error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return 0, cTypes.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}
	return acc.GetSequence(), nil
}

// GetNextAccountNumber Returns and increments the global account number counter
func (ak AccountKeeper) GetNextAccountNumber(ctx cTypes.Context) uint64 {
	var accNumber uint64
	store := ctx.KVStore(ak.key)
	bz := store.Get(types.GlobalAccountNumberKey)
	if bz == nil {
		accNumber = 0
	} else {
		err := ak.cdc.UnmarshalBinaryLengthPrefixed(bz, &accNumber)
		if err != nil {
			panic(err)
		}
	}

	bz = ak.cdc.MustMarshalBinaryLengthPrefixed(accNumber + 1)
	store.Set(types.GlobalAccountNumberKey, bz)

	return accNumber
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the auth module's parameters.
func (ak AccountKeeper) SetParams(ctx cTypes.Context, params types.Params) {
	ak.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (ak AccountKeeper) GetParams(ctx cTypes.Context) (params types.Params) {
	ak.paramSubspace.GetParamSet(ctx, &params)
	return
}

// -----------------------------------------------------------------------------
// Misc.

func (ak AccountKeeper) decodeAccount(bz []byte) (acc exported.Account) {
	err := ak.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}

func (ak AccountKeeper) GetNextAssetPegHash(ctx cTypes.Context) int {
	var assetNumber int
	store := ctx.KVStore(ak.key)
	bz := store.Get(types.AssetPegHashKey)
	if bz == nil {
		assetNumber = 0
	} else {
		ak.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &assetNumber)
	}

	bz = ak.cdc.MustMarshalBinaryLengthPrefixed(assetNumber + 1)
	store.Set(types.AssetPegHashKey, bz)

	return assetNumber
}

// GetNextFiatPegHash : Returns and increments the fiatPeg counter
func (am AccountKeeper) GetNextFiatPegHash(ctx cTypes.Context) int {
	var fiatNumber int
	store := ctx.KVStore(am.key)
	bz := store.Get(types.FiatPegHashKey)
	if bz == nil {
		fiatNumber = 0
	} else {
		am.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &fiatNumber)
	}

	bz = am.cdc.MustMarshalBinaryLengthPrefixed(fiatNumber + 1)
	store.Set(types.FiatPegHashKey, bz)

	return fiatNumber
}
