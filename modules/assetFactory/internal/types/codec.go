package types

import (
	"github.com/commitHub/commitBlockchain/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgFactoryIssueAssets{}, "commit-blockchain/MsgFactoryIssueAssets", nil)
	cdc.RegisterConcrete(MsgFactoryRedeemAssets{}, "commit-blockchain/MsgFactoryRedeemAssets", nil)
	cdc.RegisterConcrete(MsgFactorySendAssets{}, "commit-blockchain/MsgFactorySendAssets", nil)
	cdc.RegisterConcrete(MsgFactoryExecuteAssets{}, "commit-blockchain/MsgFactoryExecuteAssets", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
