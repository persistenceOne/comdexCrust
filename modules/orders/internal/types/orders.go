package types

import (
	"github.com/commitHub/commitBlockchain/modules/negotiation"
	"github.com/commitHub/commitBlockchain/types"
)

type Order interface {
	GetNegotiationID() negotiation.NegotiationID
	SetNegotiationID(negotiation.NegotiationID)

	GetAssetPegWallet() types.AssetPegWallet
	SetAssetPegWallet(types.AssetPegWallet)

	GetFiatPegWallet() types.FiatPegWallet
	SetFiatPegWallet(types.FiatPegWallet)

	GetFiatProofHash() string
	SetFiatProofHash(string)

	GetAWBProofHash() string
	SetAWBProofHash(string)
}

type BaseOrder struct {
	NegotiationID  negotiation.NegotiationID `json:"negotiationID"`
	FiatPegWallet  types.FiatPegWallet       `json:"fiatPegWallet"`
	AssetPegWallet types.AssetPegWallet      `json:"assetPegWallet"`
	FiatProofHash  string                    `json:"fiatProofHash"`
	AWBProofHash   string                    `json:"awbProofHash"`
}

var _ Order = (*BaseOrder)(nil)

func ProtoBaseOrder() Order {
	return &BaseOrder{}
}

func (baseOrder BaseOrder) GetNegotiationID() negotiation.NegotiationID {
	return baseOrder.NegotiationID
}

func (baseOrder *BaseOrder) SetNegotiationID(negotiationID negotiation.NegotiationID) {
	baseOrder.NegotiationID = negotiationID
}

func (baseOrder BaseOrder) GetAssetPegWallet() types.AssetPegWallet {
	return baseOrder.AssetPegWallet
}

func (baseOrder *BaseOrder) SetAssetPegWallet(assetPegWallet types.AssetPegWallet) {
	baseOrder.AssetPegWallet = assetPegWallet
}

func (baseOrder BaseOrder) GetFiatPegWallet() types.FiatPegWallet {
	return baseOrder.FiatPegWallet
}

func (baseOrder *BaseOrder) SetFiatPegWallet(fiatPegWallet types.FiatPegWallet) {
	baseOrder.FiatPegWallet = fiatPegWallet
}

func (baseOrder BaseOrder) GetFiatProofHash() string {
	return baseOrder.FiatProofHash
}

func (baseOrder *BaseOrder) SetFiatProofHash(fiatProofHash string) {
	baseOrder.FiatProofHash = fiatProofHash
}

func (baseOrder BaseOrder) GetAWBProofHash() string {
	return baseOrder.AWBProofHash
}

func (baseOrder *BaseOrder) SetAWBProofHash(awbProofHash string) {
	baseOrder.AWBProofHash = awbProofHash
}

type OrderDecoder func(orderBytes []byte) (Order, error)
