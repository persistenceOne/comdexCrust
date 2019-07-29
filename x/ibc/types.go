package ibc

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"

	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/x/assetFactory"
	"github.com/commitHub/commitBlockchain/x/bank"
	"github.com/commitHub/commitBlockchain/x/fiatFactory"
)

// IBCPacket :
// nolint - TODO rename to Packet as IBCPacket stutters (golint)
// IBCPacket defines a piece of data that can be send between two separate
// blockchains.
type IBCPacket struct {
	SrcAddr   sdk.AccAddress
	DestAddr  sdk.AccAddress
	Coins     sdk.Coins
	SrcChain  string
	DestChain string
}

// NewIBCPacket : returns new ibs packet
func NewIBCPacket(srcAddr sdk.AccAddress, destAddr sdk.AccAddress, coins sdk.Coins,
	srcChain string, destChain string) IBCPacket {

	return IBCPacket{
		SrcAddr:   srcAddr,
		DestAddr:  destAddr,
		Coins:     coins,
		SrcChain:  srcChain,
		DestChain: destChain,
	}
}

// GetSignBytes :
//nolint
func (p IBCPacket) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(p)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic : validate the ibc packet
func (p IBCPacket) ValidateBasic() sdk.Error {
	if p.SrcChain == p.DestChain {
		return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
	}
	if !p.Coins.IsValid() {
		return sdk.ErrInvalidCoins("")
	}
	return nil
}

// IBCTransferMsg :
// nolint - TODO rename to TransferMsg as folks will reference with ibc.TransferMsg
// IBCTransferMsg defines how another module can send an IBCPacket.
type IBCTransferMsg struct {
	IBCPacket
}

// Type :
// nolint
func (msg IBCTransferMsg) Type() string { return "ibc" }

// GetSigners : x/bank/tx.go MsgSend.GetSigners()
func (msg IBCTransferMsg) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.SrcAddr} }

//GetSignBytes : get the sign bytes for ibc transfer message
func (msg IBCTransferMsg) GetSignBytes() []byte {
	return msg.IBCPacket.GetSignBytes()
}

//ValidateBasic : validate ibc transfer message
func (msg IBCTransferMsg) ValidateBasic() sdk.Error {
	return msg.IBCPacket.ValidateBasic()
}

// IBCReceiveMsg :
// nolint - TODO rename to ReceiveMsg as folks will reference with ibc.ReceiveMsg
// IBCReceiveMsg defines the message that a relayer uses to post an IBCPacket
// to the destination chain.
type IBCReceiveMsg struct {
	IBCPacket
	Relayer  sdk.AccAddress
	Sequence int64
}

// Type :
// nolint
func (msg IBCReceiveMsg) Type() string { return "ibc" }

// ValidateBasic : alidate the ibc packet
func (msg IBCReceiveMsg) ValidateBasic() sdk.Error { return msg.IBCPacket.ValidateBasic() }

//GetSigners : x/bank/tx.go MsgSend.GetSigners()
func (msg IBCReceiveMsg) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Relayer} }

// GetSignBytes : get the sign bytes for ibc receive message
func (msg IBCReceiveMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(struct {
		IBCPacket json.RawMessage
		Relayer   sdk.AccAddress
		Sequence  int64
	}{
		IBCPacket: json.RawMessage(msg.IBCPacket.GetSignBytes()),
		Relayer:   msg.Relayer,
		Sequence:  msg.Sequence,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

//*****Comdex
//*****IssueAsset

//IssueAsset - transaction input
type IssueAsset struct {
	IssuerAddress    sdk.AccAddress `json:"issuerAddress"`
	ToAddress        sdk.AccAddress `json:"toAddress"`
	AssetPeg         sdk.AssetPeg   `json:"assetPeg"`
	SourceChain      string         `json:"sourceChain"`
	DestinationChain string         `json:"destinationChain"`
}

//NewIssueAsset : initializer
func NewIssueAsset(issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, assetPeg sdk.AssetPeg, sourceChain string, destinationChain string) IssueAsset {
	return IssueAsset{issuerAddress, toAddress, assetPeg, sourceChain, destinationChain}
}

//GetSignBytes : get bytes to sign
func (in IssueAsset) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		IssuerAddress    string       `json:"issuerAddress"`
		ToAddress        string       `json:"toAddress"`
		AssetPeg         sdk.AssetPeg `json:"assetPeg"`
		SourceChain      string       `json:"sourceChain"`
		DestinationChain string       `json:"destinationChain"`
	}{
		IssuerAddress:    in.IssuerAddress.String(),
		ToAddress:        in.ToAddress.String(),
		AssetPeg:         in.AssetPeg,
		SourceChain:      in.SourceChain,
		DestinationChain: in.DestinationChain,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//#####IssueAsset

//*****MsgIssueAssets

//MsgIssueAssets : high level issuance of assets module
type MsgIssueAssets struct {
	IssueAssets []IssueAsset `json:"issueAssets"`
}

//NewMsgIssueAssets : initializer
func NewMsgIssueAssets(issueAssets []IssueAsset) MsgIssueAssets {
	return MsgIssueAssets{issueAssets}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgIssueAssets{}

//Type : implements msg
func (msg MsgIssueAssets) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgIssueAssets) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.IssueAssets {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
		_, err = govalidator.ValidateStruct(in.AssetPeg)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgIssueAssets) GetSignBytes() []byte {
	var issueAssets []json.RawMessage
	for _, issueAsset := range msg.IssueAssets {
		issueAssets = append(issueAssets, issueAsset.GetSignBytes())
	}

	b, err := msgCdc.MarshalJSON(struct {
		IssueAssets []json.RawMessage `json:"issueAssets"`
	}{
		IssueAssets: issueAssets,
	})
	if err != nil {
		panic(err)
	}
	return b
}

//GetSigners : implements msg
func (msg MsgIssueAssets) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.IssueAssets))
	for i, in := range msg.IssueAssets {
		addrs[i] = in.IssuerAddress
	}
	return addrs
}

//BuildIssueAssetMsg : butild the issueAssetTx
func BuildIssueAssetMsg(issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, assetPeg sdk.AssetPeg, sourceChain string, destinationChain string) sdk.Msg {
	issueAsset := NewIssueAsset(issuerAddress, toAddress, assetPeg, sourceChain, destinationChain)
	msg := NewMsgIssueAssets([]IssueAsset{issueAsset})
	return msg
}

//##### Implement sdk.Msg

//#####MsgIssueAssets

//*****MsgRelayIssueAssets

//MsgRelayIssueAssets : high level issuance of assets module
type MsgRelayIssueAssets struct {
	IssueAssets []IssueAsset   `json:"issueAssets"`
	Relayer     sdk.AccAddress `json:"relayer"`
	Sequence    int64          `json:"sequence"`
}

//NewMsgRelayIssueAssets : initializer
func NewMsgRelayIssueAssets(issueAssets []IssueAsset, relayer sdk.AccAddress, sequence int64) MsgRelayIssueAssets {
	return MsgRelayIssueAssets{issueAssets, relayer, sequence}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgRelayIssueAssets{}

//Type : implements msg
func (msg MsgRelayIssueAssets) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgRelayIssueAssets) ValidateBasic() sdk.Error {
	if len(msg.Relayer) == 0 {
		return sdk.ErrInvalidAddress(msg.Relayer.String())
	}
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.IssueAssets {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
		_, err = govalidator.ValidateStruct(in.AssetPeg)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgRelayIssueAssets) GetSignBytes() []byte {
	var issueAssets []json.RawMessage
	for _, issueAsset := range msg.IssueAssets {
		issueAssets = append(issueAssets, issueAsset.GetSignBytes())
	}

	b, err := msgCdc.MarshalJSON(struct {
		IssueAssets []json.RawMessage `json:"issueAssets"`
		Relayer     string            `json:"relayer"`
		Sequence    int64             `json:"sequence"`
	}{
		IssueAssets: issueAssets,
		Relayer:     msg.Relayer.String(),
		Sequence:    msg.Sequence,
	})
	if err != nil {
		panic(err)
	}
	return b
}

//GetSigners : implements msg
func (msg MsgRelayIssueAssets) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}

