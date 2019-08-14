package keeper

import (
	"fmt"
	
	cTypes "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/commitHub/commitBlockchain/codec"
	fiatFactoryTypes "github.com/commitHub/commitBlockchain/modules/fiatFactory/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

type Keeper struct {
	storeKey      cTypes.StoreKey
	accountKeeper fiatFactoryTypes.AccountKeeper
	cdc           *codec.Codec
}

func NewKeeper(storeKey cTypes.StoreKey, accountKeeper fiatFactoryTypes.AccountKeeper, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:      storeKey,
		accountKeeper: accountKeeper,
		cdc:           cdc,
	}
}

func (k Keeper) SetFiatPeg(ctx cTypes.Context, fiatPeg types.FiatPeg) {
	store := ctx.KVStore(k.storeKey)
	
	fiatPegKey := fiatFactoryTypes.FiatPegHashStoreKey(fiatPeg.GetPegHash())
	bytes := k.cdc.MustMarshalBinaryLengthPrefixed(fiatPeg)
	store.Set(fiatPegKey, bytes)
}

func (k Keeper) GetFiatPeg(ctx cTypes.Context, pegHash types.PegHash) (types.FiatPeg, cTypes.Error) {
	store := ctx.KVStore(k.storeKey)
	
	assetKey := fiatFactoryTypes.FiatPegHashStoreKey(pegHash)
	data := store.Get(assetKey)
	if data == nil {
		return nil, fiatFactoryTypes.ErrInvalidString(fiatFactoryTypes.DefaultCodeSpace, fmt.Sprintf("Fiat with pegHash %s not found", pegHash))
	}
	
	var fiatPeg types.FiatPeg
	k.cdc.MustUnmarshalBinaryLengthPrefixed(data, &fiatPeg)
	
	return fiatPeg, nil
}
