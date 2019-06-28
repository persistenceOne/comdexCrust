package assetFactory

import (
	"fmt"
	
	sdk "github.com/comdex-blockchain/types"
)

// Keeper : asset keeper
type Keeper struct {
	am AssetPegMapper
}

// NewKeeper : return a new keeper
func NewKeeper(am AssetPegMapper) Keeper {
	return Keeper{am: am}
}

// *****comdex

// set wallet
func instantiateAndAssignAsset(ctx sdk.Context, am AssetPegMapper, issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, newAsset sdk.AssetPeg) (sdk.AssetPeg, sdk.Tags, sdk.Error) {
	asset := am.GetAssetPeg(ctx, newAsset.GetPegHash())
	if asset == nil {
		return asset, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("Asset %s not found", newAsset.GetPegHash()))
	}
	if asset.GetOwnerAddress().String() != issuerAddress.String() {
		return asset, nil, sdk.ErrInvalidAddress(fmt.Sprintf("%s", issuerAddress))
	}
	newAsset.SetPegHash(asset.GetPegHash())
	newAsset.SetOwnerAddress(toAddress)
	am.SetAssetPeg(ctx, newAsset)
	tags := sdk.NewTags("recepient", []byte(toAddress.String()))
	tags = tags.AppendTag("issuer", []byte(issuerAddress.String()))
	tags = tags.AppendTag("asset", []byte(newAsset.GetPegHash().String()))
	return newAsset, tags, nil
}

// intiate the and assign each wallet
func issueAssetsToWallets(ctx sdk.Context, am AssetPegMapper, issueAssets []IssueAsset) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()
	
	for _, req := range issueAssets {
		_, tags, err := instantiateAndAssignAsset(ctx, am, req.IssuerAddress, req.ToAddress, req.AssetPeg)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

// IssueAssetsToWallets handles a list of IssueAsset messages
func (keeper Keeper) IssueAssetsToWallets(ctx sdk.Context, issueAssets []IssueAsset) (sdk.Tags, sdk.Error) {
	return issueAssetsToWallets(ctx, keeper.am, issueAssets)
}

// set wallet
func instantiateAndRedeemAsset(ctx sdk.Context, am AssetPegMapper, ownerAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash) (sdk.AssetPeg, sdk.Tags, sdk.Error) {
	asset := am.GetAssetPeg(ctx, pegHash)
	if asset == nil {
		return asset, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("Asset %s not found", pegHash))
	}
	if asset.GetOwnerAddress().String() != ownerAddress.String() {
		return asset, nil, sdk.ErrInvalidAddress(fmt.Sprintf("%s", ownerAddress))
	}
	unSetAssetPeg := sdk.NewBaseAssetPegWithPegHash(pegHash)
	unSetAssetPeg.SetOwnerAddress(toAddress)
	am.SetAssetPeg(ctx, &unSetAssetPeg)
	tags := sdk.NewTags("recepient", []byte(toAddress.String()))
	tags = tags.AppendTag("lastOwner", []byte(ownerAddress.String()))
	tags = tags.AppendTag("asset", []byte(asset.GetPegHash().String()))
	return asset, tags, nil
}

// intiate the and redeem each wallet
func redeemAssetsToWallets(ctx sdk.Context, am AssetPegMapper, redeemAssets []RedeemAsset) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()
	
	for _, req := range redeemAssets {
		_, tags, err := instantiateAndRedeemAsset(ctx, am, req.OwnerAddress, req.ToAddress, req.PegHash)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

// RedeemAssetsToWallets handles a list of RedeemAsset messages
func (keeper Keeper) RedeemAssetsToWallets(ctx sdk.Context, redeemAssets []RedeemAsset) (sdk.Tags, sdk.Error) {
	return redeemAssetsToWallets(ctx, keeper.am, redeemAssets)
}

// ***** Send Assets
func sendAssetToOrder(ctx sdk.Context, am AssetPegMapper, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash) (sdk.AssetPeg, sdk.Tags, sdk.Error) {
	asset := am.GetAssetPeg(ctx, pegHash)
	if asset == nil {
		return asset, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s", pegHash))
	}
	if asset.GetOwnerAddress().String() != fromAddress.String() {
		return asset, nil, sdk.ErrInvalidAddress(fmt.Sprintf("%s", fromAddress))
	}
	
	asset.SetOwnerAddress(sdk.AccAddress(append(append(toAddress.Bytes(), fromAddress.Bytes()...), pegHash.Bytes()...)))
	am.SetAssetPeg(ctx, asset)
	tags := sdk.NewTags("recepient", []byte(toAddress.String()))
	tags = tags.AppendTag("sender", []byte(fromAddress.String()))
	tags = tags.AppendTag("asset", []byte(pegHash))
	return asset, tags, nil
}

func sendAssetsToOrders(ctx sdk.Context, am AssetPegMapper, sendAssets []SendAsset) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()
	
	for _, req := range sendAssets {
		_, tags, err := sendAssetToOrder(ctx, am, req.FromAddress, req.ToAddress, req.PegHash)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

// SendAssetsToOrders handles a list of SendAsset messages
func (keeper Keeper) SendAssetsToOrders(ctx sdk.Context, sendAssets []SendAsset) (sdk.Tags, sdk.Error) {
	return sendAssetsToOrders(ctx, keeper.am, sendAssets)
}

// ##### Send Assets
// ***** Execute Assets
func sendAssetFromOrder(ctx sdk.Context, am AssetPegMapper, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash) (sdk.AssetPeg, sdk.Tags, sdk.Error) {
	asset := am.GetAssetPeg(ctx, pegHash)
	if asset == nil {
		return asset, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s", pegHash))
	}
	if asset.GetOwnerAddress().String() != sdk.AccAddress(append(append(toAddress.Bytes(), fromAddress.Bytes()...), pegHash.Bytes()...)).String() {
		return asset, nil, sdk.ErrInvalidAddress(fmt.Sprintf("%s", fromAddress))
	}
	
	asset.SetOwnerAddress(toAddress)
	am.SetAssetPeg(ctx, asset)
	tags := sdk.NewTags("recepient", []byte(toAddress.String()))
	tags = tags.AppendTag("sender", []byte(fromAddress.String()))
	tags = tags.AppendTag("asset", []byte(pegHash))
	return asset, tags, nil
}

func executeAssetOrders(ctx sdk.Context, am AssetPegMapper, executeAssets []SendAsset) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()
	
	for _, req := range executeAssets {
		_, tags, err := sendAssetFromOrder(ctx, am, req.FromAddress, req.ToAddress, req.PegHash)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

// ExecuteAssetOrders handles a list of ExecuteAsset messages
func (keeper Keeper) ExecuteAssetOrders(ctx sdk.Context, executeAssets []SendAsset) (sdk.Tags, sdk.Error) {
	return executeAssetOrders(ctx, keeper.am, executeAssets)
}

// ##### Execute Assets
// #####comdex
