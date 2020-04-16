package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/comdexCrust/modules/bank"
	channelexported "github.com/persistenceOne/comdexCrust/modules/ibc/04-channel/exported"
	commitment "github.com/persistenceOne/comdexCrust/modules/ibc/23-commitment"
	host "github.com/persistenceOne/comdexCrust/modules/ibc/24-host"
	ibctypes "github.com/persistenceOne/comdexCrust/modules/ibc/types"
)

type MsgTransfer struct {
	SourcePort    string         `json:"source_port" yaml:"source_port"`       // the port on which the packet will be sent
	SourceChannel string         `json:"source_channel" yaml:"source_channel"` // the channel by which the packet will be sent
	Amount        sdk.Coins      `json:"amount" yaml:"amount"`                 // the tokens to be transferred
	Sender        sdk.AccAddress `json:"sender" yaml:"sender"`                 // the sender address
	Receiver      sdk.AccAddress `json:"receiver" yaml:"receiver"`             // the recipient address on the destination chain
	Source        bool           `json:"source" yaml:"source"`                 // indicates if the sending chain is the source chain of the tokens to be transferred
}

// NewMsgTransfer creates a new MsgTransfer instance
func NewMsgTransfer(
	sourcePort, sourceChannel string, amount sdk.Coins, sender, receiver sdk.AccAddress, source bool,
) MsgTransfer {
	return MsgTransfer{
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		Amount:        amount,
		Sender:        sender,
		Receiver:      receiver,
		Source:        source,
	}
}

// Route implements sdk.Msg
func (MsgTransfer) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (MsgTransfer) Type() string {
	return "transfer"
}

// ValidateBasic implements sdk.Msg
func (msg MsgTransfer) ValidateBasic() sdk.Error {
	if err := host.DefaultPortIdentifierValidator(msg.SourcePort); err != nil {
		return sdk.NewError(host.DefaultCodespace, 1, fmt.Sprintf("invalid port ID: %s", err.Error()))
	}
	if err := host.DefaultPortIdentifierValidator(msg.SourceChannel); err != nil {
		return sdk.NewError(host.DefaultCodespace, 1, fmt.Sprintf("invalid channel ID: %s", err.Error()))
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("transfer amount is invalid")
	}
	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("transfer amount must be positive")
	}
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if msg.Receiver.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

type MsgRecvPacket struct {
	Packet channelexported.PacketI `json:"packet" yaml:"packet"`
	Proofs []commitment.Proof      `json:"proofs" yaml:"proofs"`
	Height uint64                  `json:"height" yaml:"height"`
	Signer sdk.AccAddress          `json:"signer" yaml:"signer"`
}

// NewMsgRecvPacket creates a new MsgRecvPacket instance
func NewMsgRecvPacket(packet channelexported.PacketI, proofs []commitment.Proof, height uint64, signer sdk.AccAddress) MsgRecvPacket {
	return MsgRecvPacket{
		Packet: packet,
		Proofs: proofs,
		Height: height,
		Signer: signer,
	}
}

// Route implements sdk.Msg
func (MsgRecvPacket) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (MsgRecvPacket) Type() string {
	return "recv_packet"
}

