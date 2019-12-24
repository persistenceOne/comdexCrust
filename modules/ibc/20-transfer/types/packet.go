package types

import (
	"encoding/json"
	"fmt"

	"github.com/commitHub/commitBlockchain/modules/bank"
	"github.com/commitHub/commitBlockchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PacketData defines a struct for the packet payload
type PacketData struct {
	Amount   sdk.Coins      `json:"amount" yaml:"amount"`     // the tokens to be transferred
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`     // the sender address
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"` // the recipient address on the destination chain
	Source   bool           `json:"source" yaml:"source"`     // indicates if the sending chain is the source chain of the tokens to be transferred
}

func (pd PacketData) MarshalAmino() ([]byte, error) {
	return ModuleCdc.MarshalBinaryBare(pd)
}

func (pd *PacketData) UnmarshalAmino(bz []byte) (err error) {
	return ModuleCdc.UnmarshalBinaryBare(bz, pd)
}

func (pd PacketData) Marshal() []byte {
	return ModuleCdc.MustMarshalBinaryBare(pd)
}

type PacketDataAlias PacketData

// MarshalJSON implements the json.Marshaler interface.
func (pd PacketData) MarshalJSON() ([]byte, error) {
	return json.Marshal((PacketDataAlias)(pd))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (pd *PacketData) UnmarshalJSON(bz []byte) (err error) {
	return json.Unmarshal(bz, (*PacketDataAlias)(pd))
}

func (pd PacketData) String() string {
	return fmt.Sprintf(`PacketData:
	Amount:               %s
	Sender:               %s
	Receiver:             %s
	Source:               %v`,
		pd.Amount.String(),
		pd.Sender,
		pd.Receiver,
		pd.Source,
	)
}

// ValidateBasic performs a basic check of the packet fields
func (pd PacketData) ValidateBasic() sdk.Error {
	if !pd.Amount.IsValid() {
		return sdk.ErrInvalidCoins("transfer amount is invalid")
	}
	if !pd.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("transfer amount must be positive")
	}
	if pd.Sender.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if pd.Receiver.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}
	return nil
}

// **** Issue Asset Packet Data **** //
type IssueAsset struct {
	IssuerAddress sdk.AccAddress     `json:"issuerAddress"`
	ToAddress     sdk.AccAddress     `json:"toAddress"`
	AssetPeg      types.BaseAssetPeg `json:"assetPeg"`
}

func NewIssueAsset(issueAsset bank.IssueAsset) IssueAsset {
	return IssueAsset{
		IssuerAddress: issueAsset.IssuerAddress,
		ToAddress:     issueAsset.ToAddress,
		AssetPeg:      types.ToBaseAssetPeg(issueAsset.AssetPeg),
	}
}

type IssueAssetPacketData struct {
	IssueAsset IssueAsset `json:"issueAsset" yaml:"issueAsset"` // the asset to be transferred
}

func (pd IssueAssetPacketData) MarshalAmino() ([]byte, error) {
	return ModuleCdc.MarshalBinaryBare(pd)
}

func (pd *IssueAssetPacketData) UnmarshalAmino(bz []byte) (err error) {
	return ModuleCdc.UnmarshalBinaryBare(bz, pd)
}

func (pd IssueAssetPacketData) Marshal() []byte {
	return ModuleCdc.MustMarshalBinaryBare(pd)
}

type IssueAssetPacketDataAlias IssueAssetPacketData

// MarshalJSON implements the json.Marshaler interface.
func (pd IssueAssetPacketData) MarshalJSON() ([]byte, error) {
	return json.Marshal((IssueAssetPacketDataAlias)(pd))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (pd *IssueAssetPacketData) UnmarshalJSON(bz []byte) (err error) {
	return json.Unmarshal(bz, (*IssueAssetPacketDataAlias)(pd))
}

func (pd IssueAssetPacketData) String() string {
	return fmt.Sprintf(`PacketData:
	IssueAsset:           {
		IssuerAddress: %s
		ToAddress: %s
		AssetPeg: {
			PegHash: %s
			DocumentHash: %s
			AssetType: %s
			AssetQuantity: %d
			AssetPrice: %d
			QuantityUnit: %s
			OwnerAddress: %s
			Locked: %v
			Moderated: %v
			TakerAddress: %s
		}
	}`,
		pd.IssueAsset.IssuerAddress,
		pd.IssueAsset.ToAddress,
		pd.IssueAsset.AssetPeg.GetPegHash(),
		pd.IssueAsset.AssetPeg.GetDocumentHash(),
		pd.IssueAsset.AssetPeg.GetAssetType(),
		pd.IssueAsset.AssetPeg.GetAssetQuantity(),
		pd.IssueAsset.AssetPeg.GetAssetPrice(),
		pd.IssueAsset.AssetPeg.GetQuantityUnit(),
		pd.IssueAsset.AssetPeg.GetOwnerAddress(),
		pd.IssueAsset.AssetPeg.GetLocked(),
		pd.IssueAsset.AssetPeg.GetModerated(),
		pd.IssueAsset.AssetPeg.GetTakerAddress(),
	)
}

// ValidateBasic performs a basic check of the packet fields
func (in IssueAssetPacketData) ValidateBasic() sdk.Error {
	// if err := pd.IssueAsset.ValidateBasic(); err != nil {
	// 	return err
	// }
	// return nil
	if len(in.IssueAsset.IssuerAddress) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("invalid Issuer address %s", in.IssueAsset.IssuerAddress.String()))
	} else if len(in.IssueAsset.ToAddress) == 0 {
		return sdk.ErrInvalidAddress(fmt.Sprintf("invalid To address %s", in.IssueAsset.ToAddress.String()))
	} else if in.IssueAsset.AssetPeg.GetAssetPrice() < 0 {
		return types.ErrNegativeAmount(DefaultCodespace, "Asset price should be grater than 0.")
	} else if in.IssueAsset.AssetPeg.GetAssetQuantity() < 0 {
		return types.ErrNegativeAmount(DefaultCodespace, "Asset quantity should be grater than 0.")
	} else if in.IssueAsset.AssetPeg.GetAssetType() == "" {
		return sdk.ErrUnknownRequest("asset type should not be empty")
	} else if in.IssueAsset.AssetPeg.GetDocumentHash() == "" {
		return sdk.ErrUnknownRequest("DocumentHash should not be empty")
	}
	return nil
}

// #### Issue Asset Packet Data #### //
