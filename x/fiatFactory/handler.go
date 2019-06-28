package fiatFactory

import (
	"reflect"
	
	sdk "github.com/comdex-blockchain/types"
)

// NewHandler returns a handler for "fiatFactory" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgFactoryIssueFiats:
			return handleMsgFactoryIssueFiat(ctx, k, msg)
		case MsgFactoryRedeemFiats:
			return handleMsgFactoryRedeemFiat(ctx, k, msg)
		case MsgFactorySendFiats:
			return handleMsgFactorySendFiats(ctx, k, msg)
		case MsgFactoryExecuteFiats:
			return handleMsgFactoryExecuteFiats(ctx, k, msg)
		default:
			errMsg := "Unrecognized fiatFactory Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgFactoryIssueFiat
func handleMsgFactoryIssueFiat(ctx sdk.Context, k Keeper, msg MsgFactoryIssueFiats) sdk.Result {
	tags, err := k.IssueFiatsToWallets(ctx, msg.IssueFiats)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

// Handle MsgFactoryRedeemFiat
func handleMsgFactoryRedeemFiat(ctx sdk.Context, k Keeper, msg MsgFactoryRedeemFiats) sdk.Result {
	tags, err := k.RedeemFiatsFromWallets(ctx, msg.RedeemFiats)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

// Hande MsgFactorySendFiats
func handleMsgFactorySendFiats(ctx sdk.Context, k Keeper, msg MsgFactorySendFiats) sdk.Result {
	tags, err := k.SendFiatsToOrders(ctx, msg.SendFiats)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

// Hande MsgFactoryExecuteFiats
func handleMsgFactoryExecuteFiats(ctx sdk.Context, k Keeper, msg MsgFactoryExecuteFiats) sdk.Result {
	tags, err := k.ExecuteFiatOrders(ctx, msg.SendFiats)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
