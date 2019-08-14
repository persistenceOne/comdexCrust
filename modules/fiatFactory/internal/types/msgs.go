package types

import (
	"encoding/json"
	"fmt"

	ctypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/types"
)

type IssueFiat struct {
	IssuerAddress ctypes.AccAddress `json:"issuerAddress"`
	ToAddress     ctypes.AccAddress `json:"toAddress"`
	FiatPeg       types.FiatPeg     `json:"fiatPeg"`
}

func NewIssueFiat(issuerAddress ctypes.AccAddress, toAddress ctypes.AccAddress, fiatPeg types.FiatPeg) IssueFiat {
	return IssueFiat{
		IssuerAddress: issuerAddress,
		ToAddress:     toAddress,
		FiatPeg:       fiatPeg,
	}
}

func (in IssueFiat) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		IssuerAddress string        `json:"issuerAddress"`
		ToAddress     string        `json:"toAddress"`
		FiatPeg       types.FiatPeg `json:"fiatPeg"`
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

func (in IssueFiat) ValidateBasic() ctypes.Error {
	if len(in.IssuerAddress) == 0 {
		return ctypes.ErrInvalidAddress(fmt.Sprintf("invalid address %s", in.IssuerAddress.String()))
	} else if len(in.ToAddress) == 0 {
		return ctypes.ErrInvalidAddress(fmt.Sprintf("invalid address %s", in.ToAddress.String()))
	} else if len(in.FiatPeg.GetPegHash()) == 0 {
		return ErrInvalidString(DefaultCodeSpace, "PegHash should not be empty")
	} else if in.FiatPeg.GetRedeemedAmount() < 0 {
		return ErrInvalidAmount(DefaultCodeSpace, "RedeemedAmount should not be negative")
	} else if in.FiatPeg.GetTransactionAmount() < 0 {
		return ErrInvalidAmount(DefaultCodeSpace, "TransactionAmount should not be negative")
	} else if in.FiatPeg.GetTransactionID() == "" {
		return ErrInvalidString(DefaultCodeSpace, "TransactionID should not be empty")
	}
	return nil
}

type MsgFactoryIssueFiats struct {
	IssueFiats []IssueFiat `json:"issueFiats"`
}

func NewMsgFactoryIssueFiats(issueFiats []IssueFiat) MsgFactoryIssueFiats {
	return MsgFactoryIssueFiats{IssueFiats: issueFiats}
}

var _ ctypes.Msg = MsgFactoryIssueFiats{}

func (msg MsgFactoryIssueFiats) Type() string { return "fiatfactry" }

func (msg MsgFactoryIssueFiats) Route() string { return RouterKey }

