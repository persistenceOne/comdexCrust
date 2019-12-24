package transfer

import (
	"github.com/commitHub/commitBlockchain/modules/ibc/20-transfer/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleMsgTransfer defines the sdk.Handler for MsgTransfer
func HandleMsgTransfer(ctx sdk.Context, k Keeper, msg MsgTransfer) (res sdk.Result) {
	err := k.SendTransfer(ctx, msg.SourcePort, msg.SourceChannel, msg.Amount, msg.Sender, msg.Receiver, msg.Source)
	if err != nil {
		return sdk.NewError(DefaultCodespace, CodeErrSendPacket, err.Error()).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// HandleMsgRecvPacket defines the sdk.Handler for MsgRecvPacket
func HandleMsgRecvPacket(ctx sdk.Context, k Keeper, msg MsgRecvPacket) (res sdk.Result) {
	err := k.ReceivePacket(ctx, msg.Packet, msg.Proofs[0], msg.Height)
	if err != nil {
		return sdk.NewError(DefaultCodespace, CodeErrSendPacket, err.Error()).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleMsgIssueAssetTransfer(ctx sdk.Context, k Keeper, msg MsgIssueAssetTransfer) (res sdk.Result) {
	err := k.IssueAssetTransfer(ctx, msg.SourcePort, msg.SourceChannel, msg.IssueAsset)
	if err != nil {
		return sdk.NewError(DefaultCodespace, CodeErrSendPacket, err.Error()).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.IssueAsset.IssuerAddress.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.IssueAsset.IssuerAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleMsgReceiveIssueAssetPacket(ctx sdk.Context, k Keeper, msg MsgReceiveIssueAssetPacket) (res sdk.Result) {
	err := k.ReceiveIssueAssetPacket(ctx, msg.Packet, msg.Proofs[0], msg.Height)
	if err != nil {
		return sdk.NewError(DefaultCodespace, CodeErrSendPacket, err.Error()).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
