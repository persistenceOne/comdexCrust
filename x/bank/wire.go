package bank

import (
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/acl"
	"github.com/comdex-blockchain/x/assetFactory"
	"github.com/comdex-blockchain/x/fiatFactory"
)

// RegisterWire : Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/Send", nil)
	cdc.RegisterConcrete(MsgIssue{}, "cosmos-sdk/Issue", nil)
	cdc.RegisterConcrete(MsgDefineZones{}, "comdex-blockchain/MsgDefineZones", nil)
	cdc.RegisterConcrete(MsgDefineOrganizations{}, "comdex-blockchain/MsgDefineOrganizations", nil)
	cdc.RegisterConcrete(MsgDefineACLs{}, "comdex-blockchain/MsgDefineACLs", nil)
	cdc.RegisterConcrete(MsgBankIssueAssets{}, "cosmos-sdk/MsgBankIssueAssets", nil)
	cdc.RegisterConcrete(MsgBankReleaseAssets{}, "comdex-blockchain/MsgBankReleaseAssets", nil)
	cdc.RegisterConcrete(MsgBankRedeemAssets{}, "cosmos-sdk/MsgBankRedeemAssets", nil)
	cdc.RegisterConcrete(MsgBankIssueFiats{}, "cosmos-sdk/MsgBankIssueFiats", nil)
	cdc.RegisterConcrete(MsgBankRedeemFiats{}, "cosmos-sdk/MsgBankRedeemFiats", nil)
	cdc.RegisterConcrete(MsgBankSendAssets{}, "cosmos-sdk/MsgBankSendAssets", nil)
	cdc.RegisterConcrete(MsgBankSendFiats{}, "cosmos-sdk/MsgBankSendFiats", nil)
	cdc.RegisterConcrete(MsgBankSellerExecuteOrders{}, "cosmos-sdk/MsgBankSellerExecuteOrders", nil)
	cdc.RegisterConcrete(MsgBankBuyerExecuteOrders{}, "cosmos-sdk/MsgBankBuyerExecuteOrders", nil)
	cdc.RegisterConcrete(BuyerExecuteOrderBody{}, "comdex-blockchain/BuyerExecuteOrderBody", nil)
	cdc.RegisterConcrete(SellerExecuteOrderBody{}, "comdex-blockchain/SellerExecuteOrderBody", nil)
	cdc.RegisterConcrete(IssueAssetBody{}, "comdex-blockchain/bankIssueAssetBody", nil)
	cdc.RegisterConcrete(IssueFiatBody{}, "comdex-blockchain/bankIssueFiatBody", nil)
	cdc.RegisterConcrete(SendAssetBody{}, "comdex-blockchain/bankSendAssetBody", nil)
	cdc.RegisterConcrete(SendFiatBody{}, "comdex-blockchain/bankSendFiatBody", nil)
	cdc.RegisterConcrete(SendTxBody{}, "comdex-blockchain/SendTxBody", nil)
	cdc.RegisterConcrete(RedeemAssetBody{}, "comdex-blockchain/bankRedeemAssetBody", nil)
	cdc.RegisterConcrete(RedeemFiatBody{}, "comdex-blockchain/bankRedeemFiatBody", nil)
	cdc.RegisterConcrete(ReleaseAssetBody{}, "comdex-blockchain/ReleaseAssetBody", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	assetFactory.RegisterAssetPeg(msgCdc)
	fiatFactory.RegisterFiatPeg(msgCdc)
	acl.RegisterACLAccount(msgCdc)
	RegisterWire(msgCdc)
}
