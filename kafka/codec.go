package kafka

import (
	"github.com/persistenceOne/persistenceSDK/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(KafkaCliCtx{}, "commit-blockchain/KafkaCliCtx", nil)
	cdc.RegisterConcrete(KafkaMsg{}, "commit-blockchain/KafkaMsg", nil)
}

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