// ValidateBasic implements sdk.Msg
func (msg MsgRecvPacket) ValidateBasic() sdk.Error {
	if msg.Height < 1 {
		return sdk.NewError(DefaultCodespace, CodeInvalidHeight, "invalid height")
	}

	if msg.Proofs == nil {
		return sdk.NewError(DefaultCodespace, CodeProofMissing, "proof missing")
	}

	for _, proof := range msg.Proofs {
		if proof.Proof == nil {
			return sdk.NewError(DefaultCodespace, CodeProofMissing, "proof missing")
		}
	}

	if msg.Signer.Empty() {
		return sdk.NewError(DefaultCodespace, CodeInvalidAddress, "invalid signer")
	}

	return msg.Packet.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (msg MsgRecvPacket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgRecvPacket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgIssueAssetTransfer struct {
	SourcePort    string          `json:"source_port" yaml:"source_port"`       // the port on which the packet will be sent
	SourceChannel string          `json:"source_channel" yaml:"source_channel"` // the channel by which the packet will be sent
	IssueAsset    bank.IssueAsset `json:"issueAsset" yaml:"issueAsset"`
}

// NewMsgIssueAssetTransfer creates a new MsgIssueAssetTransfer instance
func NewMsgIssueAssetTransfer(
	sourcePort, sourceChannel string, issueAsset bank.IssueAsset,
) MsgIssueAssetTransfer {
	return MsgIssueAssetTransfer{
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		IssueAsset:    issueAsset,
	}
}

// Route implements sdk.Msg
func (MsgIssueAssetTransfer) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (MsgIssueAssetTransfer) Type() string {
	return "transfer"
}

// ValidateBasic implements sdk.Msg
func (msg MsgIssueAssetTransfer) ValidateBasic() sdk.Error {
	if err := host.DefaultPortIdentifierValidator(msg.SourcePort); err != nil {
		return sdk.NewError(host.DefaultCodespace, 1, fmt.Sprintf("invalid port ID: %s", err.Error()))
	}
	if err := host.DefaultPortIdentifierValidator(msg.SourceChannel); err != nil {
		return sdk.NewError(host.DefaultCodespace, 1, fmt.Sprintf("invalid channel ID: %s", err.Error()))
	}
	if err := msg.IssueAsset.ValidateBasic(); err != nil {
		return sdk.ErrInvalidCoins("invalid asset")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgIssueAssetTransfer) GetSignBytes() []byte {
	// return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))

	b, err := ModuleCdc.MarshalJSON(struct {
		SourcePort    string          `json:"source_port"`
		SourceChannel string          `json:"source_channel"`
		IssueAsset    json.RawMessage `json:"issueAsset"`
	}{
		SourcePort:    msg.SourcePort,
		SourceChannel: msg.SourceChannel,
		IssueAsset:    msg.IssueAsset.GetSignBytes(),
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners implements sdk.Msg
func (msg MsgIssueAssetTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.IssueAsset.IssuerAddress}
}

type MsgReceiveIssueAssetPacket struct {
	Packet channelexported.PacketI `json:"packet" yaml:"packet"`
	Proofs []commitment.Proof      `json:"proofs" yaml:"proofs"`
	Height uint64                  `json:"height" yaml:"height"`
	Signer sdk.AccAddress          `json:"signer" yaml:"signer"`
}

// NewMsgReceiveIssueAssetPacket creates a new MsgReceiveIssueAssetPacket instance
func NewMsgReceiveIssueAssetPacket(packet channelexported.PacketI, proofs []commitment.Proof, height uint64, signer sdk.AccAddress) MsgReceiveIssueAssetPacket {
	return MsgReceiveIssueAssetPacket{
		Packet: packet,
		Proofs: proofs,
		Height: height,
		Signer: signer,
	}
}

// Route implements sdk.Msg
func (MsgReceiveIssueAssetPacket) Route() string {
	return ibctypes.RouterKey
}

// Type implements sdk.Msg
func (MsgReceiveIssueAssetPacket) Type() string {
	return "recv_packet"
}

// ValidateBasic implements sdk.Msg
func (msg MsgReceiveIssueAssetPacket) ValidateBasic() sdk.Error {
	if msg.Height < 1 {
		return sdk.NewError(DefaultCodespace, CodeInvalidHeight, "invalid height")
	}

	if msg.Proofs == nil {
		return sdk.NewError(DefaultCodespace, CodeProofMissing, "proof missing")
	}

	for _, proof := range msg.Proofs {
		if proof.Proof == nil {
			return sdk.NewError(DefaultCodespace, CodeProofMissing, "proof missing")
		}
	}

	if msg.Signer.Empty() {
		return sdk.NewError(DefaultCodespace, CodeInvalidAddress, "invalid signer")
	}

	return msg.Packet.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (msg MsgReceiveIssueAssetPacket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgReceiveIssueAssetPacket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
