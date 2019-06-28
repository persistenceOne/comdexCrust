package negotiation

import (
	"github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// RegisterWire : Register concrete types on wire codec for negotiation
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgChangeBuyerBids{}, "comdex-blockchain/MsgChangeBuyerBids", nil)
	cdc.RegisterConcrete(MsgChangeSellerBids{}, "comdex-blockchain/MsgChangeSellerBids", nil)
	cdc.RegisterConcrete(MsgConfirmBuyerBids{}, "comdex-blockchain/MsgConfirmBuyerBids", nil)
	cdc.RegisterConcrete(MsgConfirmSellerBids{}, "comdex-blockchain/MsgConfirmSellerBids", nil)
	cdc.RegisterConcrete(ChangeBuyerBidBody{}, "comdex-blockchain/ChangeBuyerBidBody", nil)
	cdc.RegisterConcrete(ChangeSellerBidBody{}, "comdex-blockchain/ChangeSellerBidBody", nil)
	cdc.RegisterConcrete(ConfirmBuyerBidBody{}, "comdex-blockchain/ConfirmBuyerBidBody", nil)
	cdc.RegisterConcrete(ConfirmSellerBidBody{}, "comdex-blockchain/ConfirmSellerBidBody", nil)
	cdc.RegisterConcrete(NegotiaitonBody{}, "comdex-blockchain/NegotiaitonBody", nil)
	
}

// RegisterNegotiation :
func RegisterNegotiation(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.Negotiation)(nil), nil)
	cdc.RegisterConcrete(&types.BaseNegotiation{}, "comdex-blockchain/Negotiation", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	RegisterNegotiation(msgCdc)
	wire.RegisterCrypto(msgCdc)
}
