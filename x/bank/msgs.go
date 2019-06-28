package bank

import (
	"encoding/json"
	
	"github.com/asaskevich/govalidator"
	sdk "github.com/comdex-blockchain/types"
)

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

var _ sdk.Msg = MsgSend{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(in []Input, out []Output) MsgSend {
	return MsgSend{Inputs: in, Outputs: out}
}

// Type Implements Msg.
func (msg MsgSend) Type() string { return "bank" } // TODO: "bank/send"

// ValidateBasic Implements Msg.
func (msg MsgSend) ValidateBasic() sdk.Error {
	if len(msg.Inputs) == 0 {
		return ErrNoInputs(DefaultCodespace).TraceSDK("")
	}
	if len(msg.Outputs) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	var totalIn, totalOut sdk.Coins
	for _, in := range msg.Inputs {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
		totalIn = totalIn.Plus(in.Coins)
	}
	for _, out := range msg.Outputs {
		if err := out.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
		totalOut = totalOut.Plus(out.Coins)
	}
	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return sdk.ErrInvalidCoins(totalIn.String()).TraceSDK("inputs and outputs don't match")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	var inputs, outputs []json.RawMessage
	for _, input := range msg.Inputs {
		inputs = append(inputs, input.GetSignBytes())
	}
	for _, output := range msg.Outputs {
		outputs = append(outputs, output.GetSignBytes())
	}
	b, err := msgCdc.MarshalJSON(struct {
		Inputs  []json.RawMessage `json:"inputs"`
		Outputs []json.RawMessage `json:"outputs"`
	}{
		Inputs:  inputs,
		Outputs: outputs,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}

// ----------------------------------------
// MsgIssue

// MsgIssue - high level transaction of the coin module
type MsgIssue struct {
	Banker  sdk.AccAddress `json:"banker"`
	Outputs []Output       `json:"outputs"`
}

var _ sdk.Msg = MsgIssue{}

// NewMsgIssue - construct arbitrary multi-in, multi-out send msg.
func NewMsgIssue(banker sdk.AccAddress, out []Output) MsgIssue {
	return MsgIssue{Banker: banker, Outputs: out}
}

// Type Implements Msg.
func (msg MsgIssue) Type() string { return "bank" } // TODO: "bank/issue"

// ValidateBasic Implements Msg.
func (msg MsgIssue) ValidateBasic() sdk.Error {
	// XXX
	if len(msg.Outputs) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInvalidCoins(err.Error())
	}
	for _, out := range msg.Outputs {
		if err := out.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssue) GetSignBytes() []byte {
	var outputs []json.RawMessage
	for _, output := range msg.Outputs {
		outputs = append(outputs, output.GetSignBytes())
	}
	b, err := msgCdc.MarshalJSON(struct {
		Banker  sdk.AccAddress    `json:"banker"`
		Outputs []json.RawMessage `json:"outputs"`
	}{
		Banker:  msg.Banker,
		Outputs: outputs,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgIssue) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Banker}
}

// ----------------------------------------
// Input

// Input Transaction Input
type Input struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// GetSignBytes Return bytes to sign for Input
func (in Input) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(in)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bin)
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() sdk.Error {
	if len(in.Address) == 0 {
		return sdk.ErrInvalidAddress(in.Address.String())
	}
	
	if !in.Coins.IsValid() {
		return sdk.ErrInvalidCoins(in.Coins.String())
	}
	if !in.Coins.IsPositive() {
		return sdk.ErrInvalidCoins(in.Coins.String())
	}
	return nil
}

// NewInput - create a transaction input, used with MsgSend
func NewInput(addr sdk.AccAddress, coins sdk.Coins) Input {
	input := Input{
		Address: addr,
		Coins:   coins,
	}
	return input
}

// ----------------------------------------
// Output

// Output Transaction Output
type Output struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// GetSignBytes Return bytes to sign for Output
func (out Output) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(out)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bin)
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() sdk.Error {
	if len(out.Address) == 0 {
		return sdk.ErrInvalidAddress(out.Address.String())
	}
	if !out.Coins.IsValid() {
		return sdk.ErrInvalidCoins(out.Coins.String())
	}
	if !out.Coins.IsPositive() {
		return sdk.ErrInvalidCoins(out.Coins.String())
	}
	return nil
}

// NewOutput - create a transaction output, used with MsgSend
func NewOutput(addr sdk.AccAddress, coins sdk.Coins) Output {
	output := Output{
		Address: addr,
		Coins:   coins,
	}
	return output
}

// *****Comdex

// *****IssueAsset

// IssueAsset - transaction input
type IssueAsset struct {
	IssuerAddress sdk.AccAddress `json:"issuerAddress"`
	ToAddress     sdk.AccAddress `json:"toAddress"`
	AssetPeg      sdk.AssetPeg   `json:"assetPeg"`
}

// NewIssueAsset : initializer
func NewIssueAsset(issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, assetPeg sdk.AssetPeg) IssueAsset {
	return IssueAsset{issuerAddress, toAddress, assetPeg}
}

// GetSignBytes : get bytes to sign
func (in IssueAsset) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		IssuerAddress string       `json:"issuerAddress"`
		ToAddress     string       `json:"toAddress"`
		AssetPeg      sdk.AssetPeg `json:"assetPeg"`
	}{
		IssuerAddress: in.IssuerAddress.String(),
		ToAddress:     in.ToAddress.String(),
		AssetPeg:      in.AssetPeg,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// #####IssueAsset

// *****MsgBankIssueAssets

// MsgBankIssueAssets : high level issuance of assets module
type MsgBankIssueAssets struct {
	IssueAssets []IssueAsset `json:"issueAssets"`
}

// NewMsgBankIssueAssets : initilizer
func NewMsgBankIssueAssets(issueAssets []IssueAsset) MsgBankIssueAssets {
	return MsgBankIssueAssets{issueAssets}
}

// ***** Implementing sdk.Msg

var _ sdk.Msg = MsgBankIssueAssets{}

// Type : implements msg
func (msg MsgBankIssueAssets) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgBankIssueAssets) ValidateBasic() sdk.Error {
	if len(msg.IssueAssets) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.IssueAssets {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		}
		_, err := govalidator.ValidateStruct(in.AssetPeg)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankIssueAssets) GetSignBytes() []byte {
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

// GetSigners : implements msg
func (msg MsgBankIssueAssets) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.IssueAssets))
	for i, in := range msg.IssueAssets {
		addrs[i] = in.IssuerAddress
	}
	return addrs
}

