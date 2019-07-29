package assetFactory

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"

	sdk "github.com/commitHub/commitBlockchain/types"
)

//*****Comdex

//*****IssueAsset

//IssueAsset - transaction input
type IssueAsset struct {
	IssuerAddress sdk.AccAddress `json:"issuerAddress"`
	ToAddress     sdk.AccAddress `json:"toAddress"`
	AssetPeg      sdk.AssetPeg   `json:"assetPeg"`
}

//NewIssueAsset : initializer
func NewIssueAsset(issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, assetPeg sdk.AssetPeg) IssueAsset {
	return IssueAsset{issuerAddress, toAddress, assetPeg}
}

//GetSignBytes : get bytes to sign
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

//#####IssueAsset

//*****MsgFactoryIssueAssets

//MsgFactoryIssueAssets : high level issuance of assets module
type MsgFactoryIssueAssets struct {
	IssueAssets []IssueAsset `json:"issueAssets"`
}

//NewMsgFactoryIssueAssets : initilizer
func NewMsgFactoryIssueAssets(issueAssets []IssueAsset) MsgFactoryIssueAssets {
	return MsgFactoryIssueAssets{issueAssets}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgFactoryIssueAssets{}

//Type : implements msg
func (msg MsgFactoryIssueAssets) Type() string { return "assetFactory" }

//ValidateBasic : implements msg
func (msg MsgFactoryIssueAssets) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.IssueAssets {
		if len(in.IssuerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.IssuerAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		}
		_, err = govalidator.ValidateStruct(in.AssetPeg)
		if err != nil {
			return sdk.ErrInvalidAddress(err.Error())
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgFactoryIssueAssets) GetSignBytes() []byte {
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
func (msg MsgFactoryIssueAssets) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.IssueAssets))
	for i, in := range msg.IssueAssets {
		addrs[i] = in.IssuerAddress
	}
	return addrs
}

//BuildIssueAssetMsg : butild the issueAssetTx
func BuildIssueAssetMsg(from sdk.AccAddress, to sdk.AccAddress, assetPeg sdk.AssetPeg) sdk.Msg {
	issueAsset := NewIssueAsset(from, to, assetPeg)
	msg := NewMsgFactoryIssueAssets([]IssueAsset{issueAsset})
	return msg
}

//##### Implement sdk.Msg

//#####MsgFactoryIssueAssets

//*****RedeemAsset

//RedeemAsset - transaction input
type RedeemAsset struct {
	RelayerAddress sdk.AccAddress `json:"relayerAddress"`
	OwnerAddress   sdk.AccAddress `json:"ownerAddress"`
	ToAddress      sdk.AccAddress `json:"toAddress"`
	PegHash        sdk.PegHash    `json:"pegHash"`
}

//NewRedeemAsset : initializer
func NewRedeemAsset(relayerAddress sdk.AccAddress, ownerAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash) RedeemAsset {
	return RedeemAsset{relayerAddress, ownerAddress, toAddress, pegHash}
}

//GetSignBytes : get bytes to sign
func (in RedeemAsset) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		RelayerAddress string      `json:"relayerAddress"`
		OwnerAddress   string      `json:"ownerAddress"`
		ToAddress      string      `json:"toAddress"`
		PegHash        sdk.PegHash `json:"assetPeg"`
	}{
		RelayerAddress: in.RelayerAddress.String(),
		OwnerAddress:   in.OwnerAddress.String(),
		ToAddress:      in.ToAddress.String(),
		PegHash:        in.PegHash,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//#####RedeemAsset

//*****MsgFactoryRedeemAssets

//MsgFactoryRedeemAssets : high level redeem of assets module
type MsgFactoryRedeemAssets struct {
	RedeemAssets []RedeemAsset `json:"redeemAssets"`
}

//NewMsgFactoryRedeemAssets : initilizer
func NewMsgFactoryRedeemAssets(redeemAssets []RedeemAsset) MsgFactoryRedeemAssets {
	return MsgFactoryRedeemAssets{redeemAssets}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgFactoryRedeemAssets{}

//Type : implements msg
func (msg MsgFactoryRedeemAssets) Type() string { return "assetFactory" }

//ValidateBasic : implements msg
func (msg MsgFactoryRedeemAssets) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.RedeemAssets {
		if len(in.RelayerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RelayerAddress.String())
		} else if len(in.OwnerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.OwnerAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgFactoryRedeemAssets) GetSignBytes() []byte {
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
func (msg MsgFactoryRedeemAssets) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.RedeemAssets))
	for i, in := range msg.RedeemAssets {
		addrs[i] = in.RelayerAddress
	}
	return addrs
}

//BuildRedeemAssetMsg : build the redeemAssetTx
func BuildRedeemAssetMsg(relayerAddress sdk.AccAddress, ownerAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash) sdk.Msg {
	redeemAsset := NewRedeemAsset(relayerAddress, ownerAddress, toAddress, pegHash)
	msg := NewMsgFactoryRedeemAssets([]RedeemAsset{redeemAsset})
	return msg
}

