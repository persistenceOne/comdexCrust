package types

import (
	"github.com/commitHub/commitBlockchain/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "commit-blockchain/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "commit-blockchain/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "commit-blockchain/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgBeginUnbonding{}, "commit-blockchain/BeginUnbonding", nil)
	cdc.RegisterConcrete(MsgCompleteUnbonding{}, "commit-blockchain/CompleteUnbonding", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "commit-blockchain/BeginRedelegate", nil)
	cdc.RegisterConcrete(MsgCompleteRedelegate{}, "commit-blockchain/CompleteRedelegate", nil)
}

// generic sealed codec to be used throughout sdk
var MsgCdc *wire.Codec

func init() {
	cdc := wire.NewCodec()
	RegisterWire(cdc)
	wire.RegisterCrypto(cdc)
	MsgCdc = cdc
	//MsgCdc = cdc.Seal() //TODO use when upgraded to go-amino 0.9.10
}