// ##### Implement sdk.Msg

// #####MsgBankIssueAssets

// ****RedeemAsset

// RedeemAsset : transsction input
type RedeemAsset struct {
	IssuerAddress   sdk.AccAddress `json:"issuerAddress"`
	RedeemerAddress sdk.AccAddress `json:"redeemerAddress"`
	PegHash         sdk.PegHash    `json:"pegHash"`
}

// NewRedeemAsset : initializer
func NewRedeemAsset(issuerAddress sdk.AccAddress, redeemerAddress sdk.AccAddress, pegHash sdk.PegHash) RedeemAsset {
	return RedeemAsset{issuerAddress, redeemerAddress, pegHash}
}

// GetSignBytes : get bytes to sign
func (in RedeemAsset) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		IssuerAddress   string      `json:"issuerAddress"`
		RedeemerAddress string      `json:"redeemerAddress"`
		PegHash         sdk.PegHash `json:"pegHash"`
	}{
		IssuerAddress:   in.IssuerAddress.String(),
		RedeemerAddress: in.RedeemerAddress.String(),
		PegHash:         in.PegHash,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// #####RedeemAsset

// *****MsgBankRedeemAssets

// MsgBankRedeemAssets : Message to redeem issued assets
type MsgBankRedeemAssets struct {
	RedeemAssets []RedeemAsset `json:"redeemAssets"`
}

// NewMsgBankRedeemAssets : initializer
func NewMsgBankRedeemAssets(redeemAssets []RedeemAsset) MsgBankRedeemAssets {
	return MsgBankRedeemAssets{redeemAssets}
}

// *****Implementing sdk.Msg

var _ sdk.Msg = MsgBankRedeemAssets{}

// Type : implements msg
func (msg MsgBankRedeemAssets) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgBankRedeemAssets) ValidateBasic() sdk.Error {
	if len(msg.RedeemAssets) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.RedeemAssets {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.RedeemerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RedeemerAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankRedeemAssets) GetSignBytes() []byte {
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

// GetSigners : implements msg
func (msg MsgBankRedeemAssets) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.RedeemAssets))
	for i, in := range msg.RedeemAssets {
		addrs[i] = in.RedeemerAddress
	}
	return addrs
}

// ##### Implement sdk.Msg

// ######MsgBankRedeemAssets

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

// *****MsgBankIssueFiats

// MsgBankIssueFiats : high level issuance of fiats module
type MsgBankIssueFiats struct {
	IssueFiats []IssueFiat `json:"issueFiats"`
}

