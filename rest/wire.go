package rest

import (
	"github.com/commitHub/commitBlockchain/wire"
)

var cdc = wire.NewCodec()

// RegisterWire registers structs to codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(KafkaTxCtx{}, "commit-blockchain/KafkaTxCtx", nil)
	cdc.RegisterConcrete(KafkaCliCtx{}, "commit-blockchain/KafkaCliCtx", nil)
	cdc.RegisterConcrete(KafkaMsg{}, "commit-blockchain/KafkaMsg", nil)
}
func init() {
	RegisterWire(cdc)
}
