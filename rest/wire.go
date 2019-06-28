package rest

import (
	"github.com/comdex-blockchain/wire"
)

var cdc = wire.NewCodec()

// RegisterWire registers structs to codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(KafkaTxCtx{}, "comdex-blockchain/KafkaTxCtx", nil)
	cdc.RegisterConcrete(KafkaCliCtx{}, "comdex-blockchain/KafkaCliCtx", nil)
	cdc.RegisterConcrete(KafkaMsg{}, "comdex-blockchain/KafkaMsg", nil)
}
func init() {
	RegisterWire(cdc)
}
