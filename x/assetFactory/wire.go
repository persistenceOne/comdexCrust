package assetFactory

import (
	"github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
)

//RegisterWire : Register concrete types on wire codec for default assets
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgFactoryIssueAssets{}, "commit-blockchain/MsgFactoryIssueAssets", nil)
	cdc.RegisterConcrete(MsgFactorySendAssets{}, "commit-blockchain/MsgFactorySendAssets", nil)
	cdc.RegisterConcrete(MsgFactoryExecuteAssets{}, "commit-blockchain/MsgFactoryExecuteAssets", nil)
	cdc.RegisterConcrete(IssueAssetBody{}, "commit-blockchain/IssueAssetBody", nil)
	cdc.RegisterConcrete(RedeemAssetBody{}, "commit-blockchain/RedeemAssetBody", nil)
	cdc.RegisterConcrete(ExecuteAssetBody{}, "commit-blockchain/ExecuteAssetBody", nil)
	cdc.RegisterConcrete(SendAssetBody{}, "commit-blockchain/SendAssetBody", nil)

}

//RegisterAssetPeg : Register assetFactory types and interfaces
func RegisterAssetPeg(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.AssetPeg)(nil), nil)
	cdc.RegisterConcrete(&types.BaseAssetPeg{}, "commit-blockchain/AssetPeg", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterAssetPeg(msgCdc)
	RegisterWire(msgCdc)

}
