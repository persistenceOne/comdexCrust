package genaccounts

import (
	"github.com/persistenceOne/comdexCrust/codec"
)

// generic sealed codec to be used throughout this module
var moduleCdc *codec.Codec

func init() {
	moduleCdc = codec.New()
	codec.RegisterCrypto(moduleCdc)
	moduleCdc.Seal()
}