//##### Implement sdk.Msg

//#####MsgFactoryRedeemAssets

//*****SendAsset

//SendAsset - transaction input
type SendAsset struct {
	RelayerAddress sdk.AccAddress `json:"relayerAddress"`
	FromAddress    sdk.AccAddress `json:"fromAddress"`
	ToAddress      sdk.AccAddress `json:"toAddress"`
	PegHash        sdk.PegHash    `json:"pegHash"`
}

//NewSendAsset : initializer
func NewSendAsset(relayerAddress sdk.AccAddress, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash) SendAsset {
	return SendAsset{relayerAddress, fromAddress, toAddress, pegHash}
}

//GetSignBytes : get bytes to sign
func (in SendAsset) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(struct {
		RelayerAddress string      `json:"relayerAddress"`
		FromAddress    string      `json:"fromAddress"`
		ToAddress      string      `json:"toAddress"`
		PegHash        sdk.PegHash `json:"pegHash"`
	}{
		RelayerAddress: in.RelayerAddress.String(),
		FromAddress:    in.FromAddress.String(),
		ToAddress:      in.ToAddress.String(),
		PegHash:        in.PegHash,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

//#####SendAsset

//*****MsgFactorySendAssets

//MsgFactorySendAssets : high level issuance of assets module
type MsgFactorySendAssets struct {
	SendAssets []SendAsset `json:"sendAssets"`
}

//NewMsgFactorySendAssets : initilizer
func NewMsgFactorySendAssets(sendAssets []SendAsset) MsgFactorySendAssets {
	return MsgFactorySendAssets{sendAssets}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgFactorySendAssets{}

//Type : implements msg
func (msg MsgFactorySendAssets) Type() string { return "assetFactory" }

//ValidateBasic : implements msg
func (msg MsgFactorySendAssets) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.SendAssets {
		if len(in.RelayerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RelayerAddress.String())
		} else if len(in.FromAddress) == 0 {
			return sdk.ErrInvalidAddress(in.FromAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgFactorySendAssets) GetSignBytes() []byte {
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
func (msg MsgFactorySendAssets) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SendAssets))
	for i, in := range msg.SendAssets {
		addrs[i] = in.RelayerAddress
	}
	return addrs
}

//BuildSendAssetMsg : build the sendAssetTx
func BuildSendAssetMsg(relayerAddress sdk.AccAddress, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash) sdk.Msg {
	sendAsset := NewSendAsset(relayerAddress, fromAddress, toAddress, pegHash)
	msg := NewMsgFactorySendAssets([]SendAsset{sendAsset})
	return msg
}

//##### Implement sdk.Msg

//#####MsgFactorySendAssets

//*****MsgFactoryExecuteAssets

//MsgFactoryExecuteAssets : high level issuance of assets module
type MsgFactoryExecuteAssets struct {
	SendAssets []SendAsset `json:"sendAssets"`
}

//NewMsgFactoryExecuteAssets : initilizer
func NewMsgFactoryExecuteAssets(sendAssets []SendAsset) MsgFactoryExecuteAssets {
	return MsgFactoryExecuteAssets{sendAssets}
}

//***** Implementing sdk.Msg

var _ sdk.Msg = MsgFactoryExecuteAssets{}

//Type : implements msg
func (msg MsgFactoryExecuteAssets) Type() string { return "assetFactory" }

//ValidateBasic : implements msg
func (msg MsgFactoryExecuteAssets) ValidateBasic() sdk.Error {
	_, err := govalidator.ValidateStruct(msg)
	if err != nil {
		return sdk.ErrInsufficientFunds(err.Error())
	}
	for _, in := range msg.SendAssets {
		if len(in.RelayerAddress) == 0 {
			return sdk.ErrInvalidAddress(in.RelayerAddress.String())
		} else if len(in.FromAddress) == 0 {
			return sdk.ErrInvalidAddress(in.FromAddress.String())
		} else if len(in.ToAddress) == 0 {
			return sdk.ErrInvalidAddress(in.ToAddress.String())
		} else if len(in.PegHash) == 0 {
			return sdk.ErrUnknownRequest("PegHash is Empty")
		}
	}
	return nil
}

//GetSignBytes : implements msg
func (msg MsgFactoryExecuteAssets) GetSignBytes() []byte {
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
func (msg MsgFactoryExecuteAssets) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.SendAssets))
	for i, in := range msg.SendAssets {
		addrs[i] = in.RelayerAddress
	}
	return addrs
}

//BuildExecuteAssetMsg : build the executeAssetTx
func BuildExecuteAssetMsg(relayerAddress sdk.AccAddress, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash) sdk.Msg {
	sendAsset := NewSendAsset(relayerAddress, fromAddress, toAddress, pegHash)
	msg := NewMsgFactoryExecuteAssets([]SendAsset{sendAsset})
	return msg
}

//##### Implement sdk.Msg

//#####MsgFactoryExecuteAssets

//#####Comdex
