package assetFactory

import (
	"github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// RegisterWire : Register concrete types on wire codec for default assets
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgFactoryIssueAssets{}, "comdex-blockchain/MsgFactoryIssueAssets", nil)
	cdc.RegisterConcrete(MsgFactorySendAssets{}, "comdex-blockchain/MsgFactorySendAssets", nil)
	cdc.RegisterConcrete(MsgFactoryExecuteAssets{}, "comdex-blockchain/MsgFactoryExecuteAssets", nil)
	cdc.RegisterConcrete(IssueAssetBody{}, "comdex-blockchain/IssueAssetBody", nil)
	cdc.RegisterConcrete(RedeemAssetBody{}, "comdex-blockchain/RedeemAssetBody", nil)
	cdc.RegisterConcrete(ExecuteAssetBody{}, "comdex-blockchain/ExecuteAssetBody", nil)
	cdc.RegisterConcrete(SendAssetBody{}, "comdex-blockchain/SendAssetBody", nil)
	
}

// RegisterAssetPeg : Register assetFactory types and interfaces
func RegisterAssetPeg(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.AssetPeg)(nil), nil)
	cdc.RegisterConcrete(&types.BaseAssetPeg{}, "comdex-blockchain/AssetPeg", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterAssetPeg(msgCdc)
	RegisterWire(msgCdc)
	
}
