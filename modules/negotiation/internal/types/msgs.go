package types

import (
	"encoding/json"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/types"
)

// ChangeBid - change negotiation bid
type ChangeBid struct {
	Negotiation types.Negotiation `json:"negotiation" valid:"required"`
}

// NewChangeBid : initializer
func NewChangeBid(negotiation types.Negotiation) ChangeBid {
	return ChangeBid{negotiation}
}

// GetSignBytes : get bytes to sign
func (in ChangeBid) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		Negotiation types.Negotiation `json:"negotiation"`
	}{
		Negotiation: in.Negotiation,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

func (in ChangeBid) ValidateBasic() cTypes.Error {
	if len(in.Negotiation.GetBuyerAddress()) == 0 {
		return cTypes.ErrInvalidAddress(in.Negotiation.GetBuyerAddress().String())
	} else if len(in.Negotiation.GetSellerAddress()) == 0 {
		return cTypes.ErrInvalidAddress(in.Negotiation.GetSellerAddress().String())
	} else if len(in.Negotiation.GetPegHash()) == 0 {
		return cTypes.ErrUnknownRequest("PegHash should not be empty.")
	} else if in.Negotiation.GetBid() < 0 {
		return ErrNegativeAmount(DefaultCodeSpace, "Bid should not e negative.")
	}
	return nil
}

// MsgChangeBuyerBids : high level change bid of negotiation module
type MsgChangeBuyerBids struct {
	ChangeBids []ChangeBid `json:"changeBids"`
}

// NewMsgChangeBuyerBids : initializer
func NewMsgChangeBuyerBids(changeBids []ChangeBid) MsgChangeBuyerBids {
	return MsgChangeBuyerBids{changeBids}
}

var _ cTypes.Msg = MsgChangeBuyerBids{}

// Type : implements msg
func (msg MsgChangeBuyerBids) Type() string { return "changeBuyerBids" }

// ValidateBasic : implements msg
func (msg MsgChangeBuyerBids) ValidateBasic() cTypes.Error {
	if len(msg.ChangeBids) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.ChangeBids {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgChangeBuyerBids) GetSignBytes() []byte {
	var changeBids []json.RawMessage
	for _, changeBid := range msg.ChangeBids {
		changeBids = append(changeBids, changeBid.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		ChangeBids []json.RawMessage `json:"changeBids"`
	}{
		ChangeBids: changeBids,
	})
	if err != nil {
		panic(err)

	}
	return b
}

// GetSigners : implements msg
func (msg MsgChangeBuyerBids) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.ChangeBids))
	for i, in := range msg.ChangeBids {
		addrs[i] = in.Negotiation.GetBuyerAddress()
	}
	return addrs
}

func (msg MsgChangeBuyerBids) Route() string { return RouterKey }

// BuildMsgChangeBuyerBid : build the MsgChangeBuyerBids
func BuildMsgChangeBuyerBid(negotiation types.Negotiation) cTypes.Msg {
	changeBid := NewChangeBid(negotiation)
	msg := NewMsgChangeBuyerBids([]ChangeBid{changeBid})
	return msg
}

// MsgChangeSellerBids : high level change bid of negotiation module
type MsgChangeSellerBids struct {
	ChangeBids []ChangeBid `json:"changeBids"`
}

// NewMsgChangeSellerBids : initilizer
func NewMsgChangeSellerBids(changeBids []ChangeBid) MsgChangeSellerBids {
	return MsgChangeSellerBids{changeBids}
}

var _ cTypes.Msg = MsgChangeSellerBids{}

// Type : implements msg
func (msg MsgChangeSellerBids) Type() string { return "changeSellerBids" }

// ValidateBasic : implements msg
func (msg MsgChangeSellerBids) ValidateBasic() cTypes.Error {
	if len(msg.ChangeBids) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.ChangeBids {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgChangeSellerBids) GetSignBytes() []byte {
	var changeBids []json.RawMessage
	for _, changeBid := range msg.ChangeBids {
		changeBids = append(changeBids, changeBid.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		ChangeBids []json.RawMessage `json:"changeBids"`
	}{
		ChangeBids: changeBids,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgChangeSellerBids) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.ChangeBids))
	for i, in := range msg.ChangeBids {
		addrs[i] = in.Negotiation.GetSellerAddress()
	}
	return addrs
}

func (msg MsgChangeSellerBids) Route() string {
	return RouterKey
}

// BuildMsgChangeSellerBid : build the MsgChangeSellerBids
func BuildMsgChangeSellerBid(negotiation types.Negotiation) cTypes.Msg {
	changeBid := NewChangeBid(negotiation)
	msg := NewMsgChangeSellerBids([]ChangeBid{changeBid})
	return msg
}

// ConfirmBid :
type ConfirmBid struct {
	Negotiation types.Negotiation `json:"negotiation" valid:"required"`
}

// NewConfirmBid : initializer
func NewConfirmBid(negotiation types.Negotiation) ConfirmBid {
	return ConfirmBid{negotiation}
}

