package types

import (
	"github.com/persistenceOne/comdexCrust/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgFactoryIssueFiats{}, "persistence-blockchain/MsgFactoryIssueFiats", nil)
	cdc.RegisterConcrete(MsgFactoryRedeemFiats{}, "persistence-blockchain/MsgFactoryRedeemFiats", nil)
	cdc.RegisterConcrete(MsgFactorySendFiats{}, "persistence-blockchain/MsgFactorySendFiats", nil)
	cdc.RegisterConcrete(MsgFactoryExecuteFiats{}, "persistence-blockchain/MsgFactoryExecuteFiats", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