// NewMsgBankIssueFiats : initilizer
func NewMsgBankIssueFiats(issueFiats []IssueFiat) MsgBankIssueFiats {
	return MsgBankIssueFiats{issueFiats}
}

// ***** Implementing sdk.Msg

var _ sdk.Msg = MsgBankIssueFiats{}

// Type : implements msg
func (msg MsgBankIssueFiats) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgBankIssueFiats) ValidateBasic() sdk.Error {
	if len(msg.IssueFiats) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.IssueFiats {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		}
		_, err := govalidator.ValidateStruct(in.FiatPeg)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankIssueFiats) GetSignBytes() []byte {
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
func (msg MsgBankIssueFiats) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.IssueFiats))
	for i, in := range msg.IssueFiats {
		addrs[i] = in.IssuerAddress
	}
	return addrs
}

// ##### Implement sdk.Msg

// #####MsgBankIssueFiats

// ****RedeemFiat

// RedeemFiat : transaction input
type RedeemFiat struct {
	RedeemerAddress sdk.AccAddress `json:"redeemerAddress"`
	IssuerAddress   sdk.AccAddress `json:"issuerAddress"`
	Amount          int64          `json:"amount"`
}

// NewRedeemFiat : initializer
func NewRedeemFiat(redeemerAddress sdk.AccAddress, issuerAddress sdk.AccAddress, amount int64) RedeemFiat {
	return RedeemFiat{redeemerAddress, issuerAddress, amount}
}

// GetSignBytes : get bytes to sign
func (in RedeemFiat) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		RedeemerAddress string `json:"redeemerAddress"`
		IssuerAddress   string `json:"issuerAddress"`
		Amount          int64  `json:"amount"`
	}{
		RedeemerAddress: in.RedeemerAddress.String(),
		IssuerAddress:   in.IssuerAddress.String(),
		Amount:          in.Amount,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// #####RedeemFiat

// *****MsgBankRedeemFiats

// MsgBankRedeemFiats : Message to redeem issued fiats
type MsgBankRedeemFiats struct {
	RedeemFiats []RedeemFiat `json:"redeemFiats"`
}

// NewMsgBankRedeemFiats : initializer
func NewMsgBankRedeemFiats(redeemFiats []RedeemFiat) MsgBankRedeemFiats {
	return MsgBankRedeemFiats{redeemFiats}
}

// *****Implementing sdk.Msg

var _ sdk.Msg = MsgBankRedeemFiats{}

// Type : implements msg
func (msg MsgBankRedeemFiats) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgBankRedeemFiats) ValidateBasic() sdk.Error {
	if len(msg.RedeemFiats) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.RedeemFiats {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.RedeemerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RedeemerAddress.String())
		} else if in.Amount <= 0 {
			return sdk.ErrUnknownRequest("Amount should be Positive")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankRedeemFiats) GetSignBytes() []byte {
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

// GetSigners : implements msg
func (msg MsgBankRedeemFiats) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.RedeemFiats))
	for i, in := range msg.RedeemFiats {
		addrs[i] = in.RedeemerAddress
	}
	return addrs
}

// ##### Implement sdk.Msg

// ######MsgBankRedeemFiats

// *****SendAsset

// SendAsset - transaction input
type SendAsset struct {
	FromAddress sdk.AccAddress `json:"fromAddress"`
	ToAddress   sdk.AccAddress `json:"toAddress"`
	PegHash     sdk.PegHash    `json:"pegHash"`
}

// NewSendAsset : initializer
func NewSendAsset(fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash) SendAsset {
	return SendAsset{fromAddress, toAddress, pegHash}
}

