package fiatFactory

import (
	"github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
)

//RegisterWire : Register concrete types on wire codec for default fiats
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgFactoryIssueFiats{}, "commit-blockchain/MsgFactoryIssueFiats", nil)
	cdc.RegisterConcrete(MsgFactoryRedeemFiats{}, "commit-blockchain/MsgFactoryRedeemFiats", nil)
	cdc.RegisterConcrete(MsgFactorySendFiats{}, "commit-blockchain/MsgFactorySendFiats", nil)
	cdc.RegisterConcrete(MsgFactoryExecuteFiats{}, "commit-blockchain/MsgFactoryExecuteFiats", nil)
	cdc.RegisterConcrete(IssueFiatBody{}, "commit-blockchain/IssueFiatBody", nil)
	cdc.RegisterConcrete(SendFiatBody{}, "commit-blockchain/SendFiatBody", nil)
	cdc.RegisterConcrete(ExecuteFiatBody{}, "commit-blockchain/ExecuteFiatBody", nil)
	cdc.RegisterConcrete(RedeemFiatBody{}, "commit-blockchain/RedeemFiatBody", nil)
}

// RegisterFiatPeg : register concrete  and interface types
func RegisterFiatPeg(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.FiatPeg)(nil), nil)
	cdc.RegisterConcrete(&types.BaseFiatPeg{}, "commit-blockchain/FiatPeg", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	RegisterFiatPeg(msgCdc)
}
