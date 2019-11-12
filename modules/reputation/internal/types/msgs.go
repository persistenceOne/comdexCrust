package types

import (
	"encoding/json"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/types"
)

type SubmitTraderFeedback struct {
	TraderFeedback types.TraderFeedback `json:"traderFeedback"`
}

func NewSubmitTraderFeedback(traderFeedback types.TraderFeedback) SubmitTraderFeedback {
	return SubmitTraderFeedback{TraderFeedback: traderFeedback}
}

func (submitTraderFeedback SubmitTraderFeedback) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		TraderFeedback types.TraderFeedback `json:"traderFeedback"`
	}{
		TraderFeedback: submitTraderFeedback.TraderFeedback,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

func (in SubmitTraderFeedback) ValidateBasic() cTypes.Error {
	if len(in.TraderFeedback.BuyerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.TraderFeedback.BuyerAddress.String())
	} else if len(in.TraderFeedback.SellerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.TraderFeedback.SellerAddress.String())
	} else if len(in.TraderFeedback.PegHash) == 0 {
		return cTypes.ErrUnknownRequest("peghash is empty")
	} else if in.TraderFeedback.Rating < 0 || in.TraderFeedback.Rating > 100 {
		return cTypes.ErrUnknownRequest("Rating should be 0-100")
	}
	return nil
}

type MsgBuyerFeedbacks struct {
	SubmitTraderFeedbacks []SubmitTraderFeedback `json:"submitTraderFeedbacks"`
}

var _ cTypes.Msg = MsgBuyerFeedbacks{}

func NewMsgBuyerFeedbacks(submitTraderFeedbacks []SubmitTraderFeedback) MsgBuyerFeedbacks {
	return MsgBuyerFeedbacks{SubmitTraderFeedbacks: submitTraderFeedbacks}
}

func (msg MsgBuyerFeedbacks) Type() string  { return "reputation" }
func (msg MsgBuyerFeedbacks) Route() string { return RouterKey }

func (msg MsgBuyerFeedbacks) ValidateBasic() cTypes.Error {
	if len(msg.SubmitTraderFeedbacks) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.SubmitTraderFeedbacks {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

func (msg MsgBuyerFeedbacks) GetSignBytes() []byte {
	var submitTraderFeedbacks []json.RawMessage
	for _, submitTraderFeedback := range msg.SubmitTraderFeedbacks {
		submitTraderFeedbacks = append(submitTraderFeedbacks, submitTraderFeedback.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		SubmitTraderFeedbacks []json.RawMessage `json:"submitTraderFeedbacks"`
	}{
		SubmitTraderFeedbacks: submitTraderFeedbacks,
	})
	if err != nil {
		panic(err)
	}
	return b
}

func (msg MsgBuyerFeedbacks) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.SubmitTraderFeedbacks))
	for i, r := range msg.SubmitTraderFeedbacks {
		addrs[i] = r.TraderFeedback.BuyerAddress
	}
	return addrs
}

func BuildBuyerFeedbackMsg(buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress,
	pegHash types.PegHash, score int64) cTypes.Msg {

	traderFeedback := types.NewTraderFeedback(buyerAddress, sellerAddress, pegHash, score)
	submitTraderFeedback := NewSubmitTraderFeedback(traderFeedback)
	msg := NewMsgBuyerFeedbacks([]SubmitTraderFeedback{submitTraderFeedback})
	return &msg
}

type MsgSellerFeedbacks struct {
	SubmitTraderFeedbacks []SubmitTraderFeedback `json:" submitTraderFeedbacks"`
}

var _ cTypes.Msg = MsgSellerFeedbacks{}

func NewMsgSellerFeedbacks(submitTraderFeedbacks []SubmitTraderFeedback) MsgSellerFeedbacks {
	return MsgSellerFeedbacks{SubmitTraderFeedbacks: submitTraderFeedbacks}
}

func (msg MsgSellerFeedbacks) Type() string  { return "reputation" }
func (msg MsgSellerFeedbacks) Route() string { return RouterKey }

func (msg MsgSellerFeedbacks) ValidateBasic() cTypes.Error {
	if len(msg.SubmitTraderFeedbacks) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.SubmitTraderFeedbacks {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

func (msg MsgSellerFeedbacks) GetSignBytes() []byte {
	var submitTraderFeedbacks []json.RawMessage
	for _, submitTraderFeedback := range msg.SubmitTraderFeedbacks {
		submitTraderFeedbacks = append(submitTraderFeedbacks, submitTraderFeedback.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		SubmitTraderFeedbacks []json.RawMessage `json:"submitTraderFeedbacks"`
	}{
		SubmitTraderFeedbacks: submitTraderFeedbacks,
	})
	if err != nil {
		panic(err)
	}
	return b
}

func (msg MsgSellerFeedbacks) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.SubmitTraderFeedbacks))
	for i, r := range msg.SubmitTraderFeedbacks {
		addrs[i] = r.TraderFeedback.SellerAddress
	}
	return addrs
}

func BuildSellerFeedbackMsg(buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress,
	pegHash types.PegHash, score int64) cTypes.Msg {
	traderFeedback := types.NewTraderFeedback(buyerAddress, sellerAddress, pegHash, score)
	submitTraderFeedback := NewSubmitTraderFeedback(traderFeedback)
	msg := NewMsgSellerFeedbacks([]SubmitTraderFeedback{submitTraderFeedback})
	return &msg
}
