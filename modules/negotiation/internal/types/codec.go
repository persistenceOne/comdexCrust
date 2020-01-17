package types

import (
	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/types"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*types.Negotiation)(nil), nil)
	cdc.RegisterConcrete(&types.BaseNegotiation{}, "persistence-blockchain/Negotiation", nil)
	cdc.RegisterConcrete(MsgChangeBuyerBids{}, "persistence-blockchain/MsgChangeBuyerBids", nil)
	cdc.RegisterConcrete(MsgChangeSellerBids{}, "persistence-blockchain/MsgChangeSellerBids", nil)
	cdc.RegisterConcrete(MsgConfirmBuyerBids{}, "persistence-blockchain/MsgConfirmBuyerBids", nil)
	cdc.RegisterConcrete(MsgConfirmSellerBids{}, "persistence-blockchain/MsgConfirmSellerBids", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
