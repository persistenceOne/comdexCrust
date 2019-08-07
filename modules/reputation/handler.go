package reputation

import (
	"reflect"
	"strconv"
	
	cTypes "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/commitHub/commitBlockchain/modules/negotiation"
	"github.com/commitHub/commitBlockchain/modules/reputation/internal/types"
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
	err := SetBuyerRatingToFeedback(ctx, k, msg)
	if err != nil {
		return err.Result()
	}
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgSellerFeedback(ctx cTypes.Context, k Keeper, msg MsgSellerFeedbacks) cTypes.Result {
	err := SetSellerRatingToFeedback(ctx, k, msg)
	if err != nil {
		return err.Result()
	}
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func SetBuyerRatingToFeedback(ctx cTypes.Context, k Keeper, msgFeedback MsgBuyerFeedbacks) cTypes.Error {
	
	for _, submitTraderFeedback := range msgFeedback.SubmitTraderFeedbacks {
		traderFeedback := submitTraderFeedback.TraderFeedback
		
		negotiationID := negotiation.NegotiationID(append(append(submitTraderFeedback.TraderFeedback.BuyerAddress.Bytes(),
			submitTraderFeedback.TraderFeedback.SellerAddress.Bytes()...), submitTraderFeedback.TraderFeedback.PegHash.Bytes()...))
		
		order := k.OrderKeeper.GetOrder(ctx, negotiationID)
		
		if order.GetFiatProofHash() == "" || order.GetAWBProofHash() == "" {
			return types.ErrFeedbackCannotRegister("you have not completed the transaction to give feedback")
		}
		
		err := k.SetFeedback(ctx, traderFeedback.SellerAddress, traderFeedback)
		if err != nil {
			return err
		}
		
		ctx.EventManager().EmitEvent(
			cTypes.NewEvent(
				EventTypeSetBuyerRatingToFeedback,
				cTypes.NewAttribute(AttributeKeyFrom, traderFeedback.BuyerAddress.String()),
				cTypes.NewAttribute(AttributeKeyTo, traderFeedback.SellerAddress.String()),
				cTypes.NewAttribute(AttributeKeyPegHash, traderFeedback.PegHash.String()),
				cTypes.NewAttribute(AttributeKeyRating, strconv.FormatInt(traderFeedback.Rating, 10)),
			))
	}
	
	return nil
}

func SetSellerRatingToFeedback(ctx cTypes.Context, k Keeper, msgFeedback MsgSellerFeedbacks) cTypes.Error {
	
	for _, submitTraderFeedback := range msgFeedback.SubmitTraderFeedbacks {
		traderFeedback := submitTraderFeedback.TraderFeedback
		
		negotiationID := negotiation.NegotiationID(append(append(traderFeedback.BuyerAddress.Bytes(),
			traderFeedback.SellerAddress.Bytes()...), traderFeedback.PegHash.Bytes()...))
		order := k.OrderKeeper.GetOrder(ctx, negotiationID)
		
		if order.GetFiatProofHash() == "" || order.GetAWBProofHash() == "" {
			return types.ErrFeedbackCannotRegister("you have not completed the transaction to give feedback")
		}
		
		err := k.SetFeedback(ctx, traderFeedback.BuyerAddress, traderFeedback)
		if err != nil {
			return err
		}
		
		ctx.EventManager().EmitEvent(
			cTypes.NewEvent(
				EventTypeSetSellerRatingToFeedback,
				cTypes.NewAttribute(AttributeKeyFrom, traderFeedback.SellerAddress.String()),
				cTypes.NewAttribute(AttributeKeyTo, traderFeedback.BuyerAddress.String()),
				cTypes.NewAttribute(AttributeKeyPegHash, traderFeedback.PegHash.String()),
				cTypes.NewAttribute(AttributeKeyRating, strconv.FormatInt(traderFeedback.Rating, 10)),
			))
	}
	
	return nil
}
