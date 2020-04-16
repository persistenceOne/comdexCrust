package types

import (
	"encoding/json"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/comdexCrust/types"
)

type IssueAsset struct {
	IssuerAddress cTypes.AccAddress `json:"issuerAddress"`
	ToAddress     cTypes.AccAddress `json:"toAddress"`
	AssetPeg      types.AssetPeg    `json:"assetPeg"`
}

// NewIssueAsset : returns issueAsset
func NewIssueAsset(issuerAddress cTypes.AccAddress, toAddress cTypes.AccAddress, assetPeg types.AssetPeg) IssueAsset {
	return IssueAsset{
		IssuerAddress: issuerAddress,
		ToAddress:     toAddress,
		AssetPeg:      assetPeg,
	}
}

// GetSignBytes : get bytes to sign
func (in IssueAsset) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(struct {
		IssuerAddress string         `json:"issuerAddress"`
		ToAddress     string         `json:"toAddress"`
		AssetPeg      types.AssetPeg `json:"assetPeg"`
	}{
		IssuerAddress: in.IssuerAddress.String(),
		ToAddress:     in.ToAddress.String(),
		AssetPeg:      in.AssetPeg,
	})
	if err != nil {
		return nil
	}
	return bz
}

func (in IssueAsset) ValidateBasic() cTypes.Error {
	if len(in.IssuerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.IssuerAddress.String())
	} else if len(in.ToAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.ToAddress.String())
	} else if (in.AssetPeg.GetAssetPrice()) < 0 {
		return ErrInvalidAmount(DefaultCodeSpace, "Asset price should be greater than 0")
	} else if (in.AssetPeg.GetAssetQuantity()) < 0 {
		return ErrInvalidAmount(DefaultCodeSpace, "Asset quantity should not be 0")
	} else if (in.AssetPeg.GetAssetType()) == "" {
		return ErrInvalidString(DefaultCodeSpace, "AssetType should not be empty string")
	} else if (in.AssetPeg.GetDocumentHash()) == "" {
		return ErrInvalidString(DefaultCodeSpace, "DocumentHash should not be empty string")
	} else if (in.AssetPeg.GetQuantityUnit()) == "" {
		return ErrInvalidString(DefaultCodeSpace, "QuantityUnit should not be empty string")
	} else if len(in.AssetPeg.GetPegHash()) == 0 {
		return ErrInvalidString(DefaultCodeSpace, "PegHash should not be empty")
	}
	return nil
}

type MsgFactoryIssueAssets struct {
	IssueAssets []IssueAsset
}

func NewMsgFactoryIssueAssets(issueAssets []IssueAsset) MsgFactoryIssueAssets {
	return MsgFactoryIssueAssets{IssueAssets: issueAssets}
}

var _ cTypes.Msg = MsgFactoryIssueAssets{}

func (msg MsgFactoryIssueAssets) Type() string { return "assetFactory" }

func (msg MsgFactoryIssueAssets) Route() string { return RouterKey }

