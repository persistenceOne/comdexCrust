package ibc

import (
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/assetFactory"
	"github.com/commitHub/commitBlockchain/x/fiatFactory"
)

// RegisterWire concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(IBCTransferMsg{}, "commit-blockchain/IBCTransferMsg", nil)
	cdc.RegisterConcrete(IBCReceiveMsg{}, "commit-blockchain/IBCReceiveMsg", nil)
	cdc.RegisterConcrete(MsgIssueAssets{}, "commit-blockchain/MsgIssueAssets", nil)
	cdc.RegisterConcrete(MsgRedeemAssets{}, "commit-blockchain/MsgRedeemAssets", nil)
	cdc.RegisterConcrete(MsgIssueFiats{}, "commit-blockchain/MsgIssueFiats", nil)
	cdc.RegisterConcrete(MsgRedeemFiats{}, "commit-blockchain/MsgRedeemFiats", nil)
	cdc.RegisterConcrete(MsgRelayIssueAssets{}, "commit-blockchain/MsgRelayIssueAssets", nil)
	cdc.RegisterConcrete(MsgRelayRedeemAssets{}, "commit-blockchain/MsgRelayRedeemAssets", nil)
	cdc.RegisterConcrete(MsgRelayIssueFiats{}, "commit-blockchain/MsgRelayIssueFiats", nil)
	cdc.RegisterConcrete(MsgRelayRedeemFiats{}, "commit-blockchain/MsgRelayRedeemFiats", nil)
	cdc.RegisterConcrete(IssueAssetBody{}, "commit-blockchain/ibcIssueAssetBody", nil)
	cdc.RegisterConcrete(IssueFiatBody{}, "commit-blockchain/ibcIssueFiatBody", nil)
	cdc.RegisterConcrete(MsgSendAssets{}, "comdex-blockchian/MsgSendAssets", nil)
	cdc.RegisterConcrete(MsgRelaySendAssets{}, "commit-blockchain/MsgRelaySendAssets", nil)
	cdc.RegisterConcrete(MsgSendFiats{}, "commit-blockchain/MsgSendFiats", nil)
	cdc.RegisterConcrete(MsgRelaySendFiats{}, "commit-blockchain/MsgRelaySendFiats", nil)
	cdc.RegisterConcrete(MsgBuyerExecuteOrders{}, "commit-blockchain/MsgBuyerExecuteOrders", nil)
	cdc.RegisterConcrete(MsgSellerExecuteOrders{}, "commit-blockchain/MsgSellerExecuteOrders", nil)
	cdc.RegisterConcrete(MsgRelayBuyerExecuteOrders{}, "commit-blockchain/MsgRelayBuyerExecuteOrders", nil)
	cdc.RegisterConcrete(MsgRelaySellerExecuteOrders{}, "commit-blockchain/MsgRelaySellerExecuteOrders", nil)
	/*
		cdc.RegisterConcrete(MsgBuyerExecuteOrders{}, "commit-blockchain/MsgBuyerExecuteOrders", nil)
		cdc.RegisterConcrete(MsgSellerExecuteOrders{}, "commit-blockchain/MsgSellerExecuteOrders", nil)
	*/
}

var msgCdc = wire.NewCodec()

func init() {

	assetFactory.RegisterAssetPeg(msgCdc)
	fiatFactory.RegisterFiatPeg(msgCdc)
	RegisterWire(msgCdc)
}