//BuildRelayIssueAssetMsg : build the issueAssetTx
func BuildRelayIssueAssetMsg(issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, assetPeg sdk.AssetPeg, sourceChain string, destinationChain string, relayer sdk.AccAddress, sequence int64) sdk.Msg {
	issueAsset := NewIssueAsset(issuerAddress, toAddress, assetPeg, sourceChain, destinationChain)
	msg := NewMsgRelayIssueAssets([]IssueAsset{issueAsset}, relayer, sequence)
	return msg
}

//##### Implement sdk.Msg

//#####MsgRelayIssueAssets

func toHubMsgIssueAssets(hubMsgIssueAssets MsgIssueAssets) bank.MsgBankIssueAssets {
	var newIssueAssets []bank.IssueAsset
	for _, issueAsset := range hubMsgIssueAssets.IssueAssets {
		newIssueAssets = append(newIssueAssets, bank.NewIssueAsset(issueAsset.IssuerAddress, issueAsset.ToAddress, issueAsset.AssetPeg))
	}
	return bank.NewMsgBankIssueAssets(newIssueAssets)
}

//toTypeMsgIssueAssets : from type msgIssueAsset to assetFactory msgIssueAsset
func toTypeMsgIssueAssets(bankMsgIssueAssets MsgIssueAssets) assetFactory.MsgFactoryIssueAssets {
	var newIssueAssets []assetFactory.IssueAsset
	for _, issueAssets := range bankMsgIssueAssets.IssueAssets {
		newIssueAssets = append(newIssueAssets, assetFactory.NewIssueAsset(issueAssets.IssuerAddress, issueAssets.ToAddress, issueAssets.AssetPeg))
	}
	return assetFactory.NewMsgFactoryIssueAssets(newIssueAssets)
}

//######

//*****RedeemAsset

//RedeemAsset : transaction input
type RedeemAsset struct {
	IssuerAddress    sdk.AccAddress `json:"issuerAddress"`
	RedeemerAddress  sdk.AccAddress `json:"redeemerAddress"`
	PegHash          sdk.PegHash    `json:"pegHash"`
	SourceChain      string         `json:"sourceChain"`
	DestinationChain string         `json:"destinationChain"`
}

//NewRedeemAsset : initializer
func NewRedeemAsset(issuerAddress sdk.AccAddress, redeemerAddress sdk.AccAddress, pegHash sdk.PegHash, sourceChain string, destinationChain string) RedeemAsset {
	return RedeemAsset{issuerAddress, redeemerAddress, pegHash, sourceChain, destinationChain}
}

//GetSignBytes : get bytes to sign
func (in RedeemAsset) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		IssuerAddress    string      `json:"issuerAddress"`
		RedeemerAddress  string      `json:"redeemeraddress"`
		PegHash          sdk.PegHash `json:"pegHash"`
		SourceChain      string      `json:"sourceChain"`
		DestinationChain string      `json:"destinationChain"`
	}{
		IssuerAddress:    in.IssuerAddress.String(),
		RedeemerAddress:  in.RedeemerAddress.String(),
		PegHash:          in.PegHash,
		SourceChain:      in.SourceChain,
		DestinationChain: in.DestinationChain,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//#####RedeemAsset

//*****MsgRedeemAssets

//MsgRedeemAssets : high level redeenance of assets module
type MsgRedeemAssets struct {
	RedeemAssets []RedeemAsset `json:"redeemAssets"`
}

//NewMsgRedeemAssets : initializer
func NewMsgRedeemAssets(redeemAssets []RedeemAsset) MsgRedeemAssets {
	return MsgRedeemAssets{redeemAssets}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgRedeemAssets{}

//Type : implements msg
func (msg MsgRedeemAssets) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgRedeemAssets) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.RedeemAssets {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.RedeemerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RedeemerAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgRedeemAssets) GetSignBytes() []byte {
	var redeemAssets []json.RawMessage
	for _, redeemAsset := range msg.RedeemAssets {
		redeemAssets = append(redeemAssets, redeemAsset.GetSignBytes())
	}

	bz, err := msgCdc.MarshalJSON(struct {
		RedeemAssets []json.RawMessage `json:"redeemAssets"`
	}{
		RedeemAssets: redeemAssets,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

//GetSigners : implements msg
func (msg MsgRedeemAssets) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.RedeemAssets))
	for i, in := range msg.RedeemAssets {
		addrs[i] = in.RedeemerAddress
	}
	return addrs
}

//BuildRedeemAssetMsg : butild the redeemAssetTx
func BuildRedeemAssetMsg(issuerAddress sdk.AccAddress, redeemerAddress sdk.AccAddress, pegHash sdk.PegHash, sourceChain string, destinationChain string) sdk.Msg {
	redeemAsset := NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, sourceChain, destinationChain)
	msg := NewMsgRedeemAssets([]RedeemAsset{redeemAsset})
	return msg
}

//#####MsgRedeemAssets

//*****MsgRelayRedeemAssets

//MsgRelayRedeemAssets :
type MsgRelayRedeemAssets struct {
	RedeemAssets []RedeemAsset  `json:"redeemAssets"`
	Relayer      sdk.AccAddress `json:"relayer"`
	Sequence     int64          `json:"sequence"`
}

