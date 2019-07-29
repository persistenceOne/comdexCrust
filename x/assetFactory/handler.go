package assetFactory

import (
	"reflect"

	sdk "github.com/commitHub/commitBlockchain/types"
)

// NewHandler returns a handler for "assetFactory" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgFactoryIssueAssets:
			return handleMsgFactoryIssueAsset(ctx, k, msg)
		case MsgFactoryRedeemAssets:
			return handleMsgFactoryRedeemAsset(ctx, k, msg)
		case MsgFactorySendAssets:
			return handleMsgFactorySendAssets(ctx, k, msg)
		case MsgFactoryExecuteAssets:
			return handleMsgFactoryExecuteAssets(ctx, k, msg)
		default:
			errMsg := "Unrecognized assetFactory Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

//handleMsgFactoryIssueAsset
func handleMsgFactoryIssueAsset(ctx sdk.Context, k Keeper, msg MsgFactoryIssueAssets) sdk.Result {
	tags, err := k.IssueAssetsToWallets(ctx, msg.IssueAssets)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//Handle MsgFactoryRedeemAsset
func handleMsgFactoryRedeemAsset(ctx sdk.Context, k Keeper, msg MsgFactoryRedeemAssets) sdk.Result {
	tags, err := k.RedeemAssetsToWallets(ctx, msg.RedeemAssets)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//handle MsgFactorySendAssets
func handleMsgFactorySendAssets(ctx sdk.Context, k Keeper, msg MsgFactorySendAssets) sdk.Result {
	tags, err := k.SendAssetsToOrders(ctx, msg.SendAssets)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//handle MsgFactoryExecuteAssets
func handleMsgFactoryExecuteAssets(ctx sdk.Context, k Keeper, msg MsgFactoryExecuteAssets) sdk.Result {
	tags, err := k.ExecuteAssetOrders(ctx, msg.SendAssets)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
