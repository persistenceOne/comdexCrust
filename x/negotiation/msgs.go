package negotiation

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"

	sdk "github.com/commitHub/commitBlockchain/types"
)

//ChangeBid - change negotiation bid
type ChangeBid struct {
	Negotiation sdk.Negotiation `json:"negotiation" valid:"required"`
}

//NewChangeBid : initializer
func NewChangeBid(negotiation sdk.Negotiation) ChangeBid {
	return ChangeBid{negotiation}
}

//GetSignBytes : get bytes to sign
func (in ChangeBid) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		Negotiation sdk.Negotiation `json:"negotiation"`
	}{
		Negotiation: in.Negotiation,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//MsgChangeBuyerBids : high level change bid of negotiation module
type MsgChangeBuyerBids struct {
	ChangeBids []ChangeBid `json:"changeBids"`
}

//NewMsgChangeBuyerBids : initilizer
func NewMsgChangeBuyerBids(changeBids []ChangeBid) MsgChangeBuyerBids {
	return MsgChangeBuyerBids{changeBids}
}

var _ sdk.Msg = MsgChangeBuyerBids{}

//Type : implements msg
func (msg MsgChangeBuyerBids) Type() string { return "negotiation" }

//ValidateBasic : implements msg
func (msg MsgChangeBuyerBids) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.ChangeBids {
		_, err := govalidator.ValidateStruct(in)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
		if len(in.Negotiation.GetBuyerAddress()) == 0 {
			return sdk.ErrInvalidAddress(in.Negotiation.GetBuyerAddress().String())
		} else if len(in.Negotiation.GetSellerAddress()) == 0 {
			return sdk.ErrInvalidAddress(in.Negotiation.GetSellerAddress().String())
		} else if len(in.Negotiation.GetNegotiationID()) == 0 {
			return sdk.ErrUnknownRequest("Negotiation ID is wrong")
		} else if len(in.Negotiation.GetPegHash()) == 0 {
			return sdk.ErrUnknownRequest("Peghash is empty")
		} else if in.Negotiation.GetBid() <= 0 {
			return sdk.ErrUnknownRequest("Bid amount should be greater than 0")
		} else if in.Negotiation.GetTime() <= 0 {
			return sdk.ErrUnknownRequest("Time should not be 0")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgChangeBuyerBids) GetSignBytes() []byte {
	var changeBids []json.RawMessage
	for _, changeBid := range msg.ChangeBids {
		changeBids = append(changeBids, changeBid.GetSignBytes())
	}

	b, err := msgCdc.MarshalJSON(struct {
		ChangeBids []json.RawMessage `json:"changeBids"`
	}{
		ChangeBids: changeBids,
	})
	if err != nil {
		panic(err)
	}
	return b
}

//GetSigners : implements msg
func (msg MsgChangeBuyerBids) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.ChangeBids))
	for i, in := range msg.ChangeBids {
		addrs[i] = in.Negotiation.GetBuyerAddress()
	}
	return addrs
}

//BuildMsgChangeBuyerBid : build the MsgChangeBuyerBids
func BuildMsgChangeBuyerBid(negotiation sdk.Negotiation) sdk.Msg {
	changeBid := NewChangeBid(negotiation)
	msg := NewMsgChangeBuyerBids([]ChangeBid{changeBid})
	return msg
}

//MsgChangeSellerBids : high level change bid of negotiation module
type MsgChangeSellerBids struct {
	ChangeBids []ChangeBid `json:"changeBids"`
}

//NewMsgChangeSellerBids : initilizer
func NewMsgChangeSellerBids(changeBids []ChangeBid) MsgChangeSellerBids {
	return MsgChangeSellerBids{changeBids}
}

var _ sdk.Msg = MsgChangeSellerBids{}

//Type : implements msg
func (msg MsgChangeSellerBids) Type() string { return "negotiation" }

