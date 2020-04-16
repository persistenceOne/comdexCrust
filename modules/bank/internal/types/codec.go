package types

import (
	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/modules/acl"
	"github.com/persistenceOne/comdexCrust/modules/reputation"
	"github.com/persistenceOne/comdexCrust/types"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "cosmos-sdk/MsgMultiSend", nil)
	cdc.RegisterConcrete(MsgDefineZones{}, "persistence-blockchain/MsgDefineZones", nil)
	cdc.RegisterConcrete(DefineZone{}, "persistence-blockchain/DefineZone", nil)
	cdc.RegisterConcrete(MsgDefineOrganizations{}, "persistence-blockchain/MsgDefineOrganizations", nil)
	cdc.RegisterConcrete(MsgDefineACLs{}, "persistence-blockchain/MsgDefineACLs", nil)
	cdc.RegisterConcrete(MsgBankIssueAssets{}, "cosmos-sdk/MsgBankIssueAssets", nil)
	cdc.RegisterConcrete(MsgBankReleaseAssets{}, "persistence-blockchain/MsgBankReleaseAssets", nil)
	cdc.RegisterConcrete(MsgBankRedeemAssets{}, "cosmos-sdk/MsgBankRedeemAssets", nil)
	cdc.RegisterConcrete(MsgBankIssueFiats{}, "cosmos-sdk/MsgBankIssueFiats", nil)
	cdc.RegisterConcrete(MsgBankRedeemFiats{}, "cosmos-sdk/MsgBankRedeemFiats", nil)
	cdc.RegisterConcrete(MsgBankSendAssets{}, "cosmos-sdk/MsgBankSendAssets", nil)
	cdc.RegisterConcrete(MsgBankSendFiats{}, "cosmos-sdk/MsgBankSendFiats", nil)
	cdc.RegisterConcrete(MsgBankSellerExecuteOrders{}, "cosmos-sdk/MsgBankSellerExecuteOrders", nil)
	cdc.RegisterConcrete(MsgBankBuyerExecuteOrders{}, "cosmos-sdk/MsgBankBuyerExecuteOrders", nil)
	cdc.RegisterInterface((*acl.ACLAccount)(nil), nil)
	cdc.RegisterConcrete(&acl.BaseACLAccount{}, "persistence-blockchain/AclAccount", nil)
	cdc.RegisterInterface((*types.AssetPeg)(nil), nil)
	cdc.RegisterConcrete(&types.BaseAssetPeg{}, "persistence-blockchain/AssetPeg", nil)
	cdc.RegisterInterface((*types.FiatPeg)(nil), nil)
	cdc.RegisterConcrete(&types.BaseFiatPeg{}, "persistence-blockchain/FiatPeg", nil)
	cdc.RegisterInterface((*types.AccountReputation)(nil), nil)
	cdc.RegisterConcrete(&types.BaseAccountReputation{}, "persistence-blockchain/AccountReputation", nil)
	cdc.RegisterConcrete(reputation.MsgBuyerFeedbacks{}, "persistence-blockchain/MsgBuyerFeedbacks", nil)
	cdc.RegisterConcrete(reputation.MsgSellerFeedbacks{}, "persistence-blockchain/MsgSellerFeedbacks", nil)
}

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
