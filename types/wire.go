package types

import "github.com/comdex-blockchain/wire"

// RegisterWire : Register the sdk message type
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterInterface((*Tx)(nil), nil)
}