func (msg MsgFactoryIssueAssets) ValidateBasic() cTypes.Error {
	if len(msg.IssueAssets) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.IssueAssets {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

func (msg MsgFactoryIssueAssets) GetSignBytes() []byte {
	var issueAssets []json.RawMessage
	for _, issueAsset := range msg.IssueAssets {
		issueAssets = append(issueAssets, issueAsset.GetSignBytes())
	}

	bz, err := ModuleCdc.MarshalJSON(struct {
		IssueAssets []json.RawMessage `json:"issueAsset"`
	}{
		IssueAssets: issueAssets,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MsgFactoryIssueAssets) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.IssueAssets))
	for i, in := range msg.IssueAssets {
		addrs[i] = in.IssuerAddress
	}
	return addrs
}

func BuildIssueAssetMsg(issuerAddress cTypes.AccAddress, toAddress cTypes.AccAddress, assetPeg types.AssetPeg) cTypes.Msg {
	issueAsset := NewIssueAsset(issuerAddress, toAddress, assetPeg)
	msg := NewMsgFactoryIssueAssets([]IssueAsset{issueAsset})
	return msg
}

type RedeemAsset struct {
	RelayerAddress cTypes.AccAddress `json:"relayerAddress"`
	OwnerAddress   cTypes.AccAddress `json:"ownerAddress"`
	ToAddress      cTypes.AccAddress `json:"toAddress"`
	PegHash        types.PegHash     `json:"pegHash"`
}

func NewRedeemAsset(relayerAddress cTypes.AccAddress, ownerAddress cTypes.AccAddress, toAddress cTypes.AccAddress, pegHash types.PegHash) RedeemAsset {
	return RedeemAsset{
		RelayerAddress: relayerAddress,
		OwnerAddress:   ownerAddress,
		ToAddress:      toAddress,
		PegHash:        pegHash,
	}
}

func (in RedeemAsset) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(struct {
		RelayerAddress string        `json:relayerAddress`
		OwnerAddress   string        `json:ownerAddress`
		ToAddress      string        `json:toAddress`
		PegHash        types.PegHash `json:"pegHash"`
	}{
		RelayerAddress: in.RelayerAddress.String(),
		OwnerAddress:   in.OwnerAddress.String(),
		ToAddress:      in.ToAddress.String(),
		PegHash:        in.PegHash,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

func (in RedeemAsset) ValidateBasic() cTypes.Error {
	if len(in.RelayerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.RelayerAddress.String())
	} else if len(in.OwnerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.OwnerAddress.String())
	} else if len(in.ToAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.ToAddress.String())
	} else if len(in.PegHash) == 0 {
		return ErrInvalidString(DefaultCodeSpace, "PegHash should not be empty")
	}
	return nil
}

type MsgFactoryRedeemAssets struct {
	RedeemAssets []RedeemAsset `json:"redeemAssets"`
}

func NewMsgFactoryRedeemAssets(redeemAssets []RedeemAsset) MsgFactoryRedeemAssets {
	return MsgFactoryRedeemAssets{RedeemAssets: redeemAssets}
}

var _ cTypes.Msg = MsgFactoryRedeemAssets{}

func (msg MsgFactoryRedeemAssets) Type() string { return "assetFactory" }

func (msg MsgFactoryRedeemAssets) Route() string { return RouterKey }

func (msg MsgFactoryRedeemAssets) ValidateBasic() cTypes.Error {
	if len(msg.RedeemAssets) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.RedeemAssets {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

func (msg MsgFactoryRedeemAssets) GetSignBytes() []byte {
	var redeemAssets []json.RawMessage
	for _, redeemAsset := range msg.RedeemAssets {
		redeemAssets = append(redeemAssets, redeemAsset.GetSignBytes())
	}

	bz, err := ModuleCdc.MarshalJSON(struct {
		RedeemAssets []json.RawMessage `json:"redeemAssets"`
	}{
		RedeemAssets: redeemAssets,
	})

	if err != nil {
		panic(err)
	}

	return bz
}

func (msg MsgFactoryRedeemAssets) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.RedeemAssets))

	for i, in := range msg.RedeemAssets {
		addrs[i] = in.RelayerAddress
	}

	return addrs
}

func BuildRedeemAssetMsg(relayerAddress cTypes.AccAddress, ownerAddress cTypes.AccAddress, toAddress cTypes.AccAddress, pegHash types.PegHash) cTypes.Msg {
	redeemAsset := NewRedeemAsset(relayerAddress, ownerAddress, toAddress, pegHash)
	msg := NewMsgFactoryRedeemAssets([]RedeemAsset{redeemAsset})
	return msg
}

type SendAsset struct {
	RelayerAddress cTypes.AccAddress `json:"relayerAddress"`
	FromAddress    cTypes.AccAddress `json:"fromAddress"`
	ToAddress      cTypes.AccAddress `json:"toAddress"`
	PegHash        types.PegHash     `json:"pegHash"`
}

func NewSendAsset(relayerAddress cTypes.AccAddress, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress, peghash types.PegHash) SendAsset {
	return SendAsset{
		RelayerAddress: relayerAddress,
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		PegHash:        peghash,
	}
}

