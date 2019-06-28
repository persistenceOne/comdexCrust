package reputation

import (
	"encoding/json"
	
	"github.com/asaskevich/govalidator"
	sdk "github.com/comdex-blockchain/types"
)

// SubmitTraderFeedback : TraderFeedback
type SubmitTraderFeedback struct {
	TraderFeedback sdk.TraderFeedback `json:"traderfeedback"`
}

// NewSubmitTraderFeedback : creates new
func NewSubmitTraderFeedback(traderFeedback sdk.TraderFeedback) SubmitTraderFeedback {
	return SubmitTraderFeedback{TraderFeedback: traderFeedback}
}

// GetSignBytes : get bytes to sign
func (submitTraderFeedback SubmitTraderFeedback) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		TraderFeedback sdk.TraderFeedback `json:"traderfeedback"`
	}{
		TraderFeedback: submitTraderFeedback.TraderFeedback,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// MsgBuyerFeedbacks : Msg for   traderFeedback
type MsgBuyerFeedbacks struct {
	SubmitTraderFeedbacks []SubmitTraderFeedback `json:"submitTraderFeedbacks"`
}

var _ sdk.Msg = MsgBuyerFeedbacks{}

// NewMsgBuyerFeedbacks : creates msg, buyer rates seller
func NewMsgBuyerFeedbacks(submitTraderFeedbacks []SubmitTraderFeedback) MsgBuyerFeedbacks {
	return MsgBuyerFeedbacks{SubmitTraderFeedbacks: submitTraderFeedbacks}
}

// Type : implements msg
func (msg MsgBuyerFeedbacks) Type() string { return "reputation" }

// ValidateBasic : implements msg
func (msg MsgBuyerFeedbacks) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.SubmitTraderFeedbacks {
		_, err := govalidator.ValidateStruct(in)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBuyerFeedbacks) GetSignBytes() []byte {
	var submitTraderFeedbacks []json.RawMessage
	for _, submitTraderFeedback := range msg.SubmitTraderFeedbacks {
		submitTraderFeedbacks = append(submitTraderFeedbacks, submitTraderFeedback.GetSignBytes())
	}
	
	b, err := msgCdc.MarshalJSON(struct {
		SubmitTraderFeedbacks []json.RawMessage `json:"submitTraderFeedbacks"`
	}{
		SubmitTraderFeedbacks: submitTraderFeedbacks,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgBuyerFeedbacks) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SubmitTraderFeedbacks))
	for i, r := range msg.SubmitTraderFeedbacks {
		addrs[i] = r.TraderFeedback.BuyerAddress
	}
	return addrs
}

// BuildBuyerFeedbackMsg : butild the FeedbackTx
func BuildBuyerFeedbackMsg(buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, score int64) sdk.Msg {
	traderFeedback := sdk.NewTraderFeedback(buyerAddress, sellerAddress, pegHash, score)
	submitTraderFeedback := NewSubmitTraderFeedback(traderFeedback)
	msg := NewMsgBuyerFeedbacks([]SubmitTraderFeedback{submitTraderFeedback})
	return &msg
}

// MsgSellerFeedbacks : Msg for   traderFeedback
type MsgSellerFeedbacks struct {
	SubmitTraderFeedbacks []SubmitTraderFeedback `json:" submitTraderFeedbacks"`
}

var _ sdk.Msg = MsgSellerFeedbacks{}

// NewMsgSellerFeedbacks : creates msg, Seller rates Buyer
func NewMsgSellerFeedbacks(submitTraderFeedbacks []SubmitTraderFeedback) MsgSellerFeedbacks {
	return MsgSellerFeedbacks{SubmitTraderFeedbacks: submitTraderFeedbacks}
}

// Type : implements msg
func (msg MsgSellerFeedbacks) Type() string { return "reputation" }

// ValidateBasic : implements msg
func (msg MsgSellerFeedbacks) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.SubmitTraderFeedbacks {
		_, err := govalidator.ValidateStruct(in)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgSellerFeedbacks) GetSignBytes() []byte {
	var submitTraderFeedbacks []json.RawMessage
	for _, submitTraderFeedback := range msg.SubmitTraderFeedbacks {
		submitTraderFeedbacks = append(submitTraderFeedbacks, submitTraderFeedback.GetSignBytes())
	}
	
	b, err := msgCdc.MarshalJSON(struct {
		SubmitTraderFeedbacks []json.RawMessage `json:"submitTraderFeedbacks"`
	}{
		SubmitTraderFeedbacks: submitTraderFeedbacks,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgSellerFeedbacks) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SubmitTraderFeedbacks))
	for i, r := range msg.SubmitTraderFeedbacks {
		addrs[i] = r.TraderFeedback.SellerAddress
	}
	return addrs
}

// BuildSellerFeedbackMsg : butild the FeedbackTx
func BuildSellerFeedbackMsg(buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, score int64) sdk.Msg {
	traderFeedback := sdk.NewTraderFeedback(buyerAddress, sellerAddress, pegHash, score)
	submitTraderFeedback := NewSubmitTraderFeedback(traderFeedback)
	msg := NewMsgSellerFeedbacks([]SubmitTraderFeedback{submitTraderFeedback})
	return &msg
}
