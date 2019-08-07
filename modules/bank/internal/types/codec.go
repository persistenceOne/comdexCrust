package types

import (
	"github.com/commitHub/commitBlockchain/types"
	
	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/reputation"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "cosmos-sdk/MsgMultiSend", nil)
	cdc.RegisterConcrete(MsgDefineZones{}, "commit-blockchain/MsgDefineZones", nil)
	cdc.RegisterConcrete(DefineZone{}, "commit-blockchain/DefineZone", nil)
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
	cdc.RegisterInterface((*acl.ACLAccount)(nil), nil)
	cdc.RegisterConcrete(&acl.BaseACLAccount{}, "commit-blockchain/AclAccount", nil)
	cdc.RegisterInterface((*types.AssetPeg)(nil), nil)
	cdc.RegisterConcrete(&types.BaseAssetPeg{}, "commit-blockchain/AssetPeg", nil)
	cdc.RegisterInterface((*types.FiatPeg)(nil), nil)
	cdc.RegisterConcrete(&types.BaseFiatPeg{}, "commit-blockchain/FiatPeg", nil)
	cdc.RegisterInterface((*reputation.AccountReputation)(nil), nil)
	cdc.RegisterConcrete(&reputation.BaseAccountReputation{}, "commit-blockchain/AccountReputation", nil)
	cdc.RegisterConcrete(reputation.MsgBuyerFeedbacks{}, "commit-blockchain/MsgBuyerFeedbacks", nil)
	cdc.RegisterConcrete(reputation.MsgSellerFeedbacks{}, "commit-blockchain/MsgSellerFeedbacks", nil)
}

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
