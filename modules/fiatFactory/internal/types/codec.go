package types

import (
	"github.com/commitHub/commitBlockchain/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgFactoryIssueFiats{}, "commit-blockchain/MsgFactoryIssueFiats", nil)
	cdc.RegisterConcrete(MsgFactoryRedeemFiats{}, "commit-blockchain/MsgFactoryRedeemFiats", nil)
	cdc.RegisterConcrete(MsgFactorySendFiats{}, "commit-blockchain/MsgFactorySendFiats", nil)
	cdc.RegisterConcrete(MsgFactoryExecuteFiats{}, "commit-blockchain/MsgFactoryExecuteFiats", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
