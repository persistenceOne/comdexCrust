package negotiation

import (
	"reflect"

	"github.com/commitHub/commitBlockchain/x/reputation"

	"github.com/commitHub/commitBlockchain/x/acl"

	sdk "github.com/commitHub/commitBlockchain/types"
)

// NewHandler returns a handler for "negotiation" type messages.
func NewHandler(k Keeper, aclKeeper acl.Keeper, reputationKeeper reputation.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgChangeBuyerBids:
			return handleMsgChangeBuyerBids(ctx, k, msg, aclKeeper, reputationKeeper)
		case MsgChangeSellerBids:
			return handleMsgChangeSellerBids(ctx, k, msg, aclKeeper, reputationKeeper)
		case MsgConfirmBuyerBids:
			return handleMsgConfirmBuyerBids(ctx, k, msg, aclKeeper, reputationKeeper)
		case MsgConfirmSellerBids:
			return handleMsgConfirmSellerBids(ctx, k, msg, aclKeeper, reputationKeeper)
		default:
			errMsg := "Unrecognized negotiation Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgChangeBuyerBids(ctx sdk.Context, k Keeper, msg MsgChangeBuyerBids, ak acl.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err := k.ChangeNegotiationBids(ctx, msg.ChangeBids, ak, reputationKeeper)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgChangeSellerBids(ctx sdk.Context, k Keeper, msg MsgChangeSellerBids, ak acl.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err := k.ChangeNegotiationBids(ctx, msg.ChangeBids, ak, reputationKeeper)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
func handleMsgConfirmBuyerBids(ctx sdk.Context, k Keeper, msg MsgConfirmBuyerBids, ak acl.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err := k.ConfirmNegotiationBids(ctx, msg.ConfirmBids, ak, reputationKeeper)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgConfirmSellerBids(ctx sdk.Context, k Keeper, msg MsgConfirmSellerBids, ak acl.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err := k.ConfirmNegotiationBids(ctx, msg.ConfirmBids, ak, reputationKeeper)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
