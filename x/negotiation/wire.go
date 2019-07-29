package negotiation

import (
	"github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
)

//RegisterWire : Register concrete types on wire codec for negotiation
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgChangeBuyerBids{}, "commit-blockchain/MsgChangeBuyerBids", nil)
	cdc.RegisterConcrete(MsgChangeSellerBids{}, "commit-blockchain/MsgChangeSellerBids", nil)
	cdc.RegisterConcrete(MsgConfirmBuyerBids{}, "commit-blockchain/MsgConfirmBuyerBids", nil)
	cdc.RegisterConcrete(MsgConfirmSellerBids{}, "commit-blockchain/MsgConfirmSellerBids", nil)
	cdc.RegisterConcrete(ChangeBuyerBidBody{}, "commit-blockchain/ChangeBuyerBidBody", nil)
	cdc.RegisterConcrete(ChangeSellerBidBody{}, "commit-blockchain/ChangeSellerBidBody", nil)
	cdc.RegisterConcrete(ConfirmBuyerBidBody{}, "commit-blockchain/ConfirmBuyerBidBody", nil)
	cdc.RegisterConcrete(ConfirmSellerBidBody{}, "commit-blockchain/ConfirmSellerBidBody", nil)
	cdc.RegisterConcrete(NegotiaitonBody{}, "commit-blockchain/NegotiaitonBody", nil)

}

// RegisterNegotiation :
func RegisterNegotiation(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.Negotiation)(nil), nil)
	cdc.RegisterConcrete(&types.BaseNegotiation{}, "commit-blockchain/Negotiation", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	RegisterNegotiation(msgCdc)
	wire.RegisterCrypto(msgCdc)
}
