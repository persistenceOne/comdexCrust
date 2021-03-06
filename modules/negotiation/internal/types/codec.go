package types

import "github.com/commitHub/commitBlockchain/codec"

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Negotiation)(nil), nil)
	cdc.RegisterConcrete(&BaseNegotiation{}, "commit-blockchain/Negotiation", nil)
	cdc.RegisterConcrete(MsgChangeBuyerBids{}, "commit-blockchain/MsgChangeBuyerBids", nil)
	cdc.RegisterConcrete(MsgChangeSellerBids{}, "commit-blockchain/MsgChangeSellerBids", nil)
	cdc.RegisterConcrete(MsgConfirmBuyerBids{}, "commit-blockchain/MsgConfirmBuyerBids", nil)
	cdc.RegisterConcrete(MsgConfirmSellerBids{}, "commit-blockchain/MsgConfirmSellerBids", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
