package bank

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/modules/bank/internal/keeper"
	"github.com/commitHub/commitBlockchain/modules/bank/internal/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k keeper.Keeper) cTypes.Handler {
	return func(ctx cTypes.Context, msg cTypes.Msg) cTypes.Result {
		ctx = ctx.WithEventManager(cTypes.NewEventManager())

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
			return cTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx cTypes.Context, k keeper.Keeper, msg types.MsgSend) cTypes.Result {
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			types.EventTypeTransfer,
			cTypes.NewAttribute(cTypes.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return cTypes.Result{Events: ctx.EventManager().Events()}
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx cTypes.Context, k keeper.Keeper, msg types.MsgMultiSend) cTypes.Result {
	// NOTE: totalIn == totalOut should already have been checked
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			cTypes.EventTypeMessage,
			cTypes.NewAttribute(cTypes.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return cTypes.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBankIssueAssets(ctx cTypes.Context, k keeper.Keeper, msg types.MsgBankIssueAssets) cTypes.Result {

	for _, issueAsset := range msg.IssueAssets {
		err := k.IssueAssetsToWallets(ctx, issueAsset)

		if err != nil {
			return err.Result()
		}
	}
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankIssueFiats(ctx cTypes.Context, k keeper.Keeper, msg types.MsgBankIssueFiats) cTypes.Result {

	for _, issueFiat := range msg.IssueFiats {
		err := k.IssueFiatsToWallets(ctx, issueFiat)
		if err != nil {
			return err.Result()
		}
	}
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMMsgBankRedeemAssets(ctx cTypes.Context, k keeper.Keeper, msg types.MsgBankRedeemAssets) cTypes.Result {

	for _, redeemAsset := range msg.RedeemAssets {
		err := k.RedeemAssetsFromWallets(ctx, redeemAsset)
		if err != nil {
			return err.Result()
		}
	}
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMMsgBankRedeemFiats(ctx cTypes.Context, k keeper.Keeper, msg types.MsgBankRedeemFiats) cTypes.Result {

	for _, redeemFiat := range msg.RedeemFiats {
		err := k.RedeemFiatsFromWallets(ctx, redeemFiat)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankSendAssets(ctx cTypes.Context, k keeper.Keeper, msg types.MsgBankSendAssets) cTypes.Result {

	for _, sendAsset := range msg.SendAssets {
		err := k.SendAssetsToWallets(ctx, sendAsset)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankSendFiats(ctx cTypes.Context, k keeper.Keeper, msg types.MsgBankSendFiats) cTypes.Result {

	for _, sendFiat := range msg.SendFiats {
		err := k.SendFiatsToWallets(ctx, sendFiat)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankBuyerExecuteOrders(ctx cTypes.Context, k keeper.Keeper, msg types.MsgBankBuyerExecuteOrders) cTypes.Result {

	for _, buyerExecuteOrder := range msg.BuyerExecuteOrders {
		err, _ := k.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankSellerExecuteOrders(ctx cTypes.Context, k keeper.Keeper, msg types.MsgBankSellerExecuteOrders) cTypes.Result {

	for _, sellerExecuteOrder := range msg.SellerExecuteOrders {
		err, _ := k.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBankReleaseAssets(ctx cTypes.Context, k keeper.Keeper, msg types.MsgBankReleaseAssets) cTypes.Result {

	for _, releaseAsset := range msg.ReleaseAssets {
		err := k.ReleaseLockedAssets(ctx, releaseAsset)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgDefineZones(ctx cTypes.Context, k keeper.Keeper, msg types.MsgDefineZones) cTypes.Result {

	for _, defineZone := range msg.DefineZones {
		err := k.DefineZones(ctx, defineZone)
		if err != nil {
			return err.Result()
		}
	}
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgDefineOrganizations(ctx cTypes.Context, k keeper.Keeper, msg types.MsgDefineOrganizations) cTypes.Result {

	for _, defineOrganization := range msg.DefineOrganizations {
		err := k.DefineOrganizations(ctx, defineOrganization)
		if err != nil {
			return err.Result()
		}
	}
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgDefineACLs(ctx cTypes.Context, k keeper.Keeper, msg types.MsgDefineACLs) cTypes.Result {

	for _, defineACL := range msg.DefineACLs {
		err := k.DefineACLs(ctx, defineACL)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}
