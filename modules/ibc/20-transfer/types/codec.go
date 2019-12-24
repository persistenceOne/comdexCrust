package types

import (
	channel "github.com/commitHub/commitBlockchain/modules/ibc/04-channel"
	commitment "github.com/commitHub/commitBlockchain/modules/ibc/23-commitment"
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTransfer{}, "ibc/transfer/MsgTransfer", nil)
	cdc.RegisterConcrete(MsgRecvPacket{}, "ibc/transfer/MsgRecvPacket", nil)
	cdc.RegisterConcrete(PacketData{}, "ibc/transfer/PacketData", nil)

	cdc.RegisterConcrete(MsgIssueAssetTransfer{}, "ibc/transfer/MsgIssueAssetTransfer", nil)
	cdc.RegisterConcrete(MsgReceiveIssueAssetPacket{}, "ibc/transfer/MsgReceiveIssueAssetPacket", nil)
	cdc.RegisterConcrete(IssueAssetPacketData{}, "ibc/transfer/IssueAssetPacketData", nil)

}

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	channel.RegisterCodec(ModuleCdc)
	commitment.RegisterCodec(ModuleCdc)
}
