package types

import "github.com/commitHub/commitBlockchain/codec"

func RegisterCodec(cdc *codec.Codec) {}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
