package bank

import (
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/acl"
	"github.com/commitHub/commitBlockchain/x/assetFactory"
	"github.com/commitHub/commitBlockchain/x/fiatFactory"
)

//RegisterWire : Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/Send", nil)
	cdc.RegisterConcrete(MsgIssue{}, "cosmos-sdk/Issue", nil)
	cdc.RegisterConcrete(MsgDefineZones{}, "commit-blockchain/MsgDefineZones", nil)
	cdc.RegisterConcrete(MsgDefineOrganizations{}, "commit-blockchain/MsgDefineOrganizations", nil)
	cdc.RegisterConcrete(MsgDefineACLs{}, "commit-blockchain/MsgDefineACLs", nil)
	cdc.RegisterConcrete(MsgBankIssueAssets{}, "cosmos-sdk/MsgBankIssueAssets", nil)
	cdc.RegisterConcrete(MsgBankReleaseAssets{}, "commit-blockchain/MsgBankReleaseAssets", nil)
	cdc.RegisterConcrete(MsgBankRedeemAssets{}, "cosmos-sdk/MsgBankRedeemAssets", nil)
	cdc.RegisterConcrete(MsgBankIssueFiats{}, "cosmos-sdk/MsgBankIssueFiats", nil)
	cdc.RegisterConcrete(MsgBankRedeemFiats{}, "cosmos-sdk/MsgBankRedeemFiats", nil)
	cdc.RegisterConcrete(MsgBankSendAssets{}, "cosmos-sdk/MsgBankSendAssets", nil)
	cdc.RegisterConcrete(MsgBankSendFiats{}, "cosmos-sdk/MsgBankSendFiats", nil)
	cdc.RegisterConcrete(MsgBankSellerExecuteOrders{}, "cosmos-sdk/MsgBankSellerExecuteOrders", nil)
	cdc.RegisterConcrete(MsgBankBuyerExecuteOrders{}, "cosmos-sdk/MsgBankBuyerExecuteOrders", nil)
	cdc.RegisterConcrete(BuyerExecuteOrderBody{}, "commit-blockchain/BuyerExecuteOrderBody", nil)
	cdc.RegisterConcrete(SellerExecuteOrderBody{}, "commit-blockchain/SellerExecuteOrderBody", nil)
	cdc.RegisterConcrete(IssueAssetBody{}, "commit-blockchain/bankIssueAssetBody", nil)
	cdc.RegisterConcrete(IssueFiatBody{}, "commit-blockchain/bankIssueFiatBody", nil)
	cdc.RegisterConcrete(SendAssetBody{}, "commit-blockchain/bankSendAssetBody", nil)
	cdc.RegisterConcrete(SendFiatBody{}, "commit-blockchain/bankSendFiatBody", nil)
	cdc.RegisterConcrete(SendTxBody{}, "commit-blockchain/SendTxBody", nil)
	cdc.RegisterConcrete(RedeemAssetBody{}, "commit-blockchain/bankRedeemAssetBody", nil)
	cdc.RegisterConcrete(RedeemFiatBody{}, "commit-blockchain/bankRedeemFiatBody", nil)
	cdc.RegisterConcrete(ReleaseAssetBody{}, "commit-blockchain/ReleaseAssetBody", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	assetFactory.RegisterAssetPeg(msgCdc)
	fiatFactory.RegisterFiatPeg(msgCdc)
	acl.RegisterACLAccount(msgCdc)
	RegisterWire(msgCdc)
}
