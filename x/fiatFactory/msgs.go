package fiatFactory

import (
	"encoding/json"
	
	"github.com/asaskevich/govalidator"
	sdk "github.com/comdex-blockchain/types"
)

// *****Comdex

// *****IssueFiat

// IssueFiat - transaction input
type IssueFiat struct {
	IssuerAddress sdk.AccAddress `json:"issuerAddress"`
	ToAddress     sdk.AccAddress `json:"toAddress"`
	FiatPeg       sdk.FiatPeg    `json:"fiatPeg"`
}

// NewIssueFiat : initializer
func NewIssueFiat(issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, fiatPeg sdk.FiatPeg) IssueFiat {
	return IssueFiat{issuerAddress, toAddress, fiatPeg}
}

// GetSignBytes : get bytes to sign
func (in IssueFiat) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		IssuerAddress string      `json:"issuerAddress"`
		ToAddress     string      `json:"toAddress"`
		FiatPeg       sdk.FiatPeg `json:"fiatPeg"`
	}{
		IssuerAddress: in.IssuerAddress.String(),
		ToAddress:     in.ToAddress.String(),
		FiatPeg:       in.FiatPeg,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// #####IssueFiat

// *****MsgFactoryIssueFiats

// MsgFactoryIssueFiats : high level issuance of fiats module
type MsgFactoryIssueFiats struct {
	IssueFiats []IssueFiat `json:"issueFiats"`
}

// NewMsgFactoryIssueFiats : initilizer
func NewMsgFactoryIssueFiats(issueFiats []IssueFiat) MsgFactoryIssueFiats {
	return MsgFactoryIssueFiats{issueFiats}
}

// ***** Implementing sdk.Msg

var _ sdk.Msg = MsgFactoryIssueFiats{}

// Type : implements msg
func (msg MsgFactoryIssueFiats) Type() string { return "fiatFactory" }

// ValidateBasic : implements msg
func (msg MsgFactoryIssueFiats) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.IssueFiats {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		}
		_, err = govalidator.ValidateStruct(in.FiatPeg)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgFactoryIssueFiats) GetSignBytes() []byte {
	var issueFiats []json.RawMessage
	for _, issueFiat := range msg.IssueFiats {
		issueFiats = append(issueFiats, issueFiat.GetSignBytes())
	}
	
	b, err := msgCdc.MarshalJSON(struct {
		IssueFiats []json.RawMessage `json:"issueFiats"`
	}{
		IssueFiats: issueFiats,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgFactoryIssueFiats) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.IssueFiats))
	for i, in := range msg.IssueFiats {
		addrs[i] = in.IssuerAddress
	}
	return addrs
}

// BuildIssueFiatMsg : butild the issueFiatTx
func BuildIssueFiatMsg(from sdk.AccAddress, to sdk.AccAddress, fiatPeg sdk.FiatPeg) sdk.Msg {
	issueFiat := NewIssueFiat(from, to, fiatPeg)
	msg := NewMsgFactoryIssueFiats([]IssueFiat{issueFiat})
	return msg
}

// ##### Implement sdk.Msg

// #####MsgFactoryIssueFiats

// ****RedeemFiat

// RedeemFiat : transaction input
type RedeemFiat struct {
	RelayerAddress  sdk.AccAddress    `json:"relayerAddress"`
	RedeemerAddress sdk.AccAddress    `json:"redeemerAddress"`
	Amount          int64             `json:"amount"`
	FiatPegWallet   sdk.FiatPegWallet `json:"fiatPegWallet"`
}

// NewRedeemFiat : initializer
func NewRedeemFiat(relayerAddress sdk.AccAddress, redeemerAddress sdk.AccAddress, amount int64, fiatPegWallet sdk.FiatPegWallet) RedeemFiat {
	return RedeemFiat{relayerAddress, redeemerAddress, amount, fiatPegWallet}
}

