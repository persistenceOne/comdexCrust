package types

import (
	"encoding/json"
)

// Msg Transactions messages must fulfill the Msg
type Msg interface {
	// Return the message type.
	// Must be alphanumeric or empty.
	Type() string
	
	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() Error
	
	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte
	
	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigners() []AccAddress
}

// __________________________________________________________

// Tx Transactions objects must fulfill the Tx
type Tx interface {
	// Gets the Msg.
	GetMsgs() []Msg
}

// __________________________________________________________

// TxDecoder unmarshals transaction bytes
type TxDecoder func(txBytes []byte) (Tx, Error)

// __________________________________________________________

var _ Msg = (*TestMsg)(nil)

// TestMsg : msg type for testing
type TestMsg struct {
	signers []AccAddress
}

// NewTestMsg : returns new Test msg
func NewTestMsg(addrs ...AccAddress) *TestMsg {
	return &TestMsg{
		signers: addrs,
	}
}

// Type : type of test msg
// nolint
func (msg *TestMsg) Type() string { return "TestMsg" }

// GetSignBytes : gets sign byets of test msg
func (msg *TestMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg.signers)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(bz)
}

// ValidateBasic : valiadation for test msg
func (msg *TestMsg) ValidateBasic() Error { return nil }

// GetSigners : get signers of test msg
func (msg *TestMsg) GetSigners() []AccAddress {
	return msg.signers
}
