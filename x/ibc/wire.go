package ibc

import (
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/assetFactory"
	"github.com/comdex-blockchain/x/fiatFactory"
)

// RegisterWire concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(IBCTransferMsg{}, "comdex-blockchain/IBCTransferMsg", nil)
	cdc.RegisterConcrete(IBCReceiveMsg{}, "comdex-blockchain/IBCReceiveMsg", nil)
	cdc.RegisterConcrete(MsgIssueAssets{}, "comdex-blockchain/MsgIssueAssets", nil)
	cdc.RegisterConcrete(MsgRedeemAssets{}, "comdex-blockchain/MsgRedeemAssets", nil)
	cdc.RegisterConcrete(MsgIssueFiats{}, "comdex-blockchain/MsgIssueFiats", nil)
	cdc.RegisterConcrete(MsgRedeemFiats{}, "comdex-blockchain/MsgRedeemFiats", nil)
	cdc.RegisterConcrete(MsgRelayIssueAssets{}, "comdex-blockchain/MsgRelayIssueAssets", nil)
	cdc.RegisterConcrete(MsgRelayRedeemAssets{}, "comdex-blockchain/MsgRelayRedeemAssets", nil)
	cdc.RegisterConcrete(MsgRelayIssueFiats{}, "comdex-blockchain/MsgRelayIssueFiats", nil)
	cdc.RegisterConcrete(MsgRelayRedeemFiats{}, "comdex-blockchain/MsgRelayRedeemFiats", nil)
	cdc.RegisterConcrete(IssueAssetBody{}, "comdex-blockchain/ibcIssueAssetBody", nil)
	cdc.RegisterConcrete(IssueFiatBody{}, "comdex-blockchain/ibcIssueFiatBody", nil)
	cdc.RegisterConcrete(MsgSendAssets{}, "comdex-blockchian/MsgSendAssets", nil)
	cdc.RegisterConcrete(MsgRelaySendAssets{}, "comdex-blockchain/MsgRelaySendAssets", nil)
	cdc.RegisterConcrete(MsgSendFiats{}, "comdex-blockchain/MsgSendFiats", nil)
	cdc.RegisterConcrete(MsgRelaySendFiats{}, "comdex-blockchain/MsgRelaySendFiats", nil)
	cdc.RegisterConcrete(MsgBuyerExecuteOrders{}, "comdex-blockchain/MsgBuyerExecuteOrders", nil)
	cdc.RegisterConcrete(MsgSellerExecuteOrders{}, "comdex-blockchain/MsgSellerExecuteOrders", nil)
	cdc.RegisterConcrete(MsgRelayBuyerExecuteOrders{}, "comdex-blockchain/MsgRelayBuyerExecuteOrders", nil)
	cdc.RegisterConcrete(MsgRelaySellerExecuteOrders{}, "comdex-blockchain/MsgRelaySellerExecuteOrders", nil)
	/*
		cdc.RegisterConcrete(MsgBuyerExecuteOrders{}, "comdex-blockchain/MsgBuyerExecuteOrders", nil)
		cdc.RegisterConcrete(MsgSellerExecuteOrders{}, "comdex-blockchain/MsgSellerExecuteOrders", nil)
	*/
}

var msgCdc = wire.NewCodec()

func init() {
	
	assetFactory.RegisterAssetPeg(msgCdc)
	fiatFactory.RegisterFiatPeg(msgCdc)
	RegisterWire(msgCdc)
}
