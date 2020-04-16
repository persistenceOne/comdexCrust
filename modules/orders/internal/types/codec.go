package types

import (
	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/types"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*types.Order)(nil), nil)
	cdc.RegisterConcrete(&types.BaseOrder{}, "persistence-blockchain/Order", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