// GetSignBytes : get bytes to sign
func (in SendAsset) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		FromAddress string `json:"fromAddress"`
		ToAddress   string `json:"toAddress"`
		PegHash     string `json:"pegHash"`
	}{
		FromAddress: in.FromAddress.String(),
		ToAddress:   in.ToAddress.String(),
		PegHash:     in.PegHash.String(),
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// #####SendAsset

// *****MsgBankSendAssets

// MsgBankSendAssets : high level issuance of assets module
type MsgBankSendAssets struct {
	SendAssets []SendAsset `json:"sendAssets"`
}

// NewMsgBankSendAssets : initilizer
func NewMsgBankSendAssets(sendAssets []SendAsset) MsgBankSendAssets {
	return MsgBankSendAssets{sendAssets}
}

// ***** Implementing sdk.Msg

var _ sdk.Msg = MsgBankSendAssets{}

// Type : implements msg
func (msg MsgBankSendAssets) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgBankSendAssets) ValidateBasic() sdk.Error {
	if len(msg.SendAssets) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.SendAssets {
		if len(in.FromAddress) == 0 {
			return sdk.ErrInvalidAddress(in.FromAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is empty")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankSendAssets) GetSignBytes() []byte {
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

// GetSigners : implements msg
func (msg MsgBankSendAssets) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SendAssets))
	for i, in := range msg.SendAssets {
		addrs[i] = in.FromAddress
	}
	return addrs
}

// ##### Implement sdk.Msg

// #####MsgBankSendAssets

// *****SendFiat

// SendFiat - transaction input
type SendFiat struct {
	FromAddress sdk.AccAddress `json:"fromAddress"`
	ToAddress   sdk.AccAddress `json:"toAddress"`
	PegHash     sdk.PegHash    `json:"pegHash"`
	Amount      int64          `json:"amount"`
}

// NewSendFiat : initializer
func NewSendFiat(fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, amount int64) SendFiat {
	return SendFiat{fromAddress, toAddress, pegHash, amount}
}

// GetSignBytes : get bytes to sign
func (in SendFiat) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		FromAddress string `json:"fromAddress"`
		ToAddress   string `json:"toAddress"`
		PegHash     string `json:"pegHash"`
		Amount      int64  `json:"amount"`
	}{
		FromAddress: in.FromAddress.String(),
		ToAddress:   in.ToAddress.String(),
		PegHash:     in.PegHash.String(),
		Amount:      in.Amount,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// #####SendFiat

// *****MsgBankSendFiats

// MsgBankSendFiats : high level issuance of fiats module
type MsgBankSendFiats struct {
	SendFiats []SendFiat `json:"sendFiats"`
}

// NewMsgBankSendFiats : initilizer
func NewMsgBankSendFiats(sendFiats []SendFiat) MsgBankSendFiats {
	return MsgBankSendFiats{sendFiats}
}

// ***** Implementing sdk.Msg

var _ sdk.Msg = MsgBankSendFiats{}

// Type : implements msg
func (msg MsgBankSendFiats) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgBankSendFiats) ValidateBasic() sdk.Error {
	if len(msg.SendFiats) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.SendFiats {
		if len(in.FromAddress) == 0 {
			return sdk.ErrInvalidAddress(in.FromAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		} else if in.Amount <= 0 {
			return sdk.ErrUnknownRequest("Amount should be positive")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankSendFiats) GetSignBytes() []byte {
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
func (msg MsgBankSendFiats) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SendFiats))
	for i, in := range msg.SendFiats {
		addrs[i] = in.FromAddress
	}
	return addrs
}

// ##### Implement sdk.Msg

// #####MsgBankSendFiats

// *****BuyerExecuteOrder

// BuyerExecuteOrder - transaction input
type BuyerExecuteOrder struct {
	MediatorAddress sdk.AccAddress `json:"mediatorAddress"`
	BuyerAddress    sdk.AccAddress `json:"buyerAddress"`
	SellerAddress   sdk.AccAddress `json:"sellerAddress"`
	PegHash         sdk.PegHash    `json:"pegHash"`
	FiatProofHash   string         `json:"fiatProofHash"`
}

// NewBuyerExecuteOrder : initializer
func NewBuyerExecuteOrder(mediatorAddress sdk.AccAddress, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, fiatProofHash string) BuyerExecuteOrder {
	return BuyerExecuteOrder{mediatorAddress, buyerAddress, sellerAddress, pegHash, fiatProofHash}
}

// GetSignBytes : get bytes to sign
func (in BuyerExecuteOrder) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		MediatorAddress string `json:"mediatorAddress"`
		BuyerAddress    string `json:"buyerAddress"`
		SellerAddress   string `json:"sellerAddress"`
		PegHash         string `json:"pegHash"`
		FiatProofHash   string `json:"fiatProofHash"`
	}{
		MediatorAddress: in.MediatorAddress.String(),
		BuyerAddress:    in.BuyerAddress.String(),
		SellerAddress:   in.SellerAddress.String(),
		PegHash:         in.PegHash.String(),
		FiatProofHash:   in.FiatProofHash,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// #####BuyerExecuteOrder

// *****MsgBankBuyerExecuteOrders

// MsgBankBuyerExecuteOrders : high level issuance of fiats module
type MsgBankBuyerExecuteOrders struct {
	BuyerExecuteOrders []BuyerExecuteOrder `json:"buyerExecuteOrders"`
}

// NewMsgBankBuyerExecuteOrders : initilizer
func NewMsgBankBuyerExecuteOrders(buyerExecuteOrders []BuyerExecuteOrder) MsgBankBuyerExecuteOrders {
	return MsgBankBuyerExecuteOrders{buyerExecuteOrders}
}

// ***** Implementing sdk.Msg

var _ sdk.Msg = MsgBankBuyerExecuteOrders{}

// Type : implements msg
func (msg MsgBankBuyerExecuteOrders) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgBankBuyerExecuteOrders) ValidateBasic() sdk.Error {
	if len(msg.BuyerExecuteOrders) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.BuyerExecuteOrders {
		if len(in.MediatorAddress) == 0 {
			return sdk.ErrInvalidAddress(in.MediatorAddress.String())
		} else if len(in.SellerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.SellerAddress.String())
		} else if len(in.BuyerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.BuyerAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		} else if in.FiatProofHash == "" {
			return sdk.ErrUnknownRequest("FiatProofHash is Empty")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankBuyerExecuteOrders) GetSignBytes() []byte {
	var buyerExecuteOrders []json.RawMessage
	for _, buyerExecuteOrder := range msg.BuyerExecuteOrders {
		buyerExecuteOrders = append(buyerExecuteOrders, buyerExecuteOrder.GetSignBytes())
	}
	b, err := msgCdc.MarshalJSON(struct {
		BuyerExecuteOrders []json.RawMessage `json:"buyerExecuteOrders"`
	}{
		BuyerExecuteOrders: buyerExecuteOrders,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgBankBuyerExecuteOrders) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.BuyerExecuteOrders))
	for i, in := range msg.BuyerExecuteOrders {
		addrs[i] = in.MediatorAddress
	}
	return addrs
}

