package types

import (
	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/types"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*types.Order)(nil), nil)
	cdc.RegisterConcrete(&types.BaseOrder{}, "commit-blockchain/Order", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