//ValidateBasic : implements msg
func (msg MsgChangeSellerBids) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.ChangeBids {
		_, err := govalidator.ValidateStruct(in)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
		if len(in.Negotiation.GetBuyerAddress()) == 0 {
			return sdk.ErrInvalidAddress(in.Negotiation.GetBuyerAddress().String())
		} else if len(in.Negotiation.GetSellerAddress()) == 0 {
			return sdk.ErrInvalidAddress(in.Negotiation.GetSellerAddress().String())
		} else if len(in.Negotiation.GetNegotiationID()) == 0 {
			return sdk.ErrUnknownRequest("Negotiation ID is wrong")
		} else if len(in.Negotiation.GetPegHash()) == 0 {
			return sdk.ErrUnknownRequest("Peghash is empty")
		} else if in.Negotiation.GetBid() <= 0 {
			return sdk.ErrUnknownRequest("Bid amount should be greater than 0")
		} else if in.Negotiation.GetTime() <= 0 {
			return sdk.ErrUnknownRequest("Time should not be 0")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgChangeSellerBids) GetSignBytes() []byte {
	var changeBids []json.RawMessage
	for _, changeBid := range msg.ChangeBids {
		changeBids = append(changeBids, changeBid.GetSignBytes())
	}

	b, err := msgCdc.MarshalJSON(struct {
		ChangeBids []json.RawMessage `json:"changeBids"`
	}{
		ChangeBids: changeBids,
	})
	if err != nil {
		panic(err)
	}
	return b
}

//GetSigners : implements msg
func (msg MsgChangeSellerBids) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.ChangeBids))
	for i, in := range msg.ChangeBids {
		addrs[i] = in.Negotiation.GetSellerAddress()
	}
	return addrs
}

//BuildMsgChangeSellerBid : build the MsgChangeSellerBids
func BuildMsgChangeSellerBid(negotiation sdk.Negotiation) sdk.Msg {
	changeBid := NewChangeBid(negotiation)
	msg := NewMsgChangeSellerBids([]ChangeBid{changeBid})
	return msg
}

// ConfirmBid :
type ConfirmBid struct {
	Negotiation sdk.Negotiation `json:"negotiation" valid:"required"`
}

//NewConfirmBid : initializer
func NewConfirmBid(negotiation sdk.Negotiation) ConfirmBid {
	return ConfirmBid{negotiation}
}

//GetSignBytes : get bytes to sign
func (in ConfirmBid) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		Negotiation sdk.Negotiation `json:"negotiation"`
	}{
		Negotiation: in.Negotiation,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//MsgConfirmBuyerBids :
type MsgConfirmBuyerBids struct {
	ConfirmBids []ConfirmBid `json:"confirmBids"`
}

//NewMsgConfirmBuyerBids : initializer
func NewMsgConfirmBuyerBids(confirmBeds []ConfirmBid) MsgConfirmBuyerBids {
	return MsgConfirmBuyerBids{confirmBeds}
}

var _ sdk.Msg = MsgConfirmBuyerBids{}

//Type : implements msg
func (msg MsgConfirmBuyerBids) Type() string { return "negotiation" }

//ValidateBasic : implements msg
func (msg MsgConfirmBuyerBids) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.ConfirmBids {
		_, err := govalidator.ValidateStruct(in)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
		if len(in.Negotiation.GetBuyerAddress()) == 0 {
			return sdk.ErrInvalidAddress(in.Negotiation.GetBuyerAddress().String())
		} else if len(in.Negotiation.GetSellerAddress()) == 0 {
			return sdk.ErrInvalidAddress(in.Negotiation.GetSellerAddress().String())
		} else if len(in.Negotiation.GetNegotiationID()) == 0 {
			return sdk.ErrUnknownRequest("Negotiation ID is wrong")
		} else if len(in.Negotiation.GetPegHash()) == 0 {
			return sdk.ErrUnknownRequest("Peghash is empty")
		} else if in.Negotiation.GetBid() <= 0 {
			return sdk.ErrUnknownRequest("Bid amount should be greater than 0")
		} else if in.Negotiation.GetTime() <= 0 {
			return sdk.ErrUnknownRequest("Time should not be 0")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgConfirmBuyerBids) GetSignBytes() []byte {
	var confirmBids []json.RawMessage
	for _, confirmBid := range msg.ConfirmBids {
		confirmBids = append(confirmBids, confirmBid.GetSignBytes())
	}

	b, err := msgCdc.MarshalJSON(struct {
		ConfirmBids []json.RawMessage `json:"confirmBids"`
	}{
		ConfirmBids: confirmBids,
	})
	if err != nil {
		panic(err)
	}
	return b
}

//GetSigners : implements msg
func (msg MsgConfirmBuyerBids) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.ConfirmBids))
	for i, in := range msg.ConfirmBids {
		addrs[i] = in.Negotiation.GetBuyerAddress()
	}
	return addrs
}

