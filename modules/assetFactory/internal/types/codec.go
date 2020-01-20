package types

import (
	"github.com/persistenceOne/persistenceSDK/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgFactoryIssueAssets{}, "persistence-blockchain/MsgFactoryIssueAssets", nil)
	cdc.RegisterConcrete(MsgFactoryRedeemAssets{}, "persistence-blockchain/MsgFactoryRedeemAssets", nil)
	cdc.RegisterConcrete(MsgFactorySendAssets{}, "persistence-blockchain/MsgFactorySendAssets", nil)
	cdc.RegisterConcrete(MsgFactoryExecuteAssets{}, "persistence-blockchain/MsgFactoryExecuteAssets", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