func (in SendAsset) GetSignBytes() []byte {
	biz, err := ModuleCdc.MarshalJSON(struct {
		RelayerAddress string        `json:"relayerAddress"`
		FromAddress    string        `json:"fromAddress"`
		ToAddress      string        `json:"toAddress"`
		PegHash        types.PegHash `json:"pegHash"`
	}{
		RelayerAddress: in.RelayerAddress.String(),
		FromAddress:    in.FromAddress.String(),
		ToAddress:      in.ToAddress.String(),
		PegHash:        in.PegHash,
	})

	if err != nil {
		panic(err)
	}

	return biz
}

func (in SendAsset) ValidateBasic() cTypes.Error {
	if len(in.RelayerAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.RelayerAddress.String())
	} else if len(in.FromAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.FromAddress.String())
	} else if len(in.ToAddress) == 0 {
		return cTypes.ErrInvalidAddress(in.ToAddress.String())
	} else if len(in.PegHash) == 0 {
		return ErrInvalidString(DefaultCodeSpace, "PegHash should not be empty")
	}
	return nil
}

type MsgFactorySendAssets struct {
	SendAssets []SendAsset `json:"sendAssets"`
}

func NewMsgFactorySendAssets(sendAssets []SendAsset) MsgFactorySendAssets {
	return MsgFactorySendAssets{SendAssets: sendAssets}
}

var _ cTypes.Msg = MsgFactorySendAssets{}

func (msg MsgFactorySendAssets) Type() string { return "assetFactory" }

func (msg MsgFactorySendAssets) Route() string { return RouterKey }

func (msg MsgFactorySendAssets) ValidateBasic() cTypes.Error {
	if len(msg.SendAssets) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.SendAssets {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

func (msg MsgFactorySendAssets) GetSignBytes() []byte {
	var sendAssets []json.RawMessage
	for _, sendAsset := range msg.SendAssets {
		sendAssets = append(sendAssets, sendAsset.GetSignBytes())
	}
	bz, err := ModuleCdc.MarshalJSON(struct {
		SendAssets []json.RawMessage `json:"sendAssets"`
	}{
		SendAssets: sendAssets,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MsgFactorySendAssets) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.SendAssets))
	for i, in := range msg.SendAssets {
		addrs[i] = in.RelayerAddress
	}
	return addrs
}

func BuildSendAssetMsg(relayerAddress cTypes.AccAddress, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress, peghash types.PegHash) cTypes.Msg {
	sendAsset := NewSendAsset(relayerAddress, fromAddress, toAddress, peghash)
	msg := NewMsgFactorySendAssets([]SendAsset{sendAsset})
	return msg
}

type MsgFactoryExecuteAssets struct {
	SendAssets []SendAsset `json:"sendAssets"`
}

func NewMsgFactoryExecuteAssets(sendAssets []SendAsset) MsgFactoryExecuteAssets {
	return MsgFactoryExecuteAssets{SendAssets: sendAssets}
}

var _ cTypes.Msg = MsgFactoryExecuteAssets{}

func (msg MsgFactoryExecuteAssets) Type() string { return "assetFactory" }

func (msg MsgFactoryExecuteAssets) Route() string { return RouterKey }

func (msg MsgFactoryExecuteAssets) ValidateBasic() cTypes.Error {
	if len(msg.SendAssets) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.SendAssets {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

func (msg MsgFactoryExecuteAssets) GetSignBytes() []byte {
	var sendAssets []json.RawMessage
	for _, sendAsset := range msg.SendAssets {
		sendAssets = append(sendAssets, sendAsset.GetSignBytes())
	}

	bz, err := ModuleCdc.MarshalJSON(struct {
		SendAssets []json.RawMessage `json:"sendAssets"`
	}{
		SendAssets: sendAssets,
	})
	if err != nil {
		panic(err)
	}
	return bz
}

func (msg MsgFactoryExecuteAssets) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.SendAssets))
	for i, in := range msg.SendAssets {
		addrs[i] = in.RelayerAddress
	}
	return addrs
}

func BuildExecuteAssetMsg(relayerAddress cTypes.AccAddress, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress, pegHash types.PegHash) cTypes.Msg {
	executeAsset := NewSendAsset(relayerAddress, fromAddress, toAddress, pegHash)
	msg := NewMsgFactoryExecuteAssets([]SendAsset{executeAsset})
	return msg
}