// ##### Implement sdk.Msg

// #####MsgBankBuyerExecuteOrders

// *****SellerExecuteOrder

// SellerExecuteOrder - transaction input
type SellerExecuteOrder struct {
	MediatorAddress sdk.AccAddress `json:"mediatorAddress"`
	BuyerAddress    sdk.AccAddress `json:"buyerAddress"`
	SellerAddress   sdk.AccAddress `json:"sellerAddress"`
	PegHash         sdk.PegHash    `json:"pegHash"`
	AWBProofHash    string         `json:"awbProofHash"`
}

// NewSellerExecuteOrder : initializer
func NewSellerExecuteOrder(mediatorAddress sdk.AccAddress, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, awbProofHash string) SellerExecuteOrder {
	return SellerExecuteOrder{mediatorAddress, buyerAddress, sellerAddress, pegHash, awbProofHash}
}

// GetSignBytes : get bytes to sign
func (in SellerExecuteOrder) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		MediatorAddress string `json:"mediatorAddress"`
		BuyerAddress    string `json:"buyerAddress"`
		SellerAddress   string `json:"sellerAddress"`
		PegHash         string `json:"pegHash"`
		AWBProofHash    string `json:"awbProofHash"`
	}{
		MediatorAddress: in.MediatorAddress.String(),
		BuyerAddress:    in.BuyerAddress.String(),
		SellerAddress:   in.SellerAddress.String(),
		PegHash:         in.PegHash.String(),
		AWBProofHash:    in.AWBProofHash,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// #####SellerExecuteOrder

// *****MsgBankSellerExecuteOrders

// MsgBankSellerExecuteOrders : high level issuance of fiats module
type MsgBankSellerExecuteOrders struct {
	SellerExecuteOrders []SellerExecuteOrder `json:"sellerExecuteOrders"`
}

// NewMsgBankSellerExecuteOrders : initilizer
func NewMsgBankSellerExecuteOrders(sellerExecuteOrders []SellerExecuteOrder) MsgBankSellerExecuteOrders {
	return MsgBankSellerExecuteOrders{sellerExecuteOrders}
}

// ***** Implementing sdk.Msg

var _ sdk.Msg = MsgBankSellerExecuteOrders{}

// Type : implements msg
func (msg MsgBankSellerExecuteOrders) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgBankSellerExecuteOrders) ValidateBasic() sdk.Error {
	if len(msg.SellerExecuteOrders) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.SellerExecuteOrders {
		if len(in.MediatorAddress) == 0 {
			return sdk.ErrInvalidAddress(in.MediatorAddress.String())
		} else if len(in.SellerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.SellerAddress.String())
		} else if len(in.BuyerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.BuyerAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		} else if in.AWBProofHash == "" {
			return sdk.ErrUnknownRequest("ABAProofHash is Empty")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankSellerExecuteOrders) GetSignBytes() []byte {
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

// GetSigners : implements msg
func (msg MsgBankSellerExecuteOrders) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SellerExecuteOrders))
	for i, in := range msg.SellerExecuteOrders {
		addrs[i] = in.MediatorAddress
	}
	return addrs
}

