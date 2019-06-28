package reputation

import (
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// RegisterWire : Register concrete types on wire codec for order
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgBuyerFeedbacks{}, "comdex-blockchain/MsgBuyerFeedbacks", nil)
	cdc.RegisterConcrete(MsgSellerFeedbacks{}, "comdex-blockchain/MsgSellerFeedbacks", nil)
	cdc.RegisterConcrete(SubmitBuyerFeedbackBody{}, "comdex-blockchain/SubmitBuyerFeedbackBody", nil)
	cdc.RegisterConcrete(SubmitSellerFeedbackBody{}, "comdex-blockchain/SubmitSellerFeedbackBody", nil)
	
}

// RegisterReputation :
func RegisterReputation(cdc *wire.Codec) {
	cdc.RegisterInterface((*sdk.AccountReputation)(nil), nil)
	cdc.RegisterConcrete(&sdk.BaseAccountReputation{}, "comdex-blockchain/AccountReputation", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	RegisterReputation(msgCdc)
}