//NewMsgRelayRedeemAssets : initializer
func NewMsgRelayRedeemAssets(redeemAssets []RedeemAsset, relayer sdk.AccAddress, sequence int64) MsgRelayRedeemAssets {
	return MsgRelayRedeemAssets{redeemAssets, relayer, sequence}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgRelayRedeemAssets{}

//Type : implements msg
func (msg MsgRelayRedeemAssets) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgRelayRedeemAssets) ValidateBasic() sdk.Error {
	if len(msg.Relayer) == 0 {
		return sdk.ErrInvalidAddress(msg.Relayer.String())
	}
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.RedeemAssets {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.RedeemerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RedeemerAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgRelayRedeemAssets) GetSignBytes() []byte {
	var redeemAssets []json.RawMessage
	for _, redeemAsset := range msg.RedeemAssets {
		redeemAssets = append(redeemAssets, redeemAsset.GetSignBytes())
	}

	bz, err := msgCdc.MarshalJSON(struct {
		RedeemAssets []json.RawMessage `json:"redeemAssets"`
		Relayer      string            `json:"relayer"`
		Sequence     int64             `json:"sequence"`
	}{
		RedeemAssets: redeemAssets,
		Relayer:      msg.Relayer.String(),
		Sequence:     msg.Sequence,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

//GetSigners : implements msg
func (msg MsgRelayRedeemAssets) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}

//BuildRelayRedeemAssetMsg : build the issueAssetTx
func BuildRelayRedeemAssetMsg(issuerAddress sdk.AccAddress, redeemerAddress sdk.AccAddress, pegHash sdk.PegHash, sourceChain string, destinationChain string, relayer sdk.AccAddress, sequence int64) sdk.Msg {
	redeemAsset := NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, sourceChain, destinationChain)
	msg := NewMsgRelayRedeemAssets([]RedeemAsset{redeemAsset}, relayer, sequence)
	return msg
}

//##### Implement sdk.Msg

//#####MsgRelayRedeemAssets

func toHubMsgRedeemAssets(hubMsgRedeemAssets MsgRedeemAssets) bank.MsgBankRedeemAssets {
	var newRedeemAssets []bank.RedeemAsset
	for _, redeemAsset := range hubMsgRedeemAssets.RedeemAssets {
		newRedeemAssets = append(newRedeemAssets, bank.NewRedeemAsset(redeemAsset.IssuerAddress, redeemAsset.RedeemerAddress, redeemAsset.PegHash))
	}
	return bank.NewMsgBankRedeemAssets(newRedeemAssets)
}

//toTypeMsgRedeemAssets : from type msgIssueAsset to assetFactory msgIssueAsset
func toTypeMsgRedeemAssets(bankMsgRedeemAssets MsgRedeemAssets, bankMsgRelayRedeemAssets MsgRelayRedeemAssets) assetFactory.MsgFactoryRedeemAssets {
	var newRedeemAssets []assetFactory.RedeemAsset
	for _, redeemAsset := range bankMsgRedeemAssets.RedeemAssets {
		newRedeemAssets = append(newRedeemAssets, assetFactory.NewRedeemAsset(bankMsgRelayRedeemAssets.Relayer, redeemAsset.RedeemerAddress, redeemAsset.IssuerAddress, redeemAsset.PegHash))
	}
	return assetFactory.NewMsgFactoryRedeemAssets(newRedeemAssets)
}

//######

//*****IssueFiat

//IssueFiat - transaction input
type IssueFiat struct {
	IssuerAddress    sdk.AccAddress `json:"issuerAddress"`
	ToAddress        sdk.AccAddress `json:"toAddress"`
	FiatPeg          sdk.FiatPeg    `json:"fiatPeg"`
	SourceChain      string         `json:"sourceChain"`
	DestinationChain string         `json:"destinationChain"`
}

//NewIssueFiat : initializer
func NewIssueFiat(issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, fiatPeg sdk.FiatPeg, sourceChain string, destinationChain string) IssueFiat {
	return IssueFiat{issuerAddress, toAddress, fiatPeg, sourceChain, destinationChain}
}

//GetSignBytes : get bytes to sign
func (in IssueFiat) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		IssuerAddress    string      `json:"issuerAddress"`
		ToAddress        string      `json:"toAddress"`
		FiatPeg          sdk.FiatPeg `json:"fiatPeg"`
		SourceChain      string      `json:"sourceChain"`
		DestinationChain string      `json:"destinationChain"`
	}{
		IssuerAddress:    in.IssuerAddress.String(),
		ToAddress:        in.ToAddress.String(),
		FiatPeg:          in.FiatPeg,
		SourceChain:      in.SourceChain,
		DestinationChain: in.DestinationChain,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//#####IssueFiat

//*****MsgIssueFiats

//MsgIssueFiats : high level issuance of fiats module
type MsgIssueFiats struct {
	IssueFiats []IssueFiat `json:"issueFiats"`
}

//NewMsgIssueFiats : initializer
func NewMsgIssueFiats(issueFiats []IssueFiat) MsgIssueFiats {
	return MsgIssueFiats{issueFiats}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgIssueFiats{}

//Type : implements msg
func (msg MsgIssueFiats) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgIssueFiats) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.IssueFiats {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
		_, err = govalidator.ValidateStruct(in.FiatPeg)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgIssueFiats) GetSignBytes() []byte {
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

//GetSigners : implements msg
func (msg MsgIssueFiats) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.IssueFiats))
	for i, in := range msg.IssueFiats {
		addrs[i] = in.IssuerAddress
	}
	return addrs
}

//BuildIssueFiatMsg : butild the issueFiatTx
func BuildIssueFiatMsg(issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, fiatPeg sdk.FiatPeg, sourceChain string, destinationChain string) sdk.Msg {
	issueFiat := NewIssueFiat(issuerAddress, toAddress, fiatPeg, sourceChain, destinationChain)
	msg := NewMsgIssueFiats([]IssueFiat{issueFiat})
	return msg
}

//##### Implement sdk.Msg

//#####MsgIssueFiats

//*****MsgRelayIssueFiats

//MsgRelayIssueFiats : high level issuance of fiats module
type MsgRelayIssueFiats struct {
	IssueFiats []IssueFiat    `json:"issueFiats"`
	Relayer    sdk.AccAddress `json:"relayer"`
	Sequence   int64          `json:"sequence"`
}

//NewMsgRelayIssueFiats : initializer
func NewMsgRelayIssueFiats(issueFiats []IssueFiat, relayer sdk.AccAddress, sequence int64) MsgRelayIssueFiats {
	return MsgRelayIssueFiats{issueFiats, relayer, sequence}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgRelayIssueFiats{}

//Type : implements msg
func (msg MsgRelayIssueFiats) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgRelayIssueFiats) ValidateBasic() sdk.Error {
	if len(msg.Relayer) == 0 {
		return sdk.ErrInvalidAddress(msg.Relayer.String())
	}
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.IssueFiats {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
		_, err = govalidator.ValidateStruct(in.FiatPeg)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgRelayIssueFiats) GetSignBytes() []byte {
	var issueFiats []json.RawMessage
	for _, issueFiat := range msg.IssueFiats {
		issueFiats = append(issueFiats, issueFiat.GetSignBytes())
	}

	b, err := msgCdc.MarshalJSON(struct {
		IssueFiats []json.RawMessage `json:"issueFiats"`
		Relayer    string            `json:"relayer"`
		Sequence   int64             `json:"sequence"`
	}{
		IssueFiats: issueFiats,
		Relayer:    msg.Relayer.String(),
		Sequence:   msg.Sequence,
	})
	if err != nil {
		panic(err)
	}
	return b
}

//GetSigners : implements msg
func (msg MsgRelayIssueFiats) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}

//BuildRelayIssueFiatMsg : build the issueFiatTx
func BuildRelayIssueFiatMsg(issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, fiatPeg sdk.FiatPeg, sourceChain string, destinationChain string, relayer sdk.AccAddress, sequence int64) sdk.Msg {
	issueFiat := NewIssueFiat(issuerAddress, toAddress, fiatPeg, sourceChain, destinationChain)
	msg := NewMsgRelayIssueFiats([]IssueFiat{issueFiat}, relayer, sequence)
	return msg
}

//##### Implement sdk.Msg

//#####MsgRelayIssueFiats

func toHubMsgIssueFiats(hubMsgIssueFiats MsgIssueFiats) bank.MsgBankIssueFiats {
	var newIssueFiats []bank.IssueFiat
	for _, issueFiat := range hubMsgIssueFiats.IssueFiats {
		newIssueFiats = append(newIssueFiats, bank.NewIssueFiat(issueFiat.IssuerAddress, issueFiat.ToAddress, issueFiat.FiatPeg))
	}
	return bank.NewMsgBankIssueFiats(newIssueFiats)
}

//toTypeMsgIssueFiats : from type msgIssueFiat to fiatFactory msgIssueFiat
func toTypeMsgIssueFiats(bankMsgIssueFiats MsgIssueFiats) fiatFactory.MsgFactoryIssueFiats {
	var newIssueFiats []fiatFactory.IssueFiat
	for _, issueFiats := range bankMsgIssueFiats.IssueFiats {
		newIssueFiats = append(newIssueFiats, fiatFactory.NewIssueFiat(issueFiats.IssuerAddress, issueFiats.ToAddress, issueFiats.FiatPeg))
	}
	return fiatFactory.NewMsgFactoryIssueFiats(newIssueFiats)
}

//*****RedeemFiat

//RedeemFiat : transaction input
type RedeemFiat struct {
	RedeemerAddress  sdk.AccAddress    `json:"redeemerAddress"`
	IssuerAddress    sdk.AccAddress    `json:"issuerAddress"`
	Amount           int64             `json:"amount"`
	FiatPegWallet    sdk.FiatPegWallet `json:"fiatPegWallet"`
	SourceChain      string            `json:"sourceChain"`
	DestinationChain string            `json:"destinationChain"`
}

