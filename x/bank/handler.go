package bank

import (
	"reflect"

	"github.com/commitHub/commitBlockchain/x/reputation"

	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/x/acl"
	"github.com/commitHub/commitBlockchain/x/negotiation"
	"github.com/commitHub/commitBlockchain/x/order"
)

// NewAssetFiatHandler returns a handler for "bank" type messages.
func NewAssetFiatHandler(k Keeper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, aclKeeper acl.Keeper, reputationKeeper reputation.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, k, msg)
		case MsgIssue:
			return handleMsgIssue(ctx, k, msg)
		case MsgBankIssueAssets:
			return handleMsgBankIssueAsset(ctx, k, msg, aclKeeper)
		case MsgBankRedeemAssets:
			return handleMsgBankRedeemAsset(ctx, k, msg, aclKeeper)
		case MsgBankIssueFiats:
			return handleMsgBankIssueFiat(ctx, k, msg, aclKeeper)
		case MsgBankRedeemFiats:
			return handleMsgBankRedeemFiat(ctx, k, msg, aclKeeper)
		case MsgBankSendAssets:
			return handleMsgBankSendAsset(ctx, k, orderKeeper, negotiationKeeper, msg, aclKeeper, reputationKeeper)
		case MsgBankSendFiats:
			return handleMsgBankSendFiat(ctx, k, orderKeeper, negotiationKeeper, msg, aclKeeper, reputationKeeper)
		case MsgBankBuyerExecuteOrders:
			return handleMsgBankBuyerExecuteOrders(ctx, k, negotiationKeeper, orderKeeper, msg, aclKeeper, reputationKeeper)
		case MsgBankSellerExecuteOrders:
			return handleMsgBankSellerExecuteOrders(ctx, k, negotiationKeeper, orderKeeper, msg, aclKeeper, reputationKeeper)
		case MsgBankReleaseAssets:
			return handleMsgBankReleaseAssets(ctx, k, msg, aclKeeper)
		case MsgDefineZones:
			return handleMsgDefineZones(ctx, k, aclKeeper, msg)
		case MsgDefineOrganizations:
			return handMsgDefineOrganizations(ctx, k, aclKeeper, msg)
		case MsgDefineACLs:
			return handleMsgDefineACLs(ctx, k, aclKeeper, msg)
		default:
			errMsg := "Unrecognized bank Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, k, msg)
		case MsgIssue:
			return handleMsgIssue(ctx, k, msg)
		default:
			errMsg := "Unrecognized bank Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k Keeper, msg MsgSend) sdk.Result {

	tags, err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// Handle MsgIssue.
func handleMsgIssue(ctx sdk.Context, k Keeper, msg MsgIssue) sdk.Result {
	panic("not implemented yet")
}

//Handle MsgBankIssueAssets
func handleMsgBankIssueAsset(ctx sdk.Context, k Keeper, msg MsgBankIssueAssets, ak acl.Keeper) sdk.Result {
	tags, err, _ := k.IssueAssetsToWallets(ctx, msg.IssueAssets, ak)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//Handle MsgBankRedeemAsset
func handleMsgBankRedeemAsset(ctx sdk.Context, k Keeper, msg MsgBankRedeemAssets, ak acl.Keeper) sdk.Result {
	tags, err, _ := k.RedeemAssetsFromWallets(ctx, msg.RedeemAssets, ak)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//Hande MsgBankIssueFiats
func handleMsgBankIssueFiat(ctx sdk.Context, k Keeper, msg MsgBankIssueFiats, ak acl.Keeper) sdk.Result {
	tags, err, _ := k.IssueFiatsToWallets(ctx, msg.IssueFiats, ak)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//Handle MsgBankRedeemFiats
func handleMsgBankRedeemFiat(ctx sdk.Context, k Keeper, msg MsgBankRedeemFiats, ak acl.Keeper) sdk.Result {
	tags, err, _ := k.RedeemFiatsFromWallets(ctx, msg.RedeemFiats, ak)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//Handle MsgBankSendAssets
func handleMsgBankSendAsset(ctx sdk.Context, k Keeper, orderKeeper order.Keeper, negotiationKeeper negotiation.Keeper, msg MsgBankSendAssets, ak acl.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err, _ := k.SendAssetsToWallets(ctx, orderKeeper, negotiationKeeper, msg.SendAssets, ak, reputationKeeper)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//Handle MsgBankSendFiats
func handleMsgBankSendFiat(ctx sdk.Context, k Keeper, orderKeeper order.Keeper, negotiationKeeper negotiation.Keeper, msg MsgBankSendFiats, ak acl.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err, _ := k.SendFiatsToWallets(ctx, orderKeeper, negotiationKeeper, msg.SendFiats, ak, reputationKeeper)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//Hande MsgBankBuyerExecuteOrders
func handleMsgBankBuyerExecuteOrders(ctx sdk.Context, k Keeper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, msg MsgBankBuyerExecuteOrders, ak acl.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err, _ := k.BuyerExecuteTradeOrders(ctx, negotiationKeeper, orderKeeper, msg.BuyerExecuteOrders, ak, reputationKeeper)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//Hande MsgBankSellerExecuteOrders
func handleMsgBankSellerExecuteOrders(ctx sdk.Context, k Keeper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, msg MsgBankSellerExecuteOrders, ak acl.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err, _ := k.SellerExecuteTradeOrders(ctx, negotiationKeeper, orderKeeper, msg.SellerExecuteOrders, ak, reputationKeeper)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//Hande MsgBankReleaseAssets
func handleMsgBankReleaseAssets(ctx sdk.Context, k Keeper, msg MsgBankReleaseAssets, ak acl.Keeper) sdk.Result {
	tags, err := k.ReleaseLockedAssets(ctx, msg.ReleaseAssets, ak)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgDefineZones(ctx sdk.Context, k Keeper, aclKeeper acl.Keeper, msg MsgDefineZones) sdk.Result {

	tags, err := k.DefineZones(ctx, aclKeeper, msg.DefineZones)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

func handMsgDefineOrganizations(ctx sdk.Context, k Keeper, aclKeeper acl.Keeper, msg MsgDefineOrganizations) sdk.Result {

	tags, err := k.DefineOrganizations(ctx, aclKeeper, msg.DefineOrganizations)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgDefineACLs(ctx sdk.Context, k Keeper, aclKeeper acl.Keeper, msg MsgDefineACLs) sdk.Result {

	tags, err := k.DefineACLs(ctx, aclKeeper, msg.DefineACLs)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}