// GetSignBytes : get bytes to sign
func (in RedeemFiat) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		RelayerAddress  string            `json:"relayerAddress"`
		RedeemerAddress string            `json:"redeemerAddress"`
		Amount          int64             `json:"amount"`
		FiatPegWallet   sdk.FiatPegWallet `json:"fiatPegWallet"`
	}{
		RelayerAddress:  in.RelayerAddress.String(),
		RedeemerAddress: in.RedeemerAddress.String(),
		Amount:          in.Amount,
		FiatPegWallet:   in.FiatPegWallet,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// #####RedeemFiat

// *****MsgFactoryRedeemFiats

// MsgFactoryRedeemFiats : high level issuance of fiats module
type MsgFactoryRedeemFiats struct {
	RedeemFiats []RedeemFiat `json:"redeemFiats"`
}

// NewMsgFactoryRedeemFiats : initilizer
func NewMsgFactoryRedeemFiats(redeemFiats []RedeemFiat) MsgFactoryRedeemFiats {
	return MsgFactoryRedeemFiats{redeemFiats}
}

// ***** Implementing sdk.Msg

var _ sdk.Msg = MsgFactoryRedeemFiats{}

// Type : implements msg
func (msg MsgFactoryRedeemFiats) Type() string { return "fiatFactory" }

// ValidateBasic : implements msg
func (msg MsgFactoryRedeemFiats) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInvalidAddress(err.Error())
	}
	for _, in := range msg.RedeemFiats {
		if len(in.RelayerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RelayerAddress.String())
		} else if len(in.RedeemerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RedeemerAddress.String())
		} else if in.Amount <= 0 {
			return sdk.ErrUnknownRequest("Amount should be positive")
		} else if len(in.FiatPegWallet) == 0 {
			return sdk.ErrUnknownRequest("FiatPegWallet is Empty")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgFactoryRedeemFiats) GetSignBytes() []byte {
	var redeemFiats []json.RawMessage
	for _, redeemFiat := range msg.RedeemFiats {
		redeemFiats = append(redeemFiats, redeemFiat.GetSignBytes())
	}
	
	b, err := msgCdc.MarshalJSON(struct {
		RedeemFiats []json.RawMessage `json:"redeemFiats"`
	}{
		RedeemFiats: redeemFiats,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgFactoryRedeemFiats) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.RedeemFiats))
	for i, in := range msg.RedeemFiats {
		addrs[i] = in.RelayerAddress
	}
	return addrs
}

// BuildRedeemFiatMsg : butild the issueFiatTx
func BuildRedeemFiatMsg(relayerAddress sdk.AccAddress, redeemerAddress sdk.AccAddress, amount int64, fiatPegWallet sdk.FiatPegWallet) sdk.Msg {
	redeemFiat := NewRedeemFiat(relayerAddress, redeemerAddress, amount, fiatPegWallet)
	msg := NewMsgFactoryRedeemFiats([]RedeemFiat{redeemFiat})
	return msg
}

// ##### Implement sdk.Msg

// #####MsgFactoryRedeemFiats

// *****SendFiat

// SendFiat - transaction input
type SendFiat struct {
	RelayerAddress sdk.AccAddress    `json:"relayerAddress"`
	FromAddress    sdk.AccAddress    `json:"fromAddress"`
	ToAddress      sdk.AccAddress    `json:"toAddress"`
	PegHash        sdk.PegHash       `json:"pegHash"`
	FiatPegWallet  sdk.FiatPegWallet `json:"fiatPegWallet"`
}

// NewSendFiat : initializer
func NewSendFiat(relayerAddress sdk.AccAddress, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, fiatPegWallet sdk.FiatPegWallet) SendFiat {
	return SendFiat{relayerAddress, fromAddress, toAddress, pegHash, fiatPegWallet}
}

