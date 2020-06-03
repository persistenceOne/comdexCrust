package types

import (
	"github.com/commitHub/commitBlockchain/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Order)(nil), nil)
	cdc.RegisterConcrete(&BaseOrder{}, "commit-blockchain/Order", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