// ##### Implement sdk.Msg

// #####MsgBankSellerExecuteOrders

// *****ReleaseAsset

// ReleaseAsset - transaction input
type ReleaseAsset struct {
	ZoneAddress  sdk.AccAddress `json:"zoneAddress"`
	OwnerAddress sdk.AccAddress `json:"ownerAddress"`
	PegHash      sdk.PegHash    `json:"pegHash"`
}

// NewReleaseAsset : initializer
func NewReleaseAsset(zoneAddress sdk.AccAddress, ownerAddress sdk.AccAddress, pegHash sdk.PegHash) ReleaseAsset {
	return ReleaseAsset{zoneAddress, ownerAddress, pegHash}
}

// GetSignBytes : get bytes to sign
func (in ReleaseAsset) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		ZoneAddress  string `json:"zoneAddress"`
		OwnerAddress string `json:"ownerAddress"`
		PegHash      string `json:"pegHash"`
	}{
		ZoneAddress:  in.ZoneAddress.String(),
		OwnerAddress: in.OwnerAddress.String(),
		PegHash:      in.PegHash.String(),
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// #####ReleaseAsset

// *****MsgBankReleaseAssets

// MsgBankReleaseAssets : high level release of asset module
type MsgBankReleaseAssets struct {
	ReleaseAssets []ReleaseAsset `json:"releseAssets"`
}

// NewMsgBankReleaseAssets : initilizer
func NewMsgBankReleaseAssets(releseAsset []ReleaseAsset) MsgBankReleaseAssets {
	return MsgBankReleaseAssets{releseAsset}
}

// ***** Implementing sdk.Msg

var _ sdk.Msg = MsgBankReleaseAssets{}

// Type : implements msg
func (msg MsgBankReleaseAssets) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgBankReleaseAssets) ValidateBasic() sdk.Error {
	if len(msg.ReleaseAssets) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.ReleaseAssets {
		if len(in.OwnerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.OwnerAddress.String())
		} else if len(in.ZoneAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ZoneAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgBankReleaseAssets) GetSignBytes() []byte {
	var releaseAssets []json.RawMessage
	for _, releaseAsset := range msg.ReleaseAssets {
		releaseAssets = append(releaseAssets, releaseAsset.GetSignBytes())
	}
	
	b, err := msgCdc.MarshalJSON(struct {
		ReleaseAssets []json.RawMessage `json:"releaseAssets"`
	}{
		ReleaseAssets: releaseAssets,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgBankReleaseAssets) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.ReleaseAssets))
	for i, in := range msg.ReleaseAssets {
		addrs[i] = in.ZoneAddress
	}
	return addrs
}

// ##### Implement sdk.Msg

// #####MsgBankReleaseAssets

// DefineZone : singular define zone message
// *****ACL
type DefineZone struct {
	From   sdk.AccAddress `json:"from"`
	To     sdk.AccAddress `json:"to"`
	ZoneID sdk.ZoneID     `json:"zoneID"`
}

// NewDefineZone : new define zone struct
func NewDefineZone(from sdk.AccAddress, to sdk.AccAddress, zoneID sdk.ZoneID) DefineZone {
	return DefineZone{from, to, zoneID}
}

// GetSignBytes : get bytes to sign
func (in DefineZone) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		From   string `json:"from"`
		To     string `json:"to"`
		ZoneID string `json:"zoneID"`
	}{
		From:   in.From.String(),
		To:     in.To.String(),
		ZoneID: in.ZoneID.String(),
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// ValidateBasic : Validate Basic
func (in DefineZone) ValidateBasic() sdk.Error {
	if len(in.From) == 0 {
		return sdk.ErrInvalidAddress(in.From.String())
	} else if len(in.To) == 0 {
		return sdk.ErrInvalidAddress(in.To.String())
	} else if len(in.ZoneID) == 0 {
		return sdk.ErrInvalidAddress(in.ZoneID.String())
	}
	return nil
}

// MsgDefineZones : message define zones
type MsgDefineZones struct {
	DefineZones []DefineZone `json:"defineZones"`
}

// NewMsgDefineZones : new message define zones
func NewMsgDefineZones(defineZones []DefineZone) MsgDefineZones {
	return MsgDefineZones{defineZones}
}

var _ sdk.Msg = MsgDefineZones{}

// Type : implements msg
func (msg MsgDefineZones) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgDefineZones) ValidateBasic() sdk.Error {
	if len(msg.DefineZones) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.DefineZones {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgDefineZones) GetSignBytes() []byte {
	var defineZones []json.RawMessage
	for _, defineZone := range msg.DefineZones {
		defineZones = append(defineZones, defineZone.GetSignBytes())
	}
	
	b, err := msgCdc.MarshalJSON(struct {
		DefineZones []json.RawMessage `json:"defineZones"`
	}{
		DefineZones: defineZones,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgDefineZones) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.DefineZones))
	for i, in := range msg.DefineZones {
		addrs[i] = in.From
	}
	return addrs
}

