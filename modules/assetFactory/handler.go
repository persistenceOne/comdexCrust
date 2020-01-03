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
			return handleMsgFactoryExecuteAssets(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return cTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgFactoryIssueAssets(ctx cTypes.Context, keeper Keeper, msg MsgFactoryIssueAssets) cTypes.Result {
	fmt.Println("\n \n \n \n HANDLER \n \n \n")
	for _, issueAsset := range msg.IssueAssets {
		err := keeper.InstantiateAndAssignAsset(ctx, issueAsset)
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
		err := keeper.InstantiateAndRedeemAsset(ctx, redeemAsset.OwnerAddress,
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
		err := keeper.SendAssetToOrder(ctx, sendAsset.FromAddress, sendAsset.ToAddress, sendAsset.PegHash)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgFactoryExecuteAssets(ctx cTypes.Context, keeper Keeper, msg MsgFactoryExecuteAssets) cTypes.Result {
	for _, executeAsset := range msg.SendAssets {
		err := keeper.SendAssetFromOrder(ctx, executeAsset.FromAddress, executeAsset.ToAddress, executeAsset.PegHash)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}
