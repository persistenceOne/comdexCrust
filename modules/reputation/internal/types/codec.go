package types

import (
	"github.com/persistenceOne/persistenceSDK/codec"
)

func RegisterCodec(cdc *codec.Codec) {
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