// BuildMsgDefineZones : build define zones message
func BuildMsgDefineZones(from sdk.AccAddress, to sdk.AccAddress, zoneID sdk.ZoneID, msgs []DefineZone) []DefineZone {
	defineZone := NewDefineZone(from, to, zoneID)
	msgs = append(msgs, defineZone)
	return msgs
}

// BuildMsgDefineZoneWithDefineZones : build define zones message
func BuildMsgDefineZoneWithDefineZones(msgs []DefineZone) sdk.Msg {
	return NewMsgDefineZones(msgs)
}

// BuildMsgDefineZone : build define zones message
func BuildMsgDefineZone(from sdk.AccAddress, to sdk.AccAddress, zoneID sdk.ZoneID) sdk.Msg {
	defineZone := NewDefineZone(from, to, zoneID)
	return NewMsgDefineZones([]DefineZone{defineZone})
}

// DefineOrganization : singular define organization message
type DefineOrganization struct {
	From           sdk.AccAddress     `json:"from"`
	To             sdk.AccAddress     `json:"to"`
	OrganizationID sdk.OrganizationID `json:"organizationID"`
	ZoneID         sdk.ZoneID         `json:"zoneID"`
}

// NewDefineOrganization : new define organization struct
func NewDefineOrganization(from sdk.AccAddress, to sdk.AccAddress, organizationID sdk.OrganizationID, zoneID sdk.ZoneID) DefineOrganization {
	return DefineOrganization{from, to, organizationID, zoneID}
}

