package reputation

import (
	"reflect"

	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/x/order"
)

// NewFeedbackHandler returns a handler for "feedback" type messages.
func NewFeedbackHandler(k Keeper, orderKeeper order.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgBuyerFeedbacks:
			return handleMsgBuyerFeedback(ctx, k, orderKeeper, msg)
		case MsgSellerFeedbacks:
			return handleMsgSellerFeedback(ctx, k, orderKeeper, msg)

		default:
			errMsg := "Unrecognized feedback Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

//HandeBuyerFeedbacks
func handleMsgBuyerFeedback(ctx sdk.Context, k Keeper, orderKeeper order.Keeper, msg MsgBuyerFeedbacks) sdk.Result {
	tags, err := k.SetBuyerRatingToFeedback(ctx, orderKeeper, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

//HandeSellerFeedbacks
func handleMsgSellerFeedback(ctx sdk.Context, k Keeper, orderKeeper order.Keeper, msg MsgSellerFeedbacks) sdk.Result {
	tags, err := k.SetSellerRatingToFeedback(ctx, orderKeeper, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
