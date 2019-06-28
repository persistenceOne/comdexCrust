package fiatFactory

import (
	"github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// RegisterWire : Register concrete types on wire codec for default fiats
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgFactoryIssueFiats{}, "comdex-blockchain/MsgFactoryIssueFiats", nil)
	cdc.RegisterConcrete(MsgFactoryRedeemFiats{}, "comdex-blockchain/MsgFactoryRedeemFiats", nil)
	cdc.RegisterConcrete(MsgFactorySendFiats{}, "comdex-blockchain/MsgFactorySendFiats", nil)
	cdc.RegisterConcrete(MsgFactoryExecuteFiats{}, "comdex-blockchain/MsgFactoryExecuteFiats", nil)
	cdc.RegisterConcrete(IssueFiatBody{}, "comdex-blockchain/IssueFiatBody", nil)
	cdc.RegisterConcrete(SendFiatBody{}, "comdex-blockchain/SendFiatBody", nil)
	cdc.RegisterConcrete(ExecuteFiatBody{}, "comdex-blockchain/ExecuteFiatBody", nil)
	cdc.RegisterConcrete(RedeemFiatBody{}, "comdex-blockchain/RedeemFiatBody", nil)
}

// RegisterFiatPeg : register concrete  and interface types
func RegisterFiatPeg(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.FiatPeg)(nil), nil)
	cdc.RegisterConcrete(&types.BaseFiatPeg{}, "comdex-blockchain/FiatPeg", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	RegisterFiatPeg(msgCdc)
}
