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

func NewKeeper(cdc *codec.Codec, storeKey cTypes.StoreKey, accountKeeper assetFactoryTypes.AccountKeeper) Keeper {
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
		return nil, assetFactoryTypes.ErrInvalidString(assetFactoryTypes.DefaultCodeSpace, fmt.Sprintf("Asset with pegHash %s not found!!!", peghash))
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

func (k Keeper) InstantiateAndAssignAsset(ctx cTypes.Context, issueAsset assetFactoryTypes.IssueAsset) cTypes.Error {

	// asset, err := k.GetAssetPeg(ctx, issueAsset.AssetPeg.GetPegHash())
	// if err != nil {
	// 	return err
	// }
	// if !asset.GetOwnerAddress().Equals(issueAsset.IssuerAddress) {
	// 	return cTypes.ErrInvalidAddress(fmt.Sprintf("Cannot issue asset to %s address.", issueAsset.IssuerAddress.String()))
	// }
	// _ = issueAsset.AssetPeg.SetPegHash(issueAsset.AssetPeg.GetPegHash())
	_ = issueAsset.AssetPeg.SetOwnerAddress(issueAsset.ToAddress)
	_ = k.SetAssetPeg(ctx, issueAsset.AssetPeg)

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			assetFactoryTypes.EventTypeAssetFactoryIssueAsset,
			cTypes.NewAttribute("recipient", issueAsset.ToAddress.String()),
			cTypes.NewAttribute("issuer", issueAsset.IssuerAddress.String()),
			cTypes.NewAttribute("asset", issueAsset.AssetPeg.GetPegHash().String()),
		))

	return nil
}

func (k Keeper) InstantiateAndRedeemAsset(ctx cTypes.Context, ownerAddress cTypes.AccAddress,
	toAddress cTypes.AccAddress, peghash types.PegHash) cTypes.Error {

	asset, err := k.GetAssetPeg(ctx, peghash)
	if err != nil {
		return err
	}

	if !asset.GetOwnerAddress().Equals(ownerAddress) {
		return cTypes.ErrInvalidAddress(ownerAddress.String())
	}
	unsetAssetPeg := types.NewBaseAssetPegWithPegHash(peghash)
	_ = unsetAssetPeg.SetOwnerAddress(toAddress)
	_ = k.SetAssetPeg(ctx, &unsetAssetPeg)

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			assetFactoryTypes.EventTypeAssetFactoryRedeemAsset,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("lastOwner", ownerAddress.String()),
			cTypes.NewAttribute("asset", asset.GetPegHash().String()),
		))

	return nil
}

func (k Keeper) SendAssetToOrder(ctx cTypes.Context, fromAddress cTypes.AccAddress,
	toAddress cTypes.AccAddress, peghash types.PegHash) cTypes.Error {

	asset, err := k.GetAssetPeg(ctx, peghash)
	if err != nil {
		return err
	}
	if !asset.GetOwnerAddress().Equals(fromAddress) {
		return cTypes.ErrInvalidAddress(fromAddress.String())
	}
	_ = asset.SetOwnerAddress(cTypes.AccAddress(append(append(toAddress.Bytes(), fromAddress.Bytes()...), peghash.Bytes()...)))
	_ = k.SetAssetPeg(ctx, asset)

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			assetFactoryTypes.EventTypeAssetFactorySendAsset,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("sender", fromAddress.String()),
			cTypes.NewAttribute("asset", peghash.String()),
		))

	return nil
}

func (k Keeper) SendAssetFromOrder(ctx cTypes.Context, fromAddress cTypes.AccAddress,
	toAddress cTypes.AccAddress, peghash types.PegHash) cTypes.Error {

	asset, err := k.GetAssetPeg(ctx, peghash)
	if err != nil {
		return err
	}

	if !asset.GetOwnerAddress().Equals(cTypes.AccAddress(append(append(toAddress.Bytes(), fromAddress.Bytes()...), peghash.Bytes()...))) {
		return cTypes.ErrInvalidAddress(fromAddress.String())
	}

	_ = asset.SetOwnerAddress(toAddress)
	_ = k.SetAssetPeg(ctx, asset)

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			assetFactoryTypes.EventTypeAssetFactoryExecuteAsset,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("sender", fromAddress.String()),
			cTypes.NewAttribute("asset", peghash.String()),
		))

	return nil
}
