package assetFactory

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/types"
)

func instantiateAndAssignAsset(ctx cTypes.Context, k Keeper, issueAsset IssueAsset) cTypes.Error {

	asset, err := k.GetAssetPeg(ctx, issueAsset.AssetPeg.GetPegHash())
	if err != nil {
		return err
	}
	if !asset.GetOwnerAddress().Equals(issueAsset.IssuerAddress) {
		return cTypes.ErrInvalidAddress(fmt.Sprintf("Cannot issue asset to %s address.", issueAsset.IssuerAddress.String()))
	}
	_ = issueAsset.AssetPeg.SetPegHash(asset.GetPegHash())
	_ = issueAsset.AssetPeg.SetOwnerAddress(issueAsset.ToAddress)
	_ = k.SetAssetPeg(ctx, issueAsset.AssetPeg)

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			EventTypeAssetFactoryIssueAsset,
			cTypes.NewAttribute("recipient", issueAsset.ToAddress.String()),
			cTypes.NewAttribute("issuer", issueAsset.IssuerAddress.String()),
			cTypes.NewAttribute("asset", issueAsset.AssetPeg.GetPegHash().String()),
		))

	return nil
}

func instantiateAndRedeemAsset(ctx cTypes.Context, keeper Keeper, ownerAddress cTypes.AccAddress,
	toAddress cTypes.AccAddress, peghash types.PegHash) cTypes.Error {

	asset, err := keeper.GetAssetPeg(ctx, peghash)
	if err != nil {
		return err
	}

	if !asset.GetOwnerAddress().Equals(ownerAddress) {
		return cTypes.ErrInvalidAddress(ownerAddress.String())
	}
	unsetAssetPeg := types.NewBaseAssetPegWithPegHash(peghash)
	_ = unsetAssetPeg.SetOwnerAddress(toAddress)
	_ = keeper.SetAssetPeg(ctx, &unsetAssetPeg)

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			EventTypeAssetFactoryRedeemAsset,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("lastOwner", ownerAddress.String()),
			cTypes.NewAttribute("asset", asset.GetPegHash().String()),
		))

	return nil
}

func sendAssetToOrder(ctx cTypes.Context, keeper Keeper, fromAddress cTypes.AccAddress,
	toAddress cTypes.AccAddress, peghash types.PegHash) cTypes.Error {

	asset, err := keeper.GetAssetPeg(ctx, peghash)
	if err != nil {
		return err
	}
	if !asset.GetOwnerAddress().Equals(fromAddress) {
		return cTypes.ErrInvalidAddress(fromAddress.String())
	}
	_ = asset.SetOwnerAddress(cTypes.AccAddress(append(append(toAddress.Bytes(), fromAddress.Bytes()...), peghash.Bytes()...)))
	_ = keeper.SetAssetPeg(ctx, asset)

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			EventTypeAssetFactorySendAsset,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("sender", fromAddress.String()),
			cTypes.NewAttribute("asset", peghash.String()),
		))

	return nil
}

func sendAssetFromOrder(ctx cTypes.Context, keeper Keeper, fromAddress cTypes.AccAddress,
	toAddress cTypes.AccAddress, peghash types.PegHash) cTypes.Error {

	asset, err := keeper.GetAssetPeg(ctx, peghash)
	if err != nil {
		return err
	}

	if !asset.GetOwnerAddress().Equals(cTypes.AccAddress(append(append(toAddress.Bytes(), fromAddress.Bytes()...), peghash.Bytes()...))) {
		return cTypes.ErrInvalidAddress(fromAddress.String())
	}

	_ = asset.SetOwnerAddress(toAddress)
	_ = keeper.SetAssetPeg(ctx, asset)

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			EventTypeAssetFactoryExecuteAsset,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("sender", fromAddress.String()),
			cTypes.NewAttribute("asset", peghash.String()),
		))

	return nil
}
