package types

import (
	"github.com/comdex-blockchain/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "comdex-blockchain/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "comdex-blockchain/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "comdex-blockchain/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgBeginUnbonding{}, "comdex-blockchain/BeginUnbonding", nil)
	cdc.RegisterConcrete(MsgCompleteUnbonding{}, "comdex-blockchain/CompleteUnbonding", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "comdex-blockchain/BeginRedelegate", nil)
	cdc.RegisterConcrete(MsgCompleteRedelegate{}, "comdex-blockchain/CompleteRedelegate", nil)
}

// generic sealed codec to be used throughout sdk
var MsgCdc *wire.Codec

func init() {
	cdc := wire.NewCodec()
	RegisterWire(cdc)
	wire.RegisterCrypto(cdc)
	MsgCdc = cdc
	// MsgCdc = cdc.Seal() //TODO use when upgraded to go-amino 0.9.10
}