//NewRedeemFiat : initializer
func NewRedeemFiat(redeemerAddress sdk.AccAddress, issuerAddress sdk.AccAddress, amount int64, fiatPegWallet sdk.FiatPegWallet, sourceChain string, destinationChain string) RedeemFiat {
	return RedeemFiat{redeemerAddress, issuerAddress, amount, fiatPegWallet, sourceChain, destinationChain}
}

//GetSignBytes : get bytes to sign
func (in RedeemFiat) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		RedeemerAddress  string            `json:"redeemerAddress"`
		IssuerAddress    string            `json:"issuerAddress"`
		Amount           int64             `json:"amount"`
		FiatPegWallet    sdk.FiatPegWallet `json:"fiatPegWallet"`
		SourceChain      string            `json:"sourceChain"`
		DestinationChain string            `json:"destinationChain"`
	}{
		RedeemerAddress:  in.RedeemerAddress.String(),
		IssuerAddress:    in.IssuerAddress.String(),
		Amount:           in.Amount,
		FiatPegWallet:    in.FiatPegWallet,
		SourceChain:      in.SourceChain,
		DestinationChain: in.DestinationChain,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//#####RedeemFiat

//*****MsgRedeemFiats

//MsgRedeemFiats : high level redeenance of fiats module
type MsgRedeemFiats struct {
	RedeemFiats []RedeemFiat `json:"redeemFiats"`
}

//NewMsgRedeemFiats : initializer
func NewMsgRedeemFiats(redeemFiats []RedeemFiat) MsgRedeemFiats {
	return MsgRedeemFiats{redeemFiats}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgRedeemFiats{}

//Type : implements msg
func (msg MsgRedeemFiats) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgRedeemFiats) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.RedeemFiats {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.RedeemerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RedeemerAddress.String())
		} else if in.Amount <= 0 {
			return sdk.ErrUnknownRequest("Amount should be Positive")
		} else if len(in.FiatPegWallet) == 0 {
			return sdk.ErrUnknownRequest("FiatPegWallet is Empty")
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgRedeemFiats) GetSignBytes() []byte {
	var redeemFiats []json.RawMessage
	for _, redeemFiat := range msg.RedeemFiats {
		redeemFiats = append(redeemFiats, redeemFiat.GetSignBytes())
	}

	bz, err := msgCdc.MarshalJSON(struct {
		RedeemFiats []json.RawMessage `json:"redeemFiats"`
	}{
		RedeemFiats: redeemFiats,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

//GetSigners : implements msg
func (msg MsgRedeemFiats) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.RedeemFiats))
	for i, in := range msg.RedeemFiats {
		addrs[i] = in.RedeemerAddress
	}
	return addrs
}

//BuildRedeemFiatMsg : butild the redeemFiatTx
func BuildRedeemFiatMsg(redeemerAddress sdk.AccAddress, issuerAddress sdk.AccAddress, amount int64, fiatPegWallet sdk.FiatPegWallet, sourceChain string, destinationChain string) sdk.Msg {
	redeemFiat := NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, sourceChain, destinationChain)
	msg := NewMsgRedeemFiats([]RedeemFiat{redeemFiat})
	return msg
}

//#####MsgRedeemFiats

//*****MsgRelayRedeemFiats

//MsgRelayRedeemFiats :
type MsgRelayRedeemFiats struct {
	RedeemFiats []RedeemFiat   `json:"redeemFiats"`
	Relayer     sdk.AccAddress `json:"relayer"`
	Sequence    int64          `json:"sequence"`
}

//NewMsgRelayRedeemFiats : initializer
func NewMsgRelayRedeemFiats(redeemFiats []RedeemFiat, relayer sdk.AccAddress, sequence int64) MsgRelayRedeemFiats {
	return MsgRelayRedeemFiats{redeemFiats, relayer, sequence}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgRelayRedeemFiats{}

//Type : implements msg
func (msg MsgRelayRedeemFiats) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgRelayRedeemFiats) ValidateBasic() sdk.Error {
	if len(msg.Relayer) == 0 {
		return sdk.ErrInvalidAddress(msg.Relayer.String())
	}
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.RedeemFiats {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.RedeemerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RedeemerAddress.String())
		} else if in.Amount <= 0 {
			return sdk.ErrUnknownRequest("Amount should be Positive")
		} else if len(in.FiatPegWallet) == 0 {
			return sdk.ErrUnknownRequest("FiatPegWallet is Empty")
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgRelayRedeemFiats) GetSignBytes() []byte {
	var redeemFiats []json.RawMessage
	for _, redeemFiat := range msg.RedeemFiats {
		redeemFiats = append(redeemFiats, redeemFiat.GetSignBytes())
	}

	bz, err := msgCdc.MarshalJSON(struct {
		RedeemFiats []json.RawMessage `json:"redeemFiats"`
		Relayer     string            `json:"relayer"`
		Sequence    int64             `json:"sequence"`
	}{
		RedeemFiats: redeemFiats,
		Relayer:     msg.Relayer.String(),
		Sequence:    msg.Sequence,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

//GetSigners : implements msg
func (msg MsgRelayRedeemFiats) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}

//BuildRelayRedeemFiatMsg : build the issueFiatTx
func BuildRelayRedeemFiatMsg(redeemerAddress sdk.AccAddress, issuerAddress sdk.AccAddress, amount int64, fiatPegWallet sdk.FiatPegWallet, sourceChain string, destinationChain string, relayer sdk.AccAddress, sequence int64) sdk.Msg {
	redeemFiat := NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, sourceChain, destinationChain)
	msg := NewMsgRelayRedeemFiats([]RedeemFiat{redeemFiat}, relayer, sequence)
	return msg
}

//##### Implement sdk.Msg

//#####MsgRelayRedeemFiats

func toHubMsgRedeemFiats(hubMsgRedeemFiats MsgRedeemFiats) bank.MsgBankRedeemFiats {
	var newRedeemFiats []bank.RedeemFiat
	for _, redeemFiat := range hubMsgRedeemFiats.RedeemFiats {
		newRedeemFiats = append(newRedeemFiats, bank.NewRedeemFiat(redeemFiat.RedeemerAddress, redeemFiat.IssuerAddress, redeemFiat.Amount))
	}
	return bank.NewMsgBankRedeemFiats(newRedeemFiats)
}

//toTypeMsgRedeemFiats : from type msgIssueFiat to fiatFactory msgIssueFiat
func toTypeMsgRedeemFiats(bankMsgRedeemFiats MsgRedeemFiats, bankMsgRelayRedeemFiats MsgRelayRedeemFiats) fiatFactory.MsgFactoryRedeemFiats {
	var newRedeemFiats []fiatFactory.RedeemFiat
	for _, redeemFiat := range bankMsgRelayRedeemFiats.RedeemFiats {
		newRedeemFiats = append(newRedeemFiats, fiatFactory.NewRedeemFiat(bankMsgRelayRedeemFiats.Relayer, redeemFiat.RedeemerAddress, redeemFiat.Amount, redeemFiat.FiatPegWallet))
	}
	return fiatFactory.NewMsgFactoryRedeemFiats(newRedeemFiats)
}

//*****SendAsset

//SendAsset - transaction input
type SendAsset struct {
	FromAddress      sdk.AccAddress `json:"fromAddress"`
	ToAddress        sdk.AccAddress `json:"toAddress"`
	PegHash          sdk.PegHash    `json:"pegHash"`
	SourceChain      string         `json:"sourceChain"`
	DestinationChain string         `json:"destinationChain"`
}

//NewSendAsset : initializer
func NewSendAsset(fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, sourceChain string, destinationChain string) SendAsset {
	return SendAsset{fromAddress, toAddress, pegHash, sourceChain, destinationChain}
}

//GetSignBytes : get bytes to sign
func (in SendAsset) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		FromAddress      string      `json:"fromAddress"`
		ToAddress        string      `json:"toAddress"`
		PegHash          sdk.PegHash `json:"pegHash"`
		SourceChain      string      `json:"sourceChain"`
		DestinationChain string      `json:"destinationChain"`
	}{
		FromAddress:      in.FromAddress.String(),
		ToAddress:        in.ToAddress.String(),
		PegHash:          in.PegHash,
		SourceChain:      in.SourceChain,
		DestinationChain: in.DestinationChain,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//#####SendAsset

//*****MsgSendAssets

//MsgSendAssets : high level issuance of assets module
type MsgSendAssets struct {
	SendAssets []SendAsset `json:"sendAssets"`
}

//NewMsgSendAssets : initializer
func NewMsgSendAssets(sendAssets []SendAsset) MsgSendAssets {
	return MsgSendAssets{sendAssets}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgSendAssets{}

//Type : implements msg
func (msg MsgSendAssets) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgSendAssets) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.SendAssets {
		if len(in.FromAddress) == 0 {
			return sdk.ErrInvalidAddress(in.FromAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgSendAssets) GetSignBytes() []byte {
	var sendAssets []json.RawMessage
	for _, sendAsset := range msg.SendAssets {
		sendAssets = append(sendAssets, sendAsset.GetSignBytes())
	}

	b, err := msgCdc.MarshalJSON(struct {
		SendAssets []json.RawMessage `json:"sendAssets"`
	}{
		SendAssets: sendAssets,
	})
	if err != nil {
		panic(err)
	}
	return b
}

//GetSigners : implements msg
func (msg MsgSendAssets) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SendAssets))
	for i, in := range msg.SendAssets {
		addrs[i] = in.FromAddress
	}
	return addrs
}