//BuildMsgConfirmBuyerBid : build the MsgConfirmBuyerBids
func BuildMsgConfirmBuyerBid(negotiation sdk.Negotiation) sdk.Msg {
	confirmBid := NewConfirmBid(negotiation)
	msg := NewMsgConfirmBuyerBids([]ConfirmBid{confirmBid})
	return msg
}

//******MsgConfirmBuyerBids

//#######MsgConfirmSellerBids

//MsgConfirmSellerBids :
type MsgConfirmSellerBids struct {
	ConfirmBids []ConfirmBid `json:"confirmBids"`
}

//NewMsgConfirmSellerBids : initializer
func NewMsgConfirmSellerBids(confirmBids []ConfirmBid) MsgConfirmSellerBids {
	return MsgConfirmSellerBids{confirmBids}
}

var _ sdk.Msg = MsgConfirmSellerBids{}

//Type : implements msg
func (msg MsgConfirmSellerBids) Type() string { return "negotiation" }

//ValidateBasic : implements msg
func (msg MsgConfirmSellerBids) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.ConfirmBids {
		_, err := govalidator.ValidateStruct(in)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
		if len(in.Negotiation.GetBuyerAddress()) == 0 {
			return sdk.ErrInvalidAddress(in.Negotiation.GetBuyerAddress().String())
		} else if len(in.Negotiation.GetSellerAddress()) == 0 {
			return sdk.ErrInvalidAddress(in.Negotiation.GetSellerAddress().String())
		} else if len(in.Negotiation.GetNegotiationID()) == 0 {
			return sdk.ErrUnknownRequest("Negotiation ID is wrong")
		} else if len(in.Negotiation.GetPegHash()) == 0 {
			return sdk.ErrUnknownRequest("Peghash is empty")
		} else if in.Negotiation.GetBid() <= 0 {
			return sdk.ErrUnknownRequest("Bid amount should be greater than 0")
		} else if in.Negotiation.GetTime() <= 0 {
			return sdk.ErrUnknownRequest("Time should not be 0")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgConfirmSellerBids) GetSignBytes() []byte {
	var confirmBids []json.RawMessage
	for _, confirmBid := range msg.ConfirmBids {
		confirmBids = append(confirmBids, confirmBid.GetSignBytes())
	}

	b, err := msgCdc.MarshalJSON(struct {
		ConfirmBids []json.RawMessage `json:"confirmBids"`
	}{
		ConfirmBids: confirmBids,
	})
	if err != nil {
		panic(err)
	}
	return b
}

//GetSigners : implements msg
func (msg MsgConfirmSellerBids) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.ConfirmBids))
	for i, in := range msg.ConfirmBids {
		addrs[i] = in.Negotiation.GetSellerAddress()
	}
	return addrs
}

//BuildMsgConfirmSellerBid : build the MsgConfirmBuyerBids
func BuildMsgConfirmSellerBid(negotiation sdk.Negotiation) sdk.Msg {
	confirmBid := NewConfirmBid(negotiation)
	msg := NewMsgConfirmSellerBids([]ConfirmBid{confirmBid})
	return msg
}

//#####MsgSellerBids

//******SignNegotiationBody

// SignNegotiationBody :
type SignNegotiationBody struct {
	BuyerAddress  sdk.AccAddress `json:"buyer_address"`
	SellerAddress sdk.AccAddress `json:"seller_address"`
	PegHash       sdk.PegHash    `json:"peg_hash"`
	Bid           int64          `json:"bid"`
	Time          int64          `json:"time"`
}

// NewSignNegotiationBody :
func NewSignNegotiationBody(buyerAddress, sellerAddress sdk.AccAddress, peghash sdk.PegHash, bid, time int64) *SignNegotiationBody {
	return &SignNegotiationBody{
		BuyerAddress:  buyerAddress,
		SellerAddress: sellerAddress,
		PegHash:       peghash,
		Bid:           bid,
		Time:          time,
	}
}

//GetSignBytes :
func (bytes SignNegotiationBody) GetSignBytes() []byte {
	bz, err := msgCdc.MarshalJSON(bytes)
	if err != nil {
		panic(err)
	}
	return bz
}

//#########SignNegotiationBody
