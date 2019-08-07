package assetFactory

import (
	"fmt"
	
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) cTypes.Handler {
	return func(ctx cTypes.Context, msg cTypes.Msg) cTypes.Result {
		
		ctx = ctx.WithEventManager(cTypes.NewEventManager())
		switch msg := msg.(type) {
		case MsgFactoryIssueAssets:
			return handleMsgFactoryIssueAssets(ctx, k, msg)
		case MsgFactoryRedeemAssets:
			return handleMsgFactoryRedeemAssets(ctx, k, msg)
		case MsgFactorySendAssets:
			return handleMsgFactorySendAssets(ctx, k, msg)
		case MsgFactoryExecuteAssets:
			return handleMsgFactoryExecuteAsses(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return cTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgFactoryIssueAssets(ctx cTypes.Context, keeper Keeper, msg MsgFactoryIssueAssets) cTypes.Result {
	for _, issueAsset := range msg.IssueAssets {
		err := instantiateAndAssignAsset(ctx, keeper, issueAsset)
		if err != nil {
			return err.Result()
		}
	}
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}
func handleMsgFactoryRedeemAssets(ctx cTypes.Context, keeper Keeper, msg MsgFactoryRedeemAssets) cTypes.Result {
	
	for _, redeemAsset := range msg.RedeemAssets {
		err := instantiateAndRedeemAsset(ctx, keeper, redeemAsset.OwnerAddress,
			redeemAsset.ToAddress, redeemAsset.PegHash)
		
		if err != nil {
			return err.Result()
		}
	}
	
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgFactorySendAssets(ctx cTypes.Context, keeper Keeper, msg MsgFactorySendAssets) cTypes.Result {
	for _, sendAsset := range msg.SendAssets {
		err := sendAssetToOrder(ctx, keeper, sendAsset.FromAddress, sendAsset.ToAddress, sendAsset.PegHash)
		if err != nil {
			return err.Result()
		}
	}
	
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgFactoryExecuteAsses(ctx cTypes.Context, keeper Keeper, msg MsgFactoryExecuteAssets) cTypes.Result {
	for _, executeAsset := range msg.SendAssets {
		err := sendAssetFromOrder(ctx, keeper, executeAsset.FromAddress, executeAsset.ToAddress, executeAsset.PegHash)
		if err != nil {
			return err.Result()
		}
	}
	
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}