// GetSignBytes : get bytes to sign
func (in DefineOrganization) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		From           string `json:"from"`
		To             string `json:"to"`
		OrganizationID string `json:"organizationID"`
		ZoneID         string `json:"zoneID"`
	}{
		From:           in.From.String(),
		To:             in.To.String(),
		OrganizationID: in.OrganizationID.String(),
		ZoneID:         in.ZoneID.String(),
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// ValidateBasic : Validate Basic
func (in DefineOrganization) ValidateBasic() sdk.Error {
	if len(in.From) == 0 {
		return sdk.ErrInvalidAddress(in.From.String())
	} else if len(in.To) == 0 {
		return sdk.ErrInvalidAddress(in.To.String())
	} else if len(in.OrganizationID) == 0 {
		return sdk.ErrInvalidAddress(in.OrganizationID.String())
	} else if len(in.ZoneID) == 0 {
		return sdk.ErrInvalidAddress(in.ZoneID.String())
	}
	return nil
}

// MsgDefineOrganizations : message define organizations
type MsgDefineOrganizations struct {
	DefineOrganizations []DefineOrganization `json:"defineOrganizations"`
}

// NewMsgDefineOrganizations : new message define organizations
func NewMsgDefineOrganizations(defineOrganizations []DefineOrganization) MsgDefineOrganizations {
	return MsgDefineOrganizations{defineOrganizations}
}

var _ sdk.Msg = MsgDefineOrganizations{}

// Type : implements msg
func (msg MsgDefineOrganizations) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgDefineOrganizations) ValidateBasic() sdk.Error {
	if len(msg.DefineOrganizations) == 0 {
		return ErrNoInputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.DefineOrganizations {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgDefineOrganizations) GetSignBytes() []byte {
	var defineOrganizations []json.RawMessage
	for _, defineOrganization := range msg.DefineOrganizations {
		defineOrganizations = append(defineOrganizations, defineOrganization.GetSignBytes())
	}
	
	b, err := msgCdc.MarshalJSON(struct {
		DefineOrganizations []json.RawMessage `json:"defineOrganizations"`
	}{
		DefineOrganizations: defineOrganizations,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgDefineOrganizations) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.DefineOrganizations))
	for i, in := range msg.DefineOrganizations {
		addrs[i] = in.From
	}
	return addrs
}

// BuildMsgDefineOrganizations : build define organization message
func BuildMsgDefineOrganizations(from sdk.AccAddress, to sdk.AccAddress, organizationID sdk.OrganizationID, zoneID sdk.ZoneID, msgs []DefineOrganization) []DefineOrganization {
	defineOrganization := NewDefineOrganization(from, to, organizationID, zoneID)
	msgs = append(msgs, defineOrganization)
	return msgs
}

// BuildMsgDefineOrganizationWithMsgs : build define organization message
func BuildMsgDefineOrganizationWithMsgs(msgs []DefineOrganization) sdk.Msg {
	return NewMsgDefineOrganizations(msgs)
}

// BuildMsgDefineOrganization : build define organization message
func BuildMsgDefineOrganization(from sdk.AccAddress, to sdk.AccAddress, organizationID sdk.OrganizationID, zoneID sdk.ZoneID) sdk.Msg {
	defineOrganization := NewDefineOrganization(from, to, organizationID, zoneID)
	return NewMsgDefineOrganizations([]DefineOrganization{defineOrganization})
}

// DefineACL : indular define acl message
type DefineACL struct {
	From       sdk.AccAddress `json:"from"`
	To         sdk.AccAddress `json:"to"`
	ACLAccount sdk.ACLAccount `json:"aclAccount"`
}

// NewDefineACL : new define acl struct
func NewDefineACL(from sdk.AccAddress, to sdk.AccAddress, aclAccount sdk.ACLAccount) DefineACL {
	return DefineACL{from, to, aclAccount}
}

// GetSignBytes : get bytes to sign
func (in DefineACL) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		From       string         `json:"from"`
		To         string         `json:"to"`
		ACLAccount sdk.ACLAccount `json:"aclAccount"`
	}{
		From:       in.From.String(),
		To:         in.To.String(),
		ACLAccount: in.ACLAccount,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// ValidateBasic : Validate Basic
func (in DefineACL) ValidateBasic() sdk.Error {
	if len(in.From) == 0 {
		return sdk.ErrInvalidAddress(in.From.String())
	} else if len(in.To) == 0 {
		return sdk.ErrInvalidAddress(in.To.String())
	}
	return nil
}

// MsgDefineACLs : message define acls
type MsgDefineACLs struct {
	DefineACLs []DefineACL `json:"defineACLs"`
}

// NewMsgDefineACLs : new message define acls
func NewMsgDefineACLs(defineACLs []DefineACL) MsgDefineACLs {
	return MsgDefineACLs{defineACLs}
}

var _ sdk.Msg = MsgDefineACLs{}

// Type : implements msg
func (msg MsgDefineACLs) Type() string { return "bank" }

// ValidateBasic : implements msg
func (msg MsgDefineACLs) ValidateBasic() sdk.Error {
	if len(msg.DefineACLs) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}
	for _, in := range msg.DefineACLs {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgDefineACLs) GetSignBytes() []byte {
	var defineACLs []json.RawMessage
	for _, defineACL := range msg.DefineACLs {
		defineACLs = append(defineACLs, defineACL.GetSignBytes())
	}
	
	b, err := msgCdc.MarshalJSON(struct {
		DefineACLs []json.RawMessage `json:"defineACLs"`
	}{
		DefineACLs: defineACLs,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgDefineACLs) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.DefineACLs))
	for i, in := range msg.DefineACLs {
		addrs[i] = in.From
	}
	return addrs
}

// BuildMsgDefineACLs : build define acls message
func BuildMsgDefineACLs(from sdk.AccAddress, to sdk.AccAddress, aclAccount sdk.ACLAccount, msgs []DefineACL) []DefineACL {
	defineACL := NewDefineACL(from, to, aclAccount)
	msgs = append(msgs, defineACL)
	return msgs
}

// BuildMsgDefineACLWithACLs : build define acls message
func BuildMsgDefineACLWithACLs(msgs []DefineACL) sdk.Msg {
	return NewMsgDefineACLs(msgs)
}

// BuildMsgDefineACL : build define acls message
func BuildMsgDefineACL(from sdk.AccAddress, to sdk.AccAddress, aclAccount sdk.ACLAccount) sdk.Msg {
	defineACL := NewDefineACL(from, to, aclAccount)
	return NewMsgDefineACLs([]DefineACL{defineACL})
}

// #####ACL

// #####Comdex
