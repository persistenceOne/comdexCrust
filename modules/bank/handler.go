package bank

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/modules/bank/internal/keeper"
	"github.com/commitHub/commitBlockchain/modules/bank/internal/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSend:
			return handleMsgSend(ctx, k, msg)

		case types.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, msg)

		case types.MsgBankIssueAssets:
			return handleMsgBankIssueAssets(ctx, k, msg)

		case types.MsgBankIssueFiats:
			return handleMsgBankIssueFiats(ctx, k, msg)

		case types.MsgBankRedeemAssets:
			return handleMMsgBankRedeemAssets(ctx, k, msg)

		case types.MsgBankRedeemFiats:
			return handleMMsgBankRedeemFiats(ctx, k, msg)

		case types.MsgBankSendAssets:
			return handleMsgBankSendAssets(ctx, k, msg)

		case types.MsgBankSendFiats:
			return handleMsgBankSendFiats(ctx, k, msg)

		case types.MsgBankBuyerExecuteOrders:
			return handleMsgBankBuyerExecuteOrders(ctx, k, msg)

		case types.MsgBankSellerExecuteOrders:
			return handleMsgBankSellerExecuteOrders(ctx, k, msg)

		case types.MsgBankReleaseAssets:
			return handleMsgBankReleaseAssets(ctx, k, msg)

		case types.MsgDefineZones:
			return handleMsgDefineZones(ctx, k, msg)

		case types.MsgDefineOrganizations:
			return handleMsgDefineOrganizations(ctx, k, msg)

		case types.MsgDefineACLs:
			return handleMsgDefineACLs(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized bank message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgSend) sdk.Result {
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgMultiSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBankIssueAssets(ctx sdk.Context, k keeper.Keeper, msg types.MsgBankIssueAssets) sdk.Result {

	for _, issueAsset := range msg.IssueAssets {
		err := k.IssueAssetsToWallets(ctx, issueAsset)

		if err != nil {
			return err.Result()
		}
	}
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankIssueFiats(ctx sdk.Context, k keeper.Keeper, msg types.MsgBankIssueFiats) sdk.Result {

	for _, issueFiat := range msg.IssueFiats {
		err := k.IssueFiatsToWallets(ctx, issueFiat)
		if err != nil {
			return err.Result()
		}
	}
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMMsgBankRedeemAssets(ctx sdk.Context, k keeper.Keeper, msg types.MsgBankRedeemAssets) sdk.Result {

	for _, redeemAsset := range msg.RedeemAssets {
		err := k.RedeemAssetsFromWallets(ctx, redeemAsset)
		if err != nil {
			return err.Result()
		}
	}
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMMsgBankRedeemFiats(ctx sdk.Context, k keeper.Keeper, msg types.MsgBankRedeemFiats) sdk.Result {

	for _, redeemFiat := range msg.RedeemFiats {
		err := k.RedeemFiatsFromWallets(ctx, redeemFiat)
		if err != nil {
			return err.Result()
		}
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankSendAssets(ctx sdk.Context, k keeper.Keeper, msg types.MsgBankSendAssets) sdk.Result {

	for _, sendAsset := range msg.SendAssets {
		err := k.SendAssetsToWallets(ctx, sendAsset)
		if err != nil {
			return err.Result()
		}
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankSendFiats(ctx sdk.Context, k keeper.Keeper, msg types.MsgBankSendFiats) sdk.Result {

	for _, sendFiat := range msg.SendFiats {
		err := k.SendFiatsToWallets(ctx, sendFiat)
		if err != nil {
			return err.Result()
		}
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankBuyerExecuteOrders(ctx sdk.Context, k keeper.Keeper, msg types.MsgBankBuyerExecuteOrders) sdk.Result {

	for _, buyerExecuteOrder := range msg.BuyerExecuteOrders {
		err, _ := k.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
		if err != nil {
			return err.Result()
		}
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankSellerExecuteOrders(ctx sdk.Context, k keeper.Keeper, msg types.MsgBankSellerExecuteOrders) sdk.Result {

	for _, sellerExecuteOrder := range msg.SellerExecuteOrders {
		err, _ := k.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
		if err != nil {
			return err.Result()
		}
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankReleaseAssets(ctx sdk.Context, k keeper.Keeper, msg types.MsgBankReleaseAssets) sdk.Result {

	for _, releaseAsset := range msg.ReleaseAssets {
		err := k.ReleaseLockedAssets(ctx, releaseAsset)
		if err != nil {
			return err.Result()
		}
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgDefineZones(ctx sdk.Context, k keeper.Keeper, msg types.MsgDefineZones) sdk.Result {

	for _, defineZone := range msg.DefineZones {
		err := k.DefineZones(ctx, defineZone)
		if err != nil {
			return err.Result()
		}
	}
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgDefineOrganizations(ctx sdk.Context, k keeper.Keeper, msg types.MsgDefineOrganizations) sdk.Result {

	for _, defineOrganization := range msg.DefineOrganizations {
		err := k.DefineOrganizations(ctx, defineOrganization)
		if err != nil {
			return err.Result()
		}
	}
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgDefineACLs(ctx sdk.Context, k keeper.Keeper, msg types.MsgDefineACLs) sdk.Result {

	for _, defineACL := range msg.DefineACLs {
		err := k.DefineACLs(ctx, defineACL)
		if err != nil {
			return err.Result()
		}
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
