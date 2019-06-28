package gov

import (
	"github.com/comdex-blockchain/wire"
)

// RegisterWire : Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	
	cdc.RegisterConcrete(MsgSubmitProposal{}, "comdex-blockchain/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "comdex-blockchain/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "comdex-blockchain/MsgVote", nil)
	
	cdc.RegisterInterface((*Proposal)(nil), nil)
	cdc.RegisterConcrete(&TextProposal{}, "gov/TextProposal", nil)
}

var msgCdc = wire.NewCodec()
