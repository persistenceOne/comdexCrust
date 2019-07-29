package order

import (
	"github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
)

//RegisterWire : Register concrete types on wire codec for order
func RegisterWire(cdc *wire.Codec) {
}

//RegisterOrder : registe order an baseOrder from types for order module
func RegisterOrder(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.Order)(nil), nil)
	cdc.RegisterConcrete(&types.BaseOrder{}, "commit-blockchain/Order", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	RegisterOrder(msgCdc)
	wire.RegisterCrypto(msgCdc)
}