func (msg MsgFactoryIssueFiats) ValidateBasic() ctypes.Error {
	if len(msg.IssueFiats) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.IssueFiats {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

func (msg MsgFactoryIssueFiats) GetSignBytes() []byte {
	var issueFiats []json.RawMessage
	for _, issueFiat := range msg.IssueFiats {
		issueFiats = append(issueFiats, issueFiat.GetSignBytes())
	}
	bz, err := ModuleCdc.MarshalJSON(struct {
		IssueFiats []json.RawMessage `json:"issueFiats"`
	}{
		IssueFiats: issueFiats,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MsgFactoryIssueFiats) GetSigners() []ctypes.AccAddress {
	addr := make([]ctypes.AccAddress, len(msg.IssueFiats))
	for i, in := range msg.IssueFiats {
		addr[i] = in.IssuerAddress
	}
	return addr
}

func BuildIssueFiatMsg(issuerAddress ctypes.AccAddress, toAddress ctypes.AccAddress, fiatPeg types.FiatPeg) ctypes.Msg {
	issueFiat := NewIssueFiat(issuerAddress, toAddress, fiatPeg)
	msg := NewMsgFactoryIssueFiats([]IssueFiat{issueFiat})
	return msg
}

type RedeemFiat struct {
	RelayerAddress  ctypes.AccAddress   `json:"relayerAddress"`
	RedeemerAddress ctypes.AccAddress   `json:"redeemerAaddress"`
	Amount          int64               `json:"amount"`
	FiatPegWallet   types.FiatPegWallet `json:"fiatPegWallet"`
}

func NewRedeemFiat(relayerAddress ctypes.AccAddress, redeemerAddress ctypes.AccAddress, amount int64, wallet types.FiatPegWallet) RedeemFiat {
	return RedeemFiat{
		RelayerAddress:  relayerAddress,
		RedeemerAddress: redeemerAddress,
		Amount:          amount,
		FiatPegWallet:   wallet,
	}
}

func (in RedeemFiat) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		RelayerAddress  string              `json:"relayerAddress"`
		RedeemerAddress string              `json:"redeemerAddress"`
		Amount          int64               `json:"amount"`
		FiatPegWallet   types.FiatPegWallet `json:"fiatPegWallet"`
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

func (in RedeemFiat) ValidateBasic() ctypes.Error {
	if len(in.RelayerAddress) == 0 {
		return ctypes.ErrInvalidAddress(fmt.Sprintf("invalid address %s", in.RelayerAddress.String()))
	} else if len(in.RedeemerAddress) == 0 {
		return ctypes.ErrInvalidAddress(fmt.Sprintf("invalid address %s", in.RedeemerAddress.String()))
	} else if in.Amount < 0 {
		return ErrInvalidAmount(DefaultCodeSpace, "Amount should not be negative")
	}
	return nil
}

type MsgFactoryRedeemFiats struct {
	RedeemFiats []RedeemFiat `json:"redeemFiats"`
}

func NewMsgFactoryRedeemFiats(redeemFiats []RedeemFiat) MsgFactoryRedeemFiats {
	return MsgFactoryRedeemFiats{RedeemFiats: redeemFiats}
}

var _ ctypes.Msg = MsgFactoryRedeemFiats{}

func (msg MsgFactoryRedeemFiats) Type() string { return "fiatFactory" }

func (msg MsgFactoryRedeemFiats) Route() string { return RouterKey }

func (msg MsgFactoryRedeemFiats) ValidateBasic() ctypes.Error {
	if len(msg.RedeemFiats) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.RedeemFiats {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

func (msg MsgFactoryRedeemFiats) GetSignBytes() []byte {
	var redeemFiats []json.RawMessage
	for _, redeemFiat := range msg.RedeemFiats {
		redeemFiats = append(redeemFiats, redeemFiat.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		RedeemFiats []json.RawMessage `json:"redeemFiats"`
	}{
		RedeemFiats: redeemFiats,
	})
	if err != nil {
		panic(err)
	}
	return b
}

func (msg MsgFactoryRedeemFiats) GetSigners() []ctypes.AccAddress {
	addr := make([]ctypes.AccAddress, len(msg.RedeemFiats))
	for i, in := range msg.RedeemFiats {
		addr[i] = in.RelayerAddress
	}
	return addr
}

func BuildRedeemFiatMsg(relayerAddress ctypes.AccAddress, redeemerAddress ctypes.AccAddress, amount int64, wallet types.FiatPegWallet) ctypes.Msg {
	redeemfiat := NewRedeemFiat(relayerAddress, redeemerAddress, amount, wallet)
	msg := NewMsgFactoryRedeemFiats([]RedeemFiat{redeemfiat})
	return msg
}

type SendFiat struct {
	RelayerAddress ctypes.AccAddress   `json:"relayerAddress"`
	FromAddress    ctypes.AccAddress   `json:"fromAddress"`
	ToAddress      ctypes.AccAddress   `json:"toAddress"`
	PegHash        types.PegHash       `json:"pegHash"`
	FiatPegWallet  types.FiatPegWallet `json:"fiatPegWallet"`
}

func NewSendFiat(relayerAddress ctypes.AccAddress, fromAddress ctypes.AccAddress, toAddress ctypes.AccAddress, hash types.PegHash, wallet types.FiatPegWallet) SendFiat {
	return SendFiat{
		RelayerAddress: relayerAddress,
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		PegHash:        hash,
		FiatPegWallet:  wallet,
	}
}

func (in SendFiat) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		RelayerAddress string              `json:"relayerAddress"`
		FromAddress    string              `json:"fromAddress"`
		ToAddress      string              `json:"toAddress"`
		PegHash        types.PegHash       `json:"pegHash"`
		FiatPegWallet  types.FiatPegWallet `json:"fiatPegWallet"`
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

func (in SendFiat) ValidateBasic() ctypes.Error {
	if len(in.RelayerAddress) == 0 {
		return ctypes.ErrInvalidAddress(fmt.Sprintf("Invalid RelayerAddress %s", in.RelayerAddress.String()))
	} else if len(in.FromAddress) == 0 {
		return ctypes.ErrInvalidAddress(fmt.Sprintf("Invalid FromAddress", in.FromAddress.String()))
	} else if len(in.ToAddress) == 0 {
		return ctypes.ErrInvalidAddress(fmt.Sprintf("Invalid ToAddress", in.ToAddress.String()))
	} else if len(in.PegHash) == 0 {
		return ErrInvalidString(DefaultCodeSpace, "PegHash should not be empty")
	}
	return nil
}

type MsgFactorySendFiats struct {
	SendFiats []SendFiat `json:"sendFiats"`
}

func NewMsgFactorySendFiats(sendFiats []SendFiat) MsgFactorySendFiats {
	return MsgFactorySendFiats{SendFiats: sendFiats}
}

var _ ctypes.Msg = MsgFactorySendFiats{}

func (msg MsgFactorySendFiats) Type() string { return "fiatFactory" }

func (msg MsgFactorySendFiats) Route() string { return RouterKey }

func (msg MsgFactorySendFiats) ValidateBasic() ctypes.Error {
	if len(msg.SendFiats) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.SendFiats {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

func (msg MsgFactorySendFiats) GetSignBytes() []byte {
	var sendFiats []json.RawMessage
	for _, sendFiat := range msg.SendFiats {
		sendFiats = append(sendFiats, sendFiat.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		SendFiats []json.RawMessage `json:"sendFiats"`
	}{
		SendFiats: sendFiats,
	})
	if err != nil {
		panic(err)
	}
	return b
}

func (msg MsgFactorySendFiats) GetSigners() []ctypes.AccAddress {
	addrs := make([]ctypes.AccAddress, len(msg.SendFiats))
	for i, in := range msg.SendFiats {
		addrs[i] = in.RelayerAddress
	}
	return addrs
}

func BuildSendFiatMsg(relayerAddress ctypes.AccAddress, fromAddress ctypes.AccAddress, toAddress ctypes.AccAddress, pegHash types.PegHash, fiatPegWallet types.FiatPegWallet) ctypes.Msg {
	sendFiat := NewSendFiat(relayerAddress, fromAddress, toAddress, pegHash, fiatPegWallet)
	msg := NewMsgFactorySendFiats([]SendFiat{sendFiat})
	return msg
}

type MsgFactoryExecuteFiats struct {
	SendFiats []SendFiat `json:"sendFiats"`
}

// NewMsgFactoryExecuteFiats : initilizer
func NewMsgFactoryExecuteFiats(sendFiats []SendFiat) MsgFactoryExecuteFiats {
	return MsgFactoryExecuteFiats{sendFiats}
}

var _ ctypes.Msg = MsgFactoryExecuteFiats{}

// Type : implements msg
func (msg MsgFactoryExecuteFiats) Type() string { return "fiatFactory" }

func (msg MsgFactoryExecuteFiats) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgFactoryExecuteFiats) ValidateBasic() ctypes.Error {
	if len(msg.SendFiats) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.SendFiats {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
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

	b, err := ModuleCdc.MarshalJSON(struct {
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
func (msg MsgFactoryExecuteFiats) GetSigners() []ctypes.AccAddress {
	addrs := make([]ctypes.AccAddress, len(msg.SendFiats))
	for i, in := range msg.SendFiats {
		addrs[i] = in.RelayerAddress
	}
	return addrs
}

// BuildExecuteFiatMsg : build the executeFiatTx
func BuildExecuteFiatMsg(relayerAddress ctypes.AccAddress, fromAddress ctypes.AccAddress, toAddress ctypes.AccAddress, pegHash types.PegHash, fiatPegWallet types.FiatPegWallet) ctypes.Msg {
	sendFiat := NewSendFiat(relayerAddress, fromAddress, toAddress, pegHash, fiatPegWallet)
	msg := NewMsgFactoryExecuteFiats([]SendFiat{sendFiat})
	return msg
}
