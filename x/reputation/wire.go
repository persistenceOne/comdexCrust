package reputation

import (
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
)

//RegisterWire : Register concrete types on wire codec for order
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgBuyerFeedbacks{}, "commit-blockchain/MsgBuyerFeedbacks", nil)
	cdc.RegisterConcrete(MsgSellerFeedbacks{}, "commit-blockchain/MsgSellerFeedbacks", nil)
	cdc.RegisterConcrete(SubmitBuyerFeedbackBody{}, "commit-blockchain/SubmitBuyerFeedbackBody", nil)
	cdc.RegisterConcrete(SubmitSellerFeedbackBody{}, "commit-blockchain/SubmitSellerFeedbackBody", nil)

}

//RegisterReputation :
func RegisterReputation(cdc *wire.Codec) {
	cdc.RegisterInterface((*sdk.AccountReputation)(nil), nil)
	cdc.RegisterConcrete(&sdk.BaseAccountReputation{}, "commit-blockchain/AccountReputation", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	RegisterReputation(msgCdc)
}
