package gov

import (
	"github.com/commitHub/commitBlockchain/wire"
)

// RegisterWire : Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {

	cdc.RegisterConcrete(MsgSubmitProposal{}, "commit-blockchain/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "commit-blockchain/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "commit-blockchain/MsgVote", nil)

	cdc.RegisterInterface((*Proposal)(nil), nil)
	cdc.RegisterConcrete(&TextProposal{}, "gov/TextProposal", nil)
}

var msgCdc = wire.NewCodec()