// GetSignBytes : get bytes to sign
func (in SendFiat) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		RelayerAddress string            `json:"relayerAddress"`
		FromAddress    string            `json:"fromAddress"`
		ToAddress      string            `json:"toAddress"`
		PegHash        sdk.PegHash       `json:"pegHash"`
		FiatPegWallet  sdk.FiatPegWallet `json:"fiatPegWallet"`
	}{
		RelayerAddress: in.RelayerAddress.String(),
		FromAddress:    in.FromAddress.String(),
		ToAddress:      in.ToAddress.String(),
		PegHash:        in.PegHash,
		FiatPegWallet:  in.FiatPegWallet,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// #####SendFiat

// *****MsgFactorySendFiats

// MsgFactorySendFiats : high level issuance of fiats module
type MsgFactorySendFiats struct {
	SendFiats []SendFiat `json:"sendFiats"`
}

// NewMsgFactorySendFiats : initilizer
func NewMsgFactorySendFiats(sendFiats []SendFiat) MsgFactorySendFiats {
	return MsgFactorySendFiats{sendFiats}
}

// ***** Implementing sdk.Msg

var _ sdk.Msg = MsgFactorySendFiats{}

// Type : implements msg
func (msg MsgFactorySendFiats) Type() string { return "fiatFactory" }

// ValidateBasic : implements msg
func (msg MsgFactorySendFiats) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.SendFiats {
		if len(in.RelayerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RelayerAddress.String())
		} else if len(in.FromAddress) == 0 {
			return sdk.ErrInvalidAddress(in.FromAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if len(in.FiatPegWallet) == 0 {
			return sdk.ErrUnknownRequest("FiatPegWallet is Empty")
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgFactorySendFiats) GetSignBytes() []byte {
	var sendFiats []json.RawMessage
	for _, sendFiat := range msg.SendFiats {
		sendFiats = append(sendFiats, sendFiat.GetSignBytes())
	}
	
	b, err := msgCdc.MarshalJSON(struct {
		SendFiats []json.RawMessage `json:"sendFiats"`
	}{
		SendFiats: sendFiats,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgFactorySendFiats) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SendFiats))
	for i, in := range msg.SendFiats {
		addrs[i] = in.RelayerAddress
	}
	return addrs
}

// BuildSendFiatMsg : build the sendFiatTx
func BuildSendFiatMsg(relayerAddress sdk.AccAddress, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, fiatPegWallet sdk.FiatPegWallet) sdk.Msg {
	sendFiat := NewSendFiat(relayerAddress, fromAddress, toAddress, pegHash, fiatPegWallet)
	msg := NewMsgFactorySendFiats([]SendFiat{sendFiat})
	return msg
}

// ##### Implement sdk.Msg

// #####MsgFactorySendFiats

// *****MsgFactoryExecuteFiats

// MsgFactoryExecuteFiats : high level issuance of fiats module
type MsgFactoryExecuteFiats struct {
	SendFiats []SendFiat `json:"sendFiats"`
}

// NewMsgFactoryExecuteFiats : initilizer
func NewMsgFactoryExecuteFiats(sendFiats []SendFiat) MsgFactoryExecuteFiats {
	return MsgFactoryExecuteFiats{sendFiats}
}

// ***** Implementing sdk.Msg

var _ sdk.Msg = MsgFactoryExecuteFiats{}

// Type : implements msg
func (msg MsgFactoryExecuteFiats) Type() string { return "fiatFactory" }

// ValidateBasic : implements msg
func (msg MsgFactoryExecuteFiats) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.SendFiats {
		if len(in.RelayerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RelayerAddress.String())
		} else if len(in.FromAddress) == 0 {
			return sdk.ErrInvalidAddress(in.FromAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if len(in.FiatPegWallet) == 0 {
			return sdk.ErrUnknownRequest("FiatPegWallet is Empty")
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgFactoryExecuteFiats) GetSignBytes() []byte {
	var sendFiats []json.RawMessage
	for _, sendFiat := range msg.SendFiats {
		sendFiats = append(sendFiats, sendFiat.GetSignBytes())
	}
	
	b, err := msgCdc.MarshalJSON(struct {
		SendFiats []json.RawMessage `json:"sendFiats"`
	}{
		SendFiats: sendFiats,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgFactoryExecuteFiats) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SendFiats))
	for i, in := range msg.SendFiats {
		addrs[i] = in.RelayerAddress
	}
	return addrs
}

// BuildExecuteFiatMsg : build the executeFiatTx
func BuildExecuteFiatMsg(relayerAddress sdk.AccAddress, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, fiatPegWallet sdk.FiatPegWallet) sdk.Msg {
	sendFiat := NewSendFiat(relayerAddress, fromAddress, toAddress, pegHash, fiatPegWallet)
	msg := NewMsgFactoryExecuteFiats([]SendFiat{sendFiat})
	return msg
}

// ##### Implement sdk.Msg

// #####MsgFactoryExecuteFiats
// #####Comdex
