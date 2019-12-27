package fiatFactory

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) cTypes.Handler {
	return func(ctx cTypes.Context, msg cTypes.Msg) cTypes.Result {
		ctx = ctx.WithEventManager(cTypes.NewEventManager())
		switch msg := msg.(type) {
		case MsgFactoryIssueFiats:
			return handleMsgFactoryIssueFiat(ctx, keeper, msg)
		case MsgFactoryRedeemFiats:
			return handleMsgFactoryRedeemFiat(ctx, keeper, msg)
		case MsgFactorySendFiats:
			return handleMsgFactorySendFiat(ctx, keeper, msg)
		case MsgFactoryExecuteFiats:
			return handleMsgFactoryExecuteFiat(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return cTypes.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgFactoryIssueFiat(ctx cTypes.Context, keeper Keeper, msg MsgFactoryIssueFiats) cTypes.Result {

	for _, issueFiat := range msg.IssueFiats {
		err := keeper.InstantiateAndAssignFiat(ctx, issueFiat.IssuerAddress, issueFiat.ToAddress, issueFiat.FiatPeg)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgFactoryRedeemFiat(ctx cTypes.Context, keeper Keeper, msg MsgFactoryRedeemFiats) cTypes.Result {
	for _, redeemFiat := range msg.RedeemFiats {
		err := keeper.InstantiateAndRedeemFiat(ctx, redeemFiat.RelayerAddress, redeemFiat.FiatPegWallet)

		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgFactorySendFiat(ctx cTypes.Context, keeper Keeper, msg MsgFactorySendFiats) cTypes.Result {
	events := cTypes.EmptyEvents()
	ctx.EventManager().EmitEvents(events)

	for _, sendFiat := range msg.SendFiats {
		err := keeper.SendFiatToOrder(ctx, sendFiat.FromAddress, sendFiat.ToAddress,
			sendFiat.PegHash, sendFiat.FiatPegWallet)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgFactoryExecuteFiat(ctx cTypes.Context, keeper Keeper, msg MsgFactoryExecuteFiats) cTypes.Result {
	events := cTypes.EmptyEvents()
	ctx.EventManager().EmitEvents(events)

	for _, executeFiat := range msg.SendFiats {
		err := keeper.SendFiatFromOrder(ctx, executeFiat.FromAddress, executeFiat.ToAddress,
			executeFiat.PegHash, executeFiat.FiatPegWallet)
		if err != nil {
			return err.Result()
		}
	}

	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}
