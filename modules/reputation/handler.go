package reputation

import (
	"reflect"

	cTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) cTypes.Handler {
	return func(ctx cTypes.Context, msg cTypes.Msg) cTypes.Result {
		switch msg := msg.(type) {
		case MsgBuyerFeedbacks:
			return handleMsgBuyerFeedback(ctx, k, msg)
		case MsgSellerFeedbacks:
			return handleMsgSellerFeedback(ctx, k, msg)

		default:
			errMsg := "Unrecognized feedback Msg type: " + reflect.TypeOf(msg).Name()
			return cTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgBuyerFeedback(ctx cTypes.Context, k Keeper, msg MsgBuyerFeedbacks) cTypes.Result {
	for _, submitTraderFeedback := range msg.SubmitTraderFeedbacks {
		err := k.SetBuyerRatingToFeedback(ctx, submitTraderFeedback)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgSellerFeedback(ctx cTypes.Context, k Keeper, msg MsgSellerFeedbacks) cTypes.Result {

	for _, submitTraderFeedback := range msg.SubmitTraderFeedbacks {
		err := k.SetSellerRatingToFeedback(ctx, submitTraderFeedback)
		if err != nil {
			return err.Result()
		}
	}
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}
