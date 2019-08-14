package keeper

import (
	"fmt"
	
	cTypes "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/commitHub/commitBlockchain/codec"
	assetFactoryTypes "github.com/commitHub/commitBlockchain/modules/assetFactory/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

type Keeper struct {
	storeKey      cTypes.StoreKey
	accountKeeper assetFactoryTypes.AccountKeeper
	cdc           *codec.Codec
}

func NewKeeper(storeKey cTypes.StoreKey, accountKeeper assetFactoryTypes.AccountKeeper, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:      storeKey,
		accountKeeper: accountKeeper,
		cdc:           cdc,
	}
}

func (k Keeper) SetAssetPeg(ctx cTypes.Context, assetPeg types.AssetPeg) cTypes.Error {
	store := ctx.KVStore(k.storeKey)
	assetPegKey := assetFactoryTypes.AssetPegHashStoreKey(assetPeg.GetPegHash())
	bytes := k.cdc.MustMarshalBinaryLengthPrefixed(assetPeg)
	store.Set(assetPegKey, bytes)
	
	return nil
}

func (k Keeper) GetAssetPeg(ctx cTypes.Context, peghash types.PegHash) (types.AssetPeg, cTypes.Error) {
	store := ctx.KVStore(k.storeKey)
	
	assetKey := assetFactoryTypes.AssetPegHashStoreKey(peghash)
	data := store.Get(assetKey)
	if data == nil {
		return nil, assetFactoryTypes.ErrInvalidString(assetFactoryTypes.DefaultCodeSpace, fmt.Sprintf("Asset with pegHash %s not found", peghash))
	}
	
	var assetPeg types.AssetPeg
	k.cdc.MustUnmarshalBinaryLengthPrefixed(data, &assetPeg)
	
	return assetPeg, nil
}

func (k Keeper) IterateAssets(ctx cTypes.Context, handler func(assetPeg types.AssetPeg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	
	iterator := cTypes.KVStorePrefixIterator(store, assetFactoryTypes.PegHashKey)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var assetPeg types.AssetPeg
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &assetPeg)
		if handler(assetPeg) {
			break
		}
	}
}