//BuildSendAssetMsg : build the sendAssetTx
func BuildSendAssetMsg(fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, sourceChain string, destinationChain string) sdk.Msg {
	sendAsset := NewSendAsset(fromAddress, toAddress, pegHash, sourceChain, destinationChain)
	msg := NewMsgSendAssets([]SendAsset{sendAsset})
	return msg
}

//##### Implement sdk.Msg

//#####MsgSendAssets

//*****MsgRelaySendAssets

//MsgRelaySendAssets :
type MsgRelaySendAssets struct {
	SendAssets []SendAsset    `json:"sendAssets"`
	Relayer    sdk.AccAddress `json:"relayer"`
	Sequence   int64          `json:"sequence"`
}

//NewMsgRelaySendAssets : initializer
func NewMsgRelaySendAssets(sendAssets []SendAsset, relayer sdk.AccAddress, sequence int64) MsgRelaySendAssets {
	return MsgRelaySendAssets{sendAssets, relayer, sequence}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgRelaySendAssets{}

//Type : implements msg
func (msg MsgRelaySendAssets) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgRelaySendAssets) ValidateBasic() sdk.Error {
	if len(msg.Relayer) == 0 {
		return sdk.ErrInvalidAddress(msg.Relayer.String())
	}
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.SendAssets {
		if len(in.FromAddress) == 0 {
			return sdk.ErrInvalidAddress(in.FromAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgRelaySendAssets) GetSignBytes() []byte {
	var sendAssets []json.RawMessage
	for _, sendAsset := range msg.SendAssets {
		sendAssets = append(sendAssets, sendAsset.GetSignBytes())
	}

	bz, err := msgCdc.MarshalJSON(struct {
		SendAssets []json.RawMessage `json:"sendAssets"`
		Relayer    string            `json:"relayer"`
		Sequence   int64             `json:"sequence"`
	}{
		SendAssets: sendAssets,
		Relayer:    msg.Relayer.String(),
		Sequence:   msg.Sequence,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

//GetSigners : implements msg
func (msg MsgRelaySendAssets) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}

//BuildRelaySendAssetMsg : build the issueAssetTx
func BuildRelaySendAssetMsg(fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, sourceChain string, destinationChain string, relayer sdk.AccAddress, sequence int64) sdk.Msg {
	sendAsset := NewSendAsset(fromAddress, toAddress, pegHash, sourceChain, destinationChain)
	msg := NewMsgRelaySendAssets([]SendAsset{sendAsset}, relayer, sequence)
	return msg
}

//##### Implement sdk.Msg

//#####MsgRelaySendAssets

func toHubMsgSendAssets(hubMsgSendAssets MsgSendAssets) bank.MsgBankSendAssets {
	var newSendAssets []bank.SendAsset
	for _, sendAsset := range hubMsgSendAssets.SendAssets {
		newSendAssets = append(newSendAssets, bank.NewSendAsset(sendAsset.FromAddress, sendAsset.ToAddress, sendAsset.PegHash))
	}
	return bank.NewMsgBankSendAssets(newSendAssets)
}

//toTypeMsgSendAssets : from type msgIssueAsset to assetFactory msgIssueAsset
func toTypeMsgSendAssets(bankMsgSendAssets MsgSendAssets, bankMsgRelaySendAsset MsgRelaySendAssets) assetFactory.MsgFactorySendAssets {
	var newSendAssets []assetFactory.SendAsset
	for _, sendAsset := range bankMsgRelaySendAsset.SendAssets {
		newSendAssets = append(newSendAssets, assetFactory.NewSendAsset(bankMsgRelaySendAsset.Relayer, sendAsset.FromAddress, sendAsset.ToAddress, sendAsset.PegHash))
	}
	return assetFactory.NewMsgFactorySendAssets(newSendAssets)
}

//*****SendFiat

//SendFiat - transaction input
type SendFiat struct {
	FromAddress      sdk.AccAddress    `json:"fromAddress"`
	ToAddress        sdk.AccAddress    `json:"toAddress"`
	PegHash          sdk.PegHash       `json:"pegHash"`
	Amount           int64             `json:"amount"`
	FiatPegWallet    sdk.FiatPegWallet `json:"fiatPegWallet"`
	SourceChain      string            `json:"sourceChain"`
	DestinationChain string            `json:"destinationChain"`
}

//NewSendFiat : initializer
func NewSendFiat(fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, amount int64, fiatPegWallet sdk.FiatPegWallet, sourceChain string, destinationChain string) SendFiat {
	return SendFiat{fromAddress, toAddress, pegHash, amount, fiatPegWallet, sourceChain, destinationChain}
}

//GetSignBytes : get bytes to sign
func (in SendFiat) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		FromAddress      string            `json:"fromAddress"`
		ToAddress        string            `json:"toAddress"`
		PegHash          sdk.PegHash       `json:"pegHash"`
		Amount           int64             `json:"amount"`
		FiatPegWallet    sdk.FiatPegWallet `json:"fiatPegWallet"`
		SourceChain      string            `json:"sourceChain"`
		DestinationChain string            `json:"destinationChain"`
	}{
		FromAddress:      in.FromAddress.String(),
		ToAddress:        in.ToAddress.String(),
		Amount:           in.Amount,
		FiatPegWallet:    in.FiatPegWallet,
		PegHash:          in.PegHash,
		SourceChain:      in.SourceChain,
		DestinationChain: in.DestinationChain,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//#####SendFiat

//*****MsgSendFiats

//MsgSendFiats : high level issuance of fiats module
type MsgSendFiats struct {
	SendFiats []SendFiat `json:"sendFiats"`
}

//NewMsgSendFiats : initializer
func NewMsgSendFiats(sendFiats []SendFiat) MsgSendFiats {
	return MsgSendFiats{sendFiats}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgSendFiats{}

//Type : implements msg
func (msg MsgSendFiats) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgSendFiats) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.SendFiats {
		if len(in.FromAddress) == 0 {
			return sdk.ErrInvalidAddress(in.FromAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		} else if in.Amount <= 0 {
			return sdk.ErrUnknownRequest("Amount should be Positive")
		} else if len(in.FiatPegWallet) == 0 {
			return sdk.ErrUnknownRequest("FiatPegWallet is Empty")
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgSendFiats) GetSignBytes() []byte {
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

//GetSigners : implements msg
func (msg MsgSendFiats) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SendFiats))
	for i, in := range msg.SendFiats {
		addrs[i] = in.FromAddress
	}
	return addrs
}

//BuildSendFiatMsg : build the sendFiatTx
func BuildSendFiatMsg(fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, amount int64, fiatPegWallet sdk.FiatPegWallet, sourceChain string, destinationChain string) sdk.Msg {
	sendFiat := NewSendFiat(fromAddress, toAddress, pegHash, amount, fiatPegWallet, sourceChain, destinationChain)
	msg := NewMsgSendFiats([]SendFiat{sendFiat})
	return msg
}

//##### Implement sdk.Msg

//#####MsgSendFiats

//*****MsgRelaySendFiats

//MsgRelaySendFiats :
type MsgRelaySendFiats struct {
	SendFiats []SendFiat     `json:"sendFiats"`
	Relayer   sdk.AccAddress `json:"relayer"`
	Sequence  int64          `json:"sequence"`
}

//NewMsgRelaySendFiats : initializer
func NewMsgRelaySendFiats(sendFiats []SendFiat, relayer sdk.AccAddress, sequence int64) MsgRelaySendFiats {
	return MsgRelaySendFiats{sendFiats, relayer, sequence}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgRelaySendFiats{}

//Type : implements msg
func (msg MsgRelaySendFiats) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgRelaySendFiats) ValidateBasic() sdk.Error {
	if len(msg.Relayer) == 0 {
		return sdk.ErrInvalidAddress(msg.Relayer.String())
	}
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.SendFiats {
		if len(in.FromAddress) == 0 {
			return sdk.ErrInvalidAddress(in.FromAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is empty")
		} else if in.Amount <= 0 {
			return sdk.ErrUnknownRequest("Amount should be Positive")
		} else if len(in.FiatPegWallet) == 0 {
			return sdk.ErrUnknownRequest("FiatPegWallet is Empty")
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgRelaySendFiats) GetSignBytes() []byte {
	var sendFiats []json.RawMessage
	for _, sendFiat := range msg.SendFiats {
		sendFiats = append(sendFiats, sendFiat.GetSignBytes())
	}

	bz, err := msgCdc.MarshalJSON(struct {
		SendFiats []json.RawMessage `json:"sendFiats"`
		Relayer   string            `json:"relayer"`
		Sequence  int64             `json:"sequence"`
	}{
		SendFiats: sendFiats,
		Relayer:   msg.Relayer.String(),
		Sequence:  msg.Sequence,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

//GetSigners : implements msg
func (msg MsgRelaySendFiats) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}

//BuildRelaySendFiatMsg : build the issueFiatTx
func BuildRelaySendFiatMsg(fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, amount int64, fiatPegWallet sdk.FiatPegWallet, sourceChain string, destinationChain string, relayer sdk.AccAddress, sequence int64) sdk.Msg {
	sendFiat := NewSendFiat(fromAddress, toAddress, pegHash, amount, fiatPegWallet, sourceChain, destinationChain)
	msg := NewMsgRelaySendFiats([]SendFiat{sendFiat}, relayer, sequence)
	return msg
}

func toHubMsgSendFiats(hubMsgSendFiats MsgSendFiats) bank.MsgBankSendFiats {
	var newSendFiats []bank.SendFiat
	for _, sendFiat := range hubMsgSendFiats.SendFiats {
		newSendFiats = append(newSendFiats, bank.NewSendFiat(sendFiat.FromAddress, sendFiat.ToAddress, sendFiat.PegHash, sendFiat.Amount))
	}
	return bank.NewMsgBankSendFiats(newSendFiats)
}

//toTypeMsgSendFiats : from type msgIssueFiat to fiatFactory msgIssueFiat
func toTypeMsgSendFiats(bankMsgSendFiats MsgSendFiats, bankMsgRelaySendFiat MsgRelaySendFiats) fiatFactory.MsgFactorySendFiats {
	var newSendFiats []fiatFactory.SendFiat
	for _, sendFiat := range bankMsgRelaySendFiat.SendFiats {
		newSendFiats = append(newSendFiats, fiatFactory.NewSendFiat(bankMsgRelaySendFiat.Relayer, sendFiat.FromAddress, sendFiat.ToAddress, sendFiat.PegHash, sendFiat.FiatPegWallet))
	}
	return fiatFactory.NewMsgFactorySendFiats(newSendFiats)
}

//#####MsgRelaySendFiats

//BuyerExecuteOrder
type BuyerExecuteOrder struct {
	MediatorAddress  sdk.AccAddress    `json:"mediatorAddress"`
	BuyerAddress     sdk.AccAddress    `json:"buyerAddress"`
	SellerAddress    sdk.AccAddress    `json:"sellerAddress"`
	PegHash          sdk.PegHash       `json:"pegHash"`
	FiatPegWallet    sdk.FiatPegWallet `json:"fiatPegWallet"`
	FiatProofHash    string            `json:"fiatProofHash"`
	SourceChain      string            `json:"sourceChain"`
	DestinationChain string            `json:"destinationChain"`
}

// NewBuyerExecuteOrder : initialization
func NewBuyerExecuteOrder(mediatorAddress sdk.AccAddress, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, fiatProofHash string, fiatPegWallet sdk.FiatPegWallet, sourceChain string, destChain string) BuyerExecuteOrder {
	return BuyerExecuteOrder{mediatorAddress, buyerAddress, sellerAddress, pegHash, fiatPegWallet, fiatProofHash, sourceChain, destChain}
}

//GetSignBytes :
func (in BuyerExecuteOrder) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		MediatorAddress  string            `json:"mediatorAddress"`
		BuyerAddress     string            `json:"buyerAddress"`
		SellerAddress    string            `json:"sellerAddress:`
		PegHash          string            `json:"pegHash"`
		FiatPegWallet    sdk.FiatPegWallet `json:"fiatPegWallet"`
		FiatProofHash    string            `json:"fiatProofHash"`
		SourceChain      string            `json:"sourceChain"`
		DestinationChain string            `json:"destinationChain"`
	}{
		MediatorAddress:  in.MediatorAddress.String(),
		BuyerAddress:     in.BuyerAddress.String(),
		SellerAddress:    in.SellerAddress.String(),
		PegHash:          in.PegHash.String(),
		FiatPegWallet:    in.FiatPegWallet,
		FiatProofHash:    in.FiatProofHash,
		SourceChain:      in.SourceChain,
		DestinationChain: in.DestinationChain,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//MsgBuyerExecuteOrders :
type MsgBuyerExecuteOrders struct {
	BuyerExecuteOrders []BuyerExecuteOrder `json:"buyerExecuteOrders"`
}

//NewMsgBuyerExecuteOrders :
func NewMsgBuyerExecuteOrders(buyerExecuteOrders []BuyerExecuteOrder) MsgBuyerExecuteOrders {
	return MsgBuyerExecuteOrders{buyerExecuteOrders}
}

var _ sdk.Msg = MsgBuyerExecuteOrders{}

//Type :
func (msg MsgBuyerExecuteOrders) Type() string { return "ibc" }

// ValidateBasic :
func (msg MsgBuyerExecuteOrders) ValidateBasic() sdk.Error {
	for _, in := range msg.BuyerExecuteOrders {
		if len(in.MediatorAddress) == 0 {
			return sdk.ErrInvalidAddress(in.MediatorAddress.String())
		} else if len(in.BuyerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.BuyerAddress.String())
		} else if len(in.SellerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.SellerAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("Invalid PegHash")
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		} else if in.FiatProofHash == "" {
			return sdk.ErrUnknownRequest("fiatProofHash is empty")
		} else if len(in.FiatPegWallet) == 0 {
			return sdk.ErrUnknownRequest("fiatPegWallet is empty")
		}
	}
	return nil
}

//GetSignBytes :
func (msg MsgBuyerExecuteOrders) GetSignBytes() []byte {
	var buyerExecuteOrders []json.RawMessage
	for _, buyerExecuteOrder := range msg.BuyerExecuteOrders {
		buyerExecuteOrders = append(buyerExecuteOrders, buyerExecuteOrder.GetSignBytes())
	}
	bz, err := msgCdc.MarshalJSON(struct {
		BuyerExecuteOrders []json.RawMessage `json:"buyerExecuteOrders"`
	}{
		BuyerExecuteOrders: buyerExecuteOrders,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners :
func (msg MsgBuyerExecuteOrders) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.BuyerExecuteOrders))
	for i, in := range msg.BuyerExecuteOrders {
		addrs[i] = in.MediatorAddress
	}
	return addrs
}

//BuildBuyerExecuteOrder : build msg
func BuildBuyerExecuteOrder(mediatorAddress sdk.AccAddress, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, fiatProofHash string, fiatPegWallet sdk.FiatPegWallet, souceChain string, destinationChain string) sdk.Msg {
	buyerExecuteOrder := NewBuyerExecuteOrder(mediatorAddress, buyerAddress, sellerAddress, pegHash, fiatProofHash, fiatPegWallet, souceChain, destinationChain)
	msg := NewMsgBuyerExecuteOrders([]BuyerExecuteOrder{buyerExecuteOrder})
	return msg
}

// MsgRelayBuyerExecuteOrders :
type MsgRelayBuyerExecuteOrders struct {
	BuyerExecuteOrders []BuyerExecuteOrder `json:"buyerExecuteOrder"`
	Relayer            sdk.AccAddress      `json:"relayerAddress"`
	Sequence           int64               `json:"sequence"`
}

// NewMsgRelayBuyerExecuteOrders :
func NewMsgRelayBuyerExecuteOrders(buyerExecuteOrders []BuyerExecuteOrder, relayerAddress sdk.AccAddress, sequence int64) MsgRelayBuyerExecuteOrders {
	return MsgRelayBuyerExecuteOrders{buyerExecuteOrders, relayerAddress, sequence}
}

func toHubMsgBuyerExecuteOrdermsg(hubMsgBuyerExecuteOrder MsgBuyerExecuteOrders) bank.MsgBankBuyerExecuteOrders {
	var newBuyerExecuteOrders []bank.BuyerExecuteOrder
	for _, buyerExecuteOrder := range hubMsgBuyerExecuteOrder.BuyerExecuteOrders {
		newBuyerExecuteOrders = append(newBuyerExecuteOrders, bank.NewBuyerExecuteOrder(buyerExecuteOrder.MediatorAddress, buyerExecuteOrder.BuyerAddress, buyerExecuteOrder.SellerAddress, buyerExecuteOrder.PegHash, buyerExecuteOrder.FiatProofHash))
	}
	return bank.NewMsgBankBuyerExecuteOrders(newBuyerExecuteOrders)
}

var _ sdk.Msg = MsgRelayBuyerExecuteOrders{}

//Type :
func (msg MsgRelayBuyerExecuteOrders) Type() string { return "ibc" }

// ValidateBasic :
func (msg MsgRelayBuyerExecuteOrders) ValidateBasic() sdk.Error {
	for _, in := range msg.BuyerExecuteOrders {
		if len(in.MediatorAddress) == 0 {
			return sdk.ErrInvalidAddress(in.MediatorAddress.String())
		} else if len(in.BuyerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.BuyerAddress.String())
		} else if len(in.SellerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.SellerAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash in Empty")
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		} else if in.FiatProofHash == "" {
			return sdk.ErrUnknownRequest("fiatProofHash is empty")
		} else if len(in.FiatPegWallet) == 0 {
			return sdk.ErrUnknownRequest("fiatPegWallet is empty")
		}
	}
	if len(msg.Relayer) == 0 {
		return sdk.ErrInvalidAddress(msg.Relayer.String())
	}
	return nil
}

// GetSignBytes :
func (msg MsgRelayBuyerExecuteOrders) GetSignBytes() []byte {
	var buyerExecuteOrders []json.RawMessage
	for _, buyerExecuteOrder := range msg.BuyerExecuteOrders {
		buyerExecuteOrders = append(buyerExecuteOrders, buyerExecuteOrder.GetSignBytes())
	}

	bz, err := msgCdc.MarshalJSON(struct {
		BuyerExecuteOrders []json.RawMessage `json:"sendFiats"`
		Relayer            string            `json:"relayer"`
		Sequence           int64             `json:"sequence"`
	}{
		BuyerExecuteOrders: buyerExecuteOrders,
		Relayer:            msg.Relayer.String(),
		Sequence:           msg.Sequence,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

//GetSigners :
func (msg MsgRelayBuyerExecuteOrders) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}

//#####MsgBuyerExecuteOrders

//*****SellerExecuteOrder

//SellerExecuteOrder - transaction input
type SellerExecuteOrder struct {
	MediatorAddress  sdk.AccAddress `json:"mediatorAddress"`
	BuyerAddress     sdk.AccAddress `json:"buyerAddress"`
	SellerAddress    sdk.AccAddress `json:"sellerAddress"`
	PegHash          sdk.PegHash    `json:"pegHash"`
	AWBProofHash     string         `json:"awbProofHash"`
	SourceChain      string         `json:"sourceChain"`
	DestinationChain string         `json:"destinationChain"`
}

//NewSellerExecuteOrder : initializer
func NewSellerExecuteOrder(mediatorAddress sdk.AccAddress, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, awbProofHash string, sourceChain string, destinationChain string) SellerExecuteOrder {
	return SellerExecuteOrder{mediatorAddress, buyerAddress, sellerAddress, pegHash, awbProofHash, sourceChain, destinationChain}
}

//GetSignBytes : get bytes to sign
func (in SellerExecuteOrder) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		MediatorAddress  string      `json:"mediatorAddress"`
		BuyerAddress     string      `json:"buyerAddress"`
		SellerAddress    string      `json:"sellerAddress"`
		PegHash          sdk.PegHash `json:"pegHash"`
		AWBProofHash     string      `json:"awbProofHash"`
		SourceChain      string      `json:"sourceChain"`
		DestinationChain string      `json:"destinationChain"`
	}{
		MediatorAddress:  in.MediatorAddress.String(),
		BuyerAddress:     in.BuyerAddress.String(),
		SellerAddress:    in.SellerAddress.String(),
		PegHash:          in.PegHash,
		AWBProofHash:     in.AWBProofHash,
		SourceChain:      in.SourceChain,
		DestinationChain: in.DestinationChain,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//#####SellerExecuteOrder

//*****MsgSellerExecuteOrders

//MsgSellerExecuteOrders : high level issuance of fiats module
type MsgSellerExecuteOrders struct {
	SellerExecuteOrders []SellerExecuteOrder `json:"sellerExecuteOrders"`
}

//NewMsgSellerExecuteOrders : initializer
func NewMsgSellerExecuteOrders(sellerExecuteOrders []SellerExecuteOrder) MsgSellerExecuteOrders {
	return MsgSellerExecuteOrders{sellerExecuteOrders}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgSellerExecuteOrders{}

//Type : implements msg
func (msg MsgSellerExecuteOrders) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgSellerExecuteOrders) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInvalidAddress(err.Error())
	}
	for _, in := range msg.SellerExecuteOrders {
		if len(in.MediatorAddress) == 0 {
			return sdk.ErrInvalidAddress(in.MediatorAddress.String())
		} else if len(in.SellerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.SellerAddress.String())
		} else if len(in.SellerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.SellerAddress.String())
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		} else if in.AWBProofHash == "" {
			return sdk.ErrUnknownRequest("ABAProofHash is Empty")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgSellerExecuteOrders) GetSignBytes() []byte {
	var sellerExecuteOrders []json.RawMessage
	for _, sellerExecuteOrder := range msg.SellerExecuteOrders {
		sellerExecuteOrders = append(sellerExecuteOrders, sellerExecuteOrder.GetSignBytes())
	}

	b, err := msgCdc.MarshalJSON(struct {
		SellerExecuteOrders []json.RawMessage `json:"sellerExecuteOrders"`
	}{
		SellerExecuteOrders: sellerExecuteOrders,
	})
	if err != nil {
		panic(err)
	}
	return b
}

//GetSigners : implements msg
func (msg MsgSellerExecuteOrders) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SellerExecuteOrders))
	for i, in := range msg.SellerExecuteOrders {
		addrs[i] = in.MediatorAddress
	}
	return addrs
}

//BuildSellerExecuteOrder : build msg
func BuildSellerExecuteOrder(mediatorAddress sdk.AccAddress, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, awbProofHash string, souceChain string, destinationChain string) sdk.Msg {
	sellerExecuteOrder := NewSellerExecuteOrder(mediatorAddress, buyerAddress, sellerAddress, pegHash, awbProofHash, souceChain, destinationChain)
	msg := NewMsgSellerExecuteOrders([]SellerExecuteOrder{sellerExecuteOrder})
	return msg
}

func toHubMsgSellerExecuteOrdermsg(hubMsgSellerExecuteOrder MsgSellerExecuteOrders) bank.MsgBankSellerExecuteOrders {
	var newSellerExecuteOrders []bank.SellerExecuteOrder
	for _, sellerExecuteOrder := range hubMsgSellerExecuteOrder.SellerExecuteOrders {
		newSellerExecuteOrders = append(newSellerExecuteOrders, bank.NewSellerExecuteOrder(sellerExecuteOrder.MediatorAddress, sellerExecuteOrder.BuyerAddress, sellerExecuteOrder.SellerAddress, sellerExecuteOrder.PegHash, sellerExecuteOrder.AWBProofHash))
	}
	return bank.NewMsgBankSellerExecuteOrders(newSellerExecuteOrders)
}

// MsgRelaySellerExecuteOrders :
type MsgRelaySellerExecuteOrders struct {
	SellerExecuteOrders []SellerExecuteOrder `json:"sellerExecuteOrder"`
	Relayer             sdk.AccAddress       `json:"relayerAddress"`
	Sequence            int64                `json:"sequence"`
}

// NewMsgRelaySellerExecuteOrders :
func NewMsgRelaySellerExecuteOrders(sellerExecuteOrders []SellerExecuteOrder, relayerAddress sdk.AccAddress, sequence int64) MsgRelaySellerExecuteOrders {
	return MsgRelaySellerExecuteOrders{sellerExecuteOrders, relayerAddress, sequence}
}

var _ sdk.Msg = MsgRelaySellerExecuteOrders{}

// Type :
func (msg MsgRelaySellerExecuteOrders) Type() string { return "ibc" }

// ValidateBasic :
func (msg MsgRelaySellerExecuteOrders) ValidateBasic() sdk.Error {
	for _, in := range msg.SellerExecuteOrders {
		if len(in.MediatorAddress) == 0 {
			return sdk.ErrInvalidAddress(in.MediatorAddress.String())
		} else if len(in.BuyerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.BuyerAddress.String())
		} else if len(in.SellerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.SellerAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("pegHash is empty")
		} else if len(in.AWBProofHash) == 0 {
			return sdk.ErrUnknownRequest("awbproofhash is empty")
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		}
	}
	if len(msg.Relayer) == 0 {
		return sdk.ErrInvalidAddress(msg.Relayer.String())
	}
	return nil
}

// GetSignBytes :

func (msg MsgRelaySellerExecuteOrders) GetSignBytes() []byte {
	var sellerExecuteOrders []json.RawMessage
	for _, sellerExecuteOrder := range msg.SellerExecuteOrders {
		sellerExecuteOrders = append(sellerExecuteOrders, sellerExecuteOrder.GetSignBytes())
	}

	bz, err := msgCdc.MarshalJSON(struct {
		SellerExecuteOrders []json.RawMessage `json:"sendFiats"`
		Relayer             string            `json:"relayer"`
		Sequence            int64             `json:"sequence"`
	}{
		SellerExecuteOrders: sellerExecuteOrders,
		Relayer:             msg.Relayer.String(),
		Sequence:            msg.Sequence,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners :
func (msg MsgRelaySellerExecuteOrders) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}

//##### Implement sdk.Msg

//#####MsgSellerExecuteOrders

//*****ExecuteOrder
/*
//ExecuteOrder - transaction input
type ExecuteOrder struct {
	MediatorAddress  sdk.AccAddress `json:"mediatorAddress"`
	BuyerAddress     sdk.AccAddress `json:"buyerAddress"`
	SellerAddress    sdk.AccAddress `json:"sellerAddress"`
	PegHash          sdk.PegHash    `json:"pegHash"`
	SourceChain      string         `json:"sourceChain"`
	DestinationChain string         `json:"destinationChain"`
}

//NewExecuteOrder : initializer
func NewExecuteOrder(mediatorAddress sdk.AccAddress, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, sourceChain string, destinationChain string) ExecuteOrder {
	return ExecuteOrder{mediatorAddress, buyerAddress, sellerAddress, pegHash, sourceChain, destinationChain}
}

//GetSignBytes : get bytes to sign
func (in ExecuteOrder) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		MediatorAddress  string `json:"mediatorAddress"`
		BuyerAddress     string `json:"buyerAddress"`
		SellerAddress    string `json:"sellerAddress"`
		PegHash          string `json:"pegHash"`
		SourceChain      string `json:"sourceChain"`
		DestinationChain string `json:"destinationChain"`
	}{
		MediatorAddress:  in.MediatorAddress.String(),
		BuyerAddress:     in.BuyerAddress.String(),
		SellerAddress:    in.SellerAddress.String(),
		PegHash:          in.PegHash.String(),
		SourceChain:      in.SourceChain,
		DestinationChain: in.DestinationChain,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//#####ExecuteOrder

//*****MsgExecuteOrders

//MsgExecuteOrders : high level issuance of fiats module
type MsgExecuteOrders struct {
	ExecuteOrders []ExecuteOrder `json:"executeOrders"`
}

//NewMsgExecuteOrders : initializer
func NewMsgExecuteOrders(executeOrders []ExecuteOrder) MsgExecuteOrders {
	return MsgExecuteOrders{executeOrders}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgExecuteOrders{}

//Type : implements msg
func (msg MsgExecuteOrders) Type() string { return "ibc" }

//ValidateBasic : implements msg
func (msg MsgExecuteOrders) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInvalidAddress(err.Error())
	}
	for _, in := range msg.ExecuteOrders {
		if len(in.MediatorAddress) == 0 {
			return sdk.ErrInvalidAddress(in.MediatorAddress.String())
		} else if len(in.SellerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.SellerAddress.String())
		} else if len(in.BuyerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.BuyerAddress.String())
		} else if in.SourceChain == in.DestinationChain {
			return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgExecuteOrders) GetSignBytes() []byte {
	var executeOrders []json.RawMessage
	for _, executeOrder := range msg.ExecuteOrders {
		executeOrders = append(executeOrders, executeOrder.GetSignBytes())
	}

	b, err := msgCdc.MarshalJSON(struct {
		ExecuteOrders []json.RawMessage `json:"executeOrders"`
	}{
		ExecuteOrders: executeOrders,
	})
	if err != nil {
		panic(err)
	}
	return b
}

//GetSigners : implements msg
func (msg MsgExecuteOrders) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.ExecuteOrders))
	for i, in := range msg.ExecuteOrders {
		addrs[i] = in.MediatorAddress
	}
	return addrs
}

//##### Implement sdk.Msg

//#####MsgExecuteOrder

//#####Comdex
*/