// GetSignBytes : get bytes to sign
func (in ConfirmBid) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		Negotiation types.Negotiation `json:"negotiation"`
	}{
		Negotiation: in.Negotiation,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

func (in ConfirmBid) ValidateBasic() cTypes.Error {
	if len(in.Negotiation.GetBuyerAddress()) == 0 {
		return cTypes.ErrInvalidAddress(in.Negotiation.GetBuyerAddress().String())
	} else if len(in.Negotiation.GetSellerAddress()) == 0 {
		return cTypes.ErrInvalidAddress(in.Negotiation.GetSellerAddress().String())
	} else if len(in.Negotiation.GetPegHash()) == 0 {
		return cTypes.ErrUnknownRequest("PegHash should not be empty.")
	} else if in.Negotiation.GetBid() < 0 {
		return ErrNegativeAmount(DefaultCodeSpace, "Bid should not e negative.")
	}
	return nil
}

// MsgConfirmBuyerBids :
type MsgConfirmBuyerBids struct {
	ConfirmBids []ConfirmBid `json:"confirmBids"`
}

// NewMsgConfirmBuyerBids : initializer
func NewMsgConfirmBuyerBids(confirmBeds []ConfirmBid) MsgConfirmBuyerBids {
	return MsgConfirmBuyerBids{confirmBeds}
}

var _ cTypes.Msg = MsgConfirmBuyerBids{}

// Type : implements msg
func (msg MsgConfirmBuyerBids) Type() string { return "confirmBuyerBids" }

// ValidateBasic : implements msg
func (msg MsgConfirmBuyerBids) ValidateBasic() cTypes.Error {
	if len(msg.ConfirmBids) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.ConfirmBids {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgConfirmBuyerBids) GetSignBytes() []byte {
	var confirmBids []json.RawMessage
	for _, confirmBid := range msg.ConfirmBids {
		confirmBids = append(confirmBids, confirmBid.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		ConfirmBids []json.RawMessage `json:"confirmBids"`
	}{
		ConfirmBids: confirmBids,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgConfirmBuyerBids) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.ConfirmBids))
	for i, in := range msg.ConfirmBids {
		addrs[i] = in.Negotiation.GetBuyerAddress()
	}
	return addrs
}

func (msg MsgConfirmBuyerBids) Route() string {
	return RouterKey
}

// BuildMsgConfirmBuyerBid : build the MsgConfirmBuyerBids
func BuildMsgConfirmBuyerBid(negotiation types.Negotiation) cTypes.Msg {
	confirmBid := NewConfirmBid(negotiation)
	msg := NewMsgConfirmBuyerBids([]ConfirmBid{confirmBid})
	return msg
}

// ******MsgConfirmBuyerBids

// #######MsgConfirmSellerBids

// MsgConfirmSellerBids :
type MsgConfirmSellerBids struct {
	ConfirmBids []ConfirmBid `json:"confirmBids"`
}

// NewMsgConfirmSellerBids : initializer
func NewMsgConfirmSellerBids(confirmBids []ConfirmBid) MsgConfirmSellerBids {
	return MsgConfirmSellerBids{confirmBids}
}

var _ cTypes.Msg = MsgConfirmSellerBids{}

// Type : implements msg
func (msg MsgConfirmSellerBids) Type() string { return "confirmSellerBids" }

// ValidateBasic : implements msg
func (msg MsgConfirmSellerBids) ValidateBasic() cTypes.Error {
	if len(msg.ConfirmBids) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.ConfirmBids {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgConfirmSellerBids) GetSignBytes() []byte {
	var confirmBids []json.RawMessage
	for _, confirmBid := range msg.ConfirmBids {
		confirmBids = append(confirmBids, confirmBid.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		ConfirmBids []json.RawMessage `json:"confirmBids"`
	}{
		ConfirmBids: confirmBids,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgConfirmSellerBids) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.ConfirmBids))
	for i, in := range msg.ConfirmBids {
		addrs[i] = in.Negotiation.GetSellerAddress()
	}
	return addrs
}

func (msg MsgConfirmSellerBids) Route() string {
	return RouterKey
}

// BuildMsgConfirmSellerBid : build the MsgConfirmBuyerBids
func BuildMsgConfirmSellerBid(negotiation types.Negotiation) cTypes.Msg {
	confirmBid := NewConfirmBid(negotiation)
	msg := NewMsgConfirmSellerBids([]ConfirmBid{confirmBid})
	return msg
}

// SignNegotiationBody :
type SignNegotiationBody struct {
	BuyerAddress  cTypes.AccAddress `json:"buyerAddress"`
	SellerAddress cTypes.AccAddress `json:"sellerAddress"`
	PegHash       types.PegHash     `json:"pegHash"`
	Bid           int64             `json:"bid"`
	Time          int64             `json:"time"`
}

// NewSignNegotiationBody :
func NewSignNegotiationBody(buyerAddress, sellerAddress cTypes.AccAddress, peghash types.PegHash, bid, time int64) *SignNegotiationBody {
	return &SignNegotiationBody{
		BuyerAddress:  buyerAddress,
		SellerAddress: sellerAddress,
		PegHash:       peghash,
		Bid:           bid,
		Time:          time,
	}
}

// GetSignBytes :
func (bytes SignNegotiationBody) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(bytes)
	if err != nil {
		panic(err)
	}
	return bz
}

// #####MsgSellerBids
