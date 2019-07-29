package slashing

import (
	"github.com/commitHub/commitBlockchain/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgUnjail{}, "commit-blockchain/MsgUnjail", nil)
}

var cdcEmpty = wire.NewCodec()
