package types

type Order interface {
	GetNegotiationID() NegotiationID
	SetNegotiationID(NegotiationID)

	GetAssetPegWallet() AssetPegWallet
	SetAssetPegWallet(AssetPegWallet)

	GetFiatPegWallet() FiatPegWallet
	SetFiatPegWallet(FiatPegWallet)

	GetFiatProofHash() string
	SetFiatProofHash(string)

	GetAWBProofHash() string
	SetAWBProofHash(string)
}

type BaseOrder struct {
	NegotiationID  NegotiationID  `json:"negotiationID"`
	FiatPegWallet  FiatPegWallet  `json:"fiatPegWallet"`
	AssetPegWallet AssetPegWallet `json:"assetPegWallet"`
	FiatProofHash  string         `json:"fiatProofHash"`
	AWBProofHash   string         `json:"awbProofHash"`
}

var _ Order = (*BaseOrder)(nil)

func ProtoBaseOrder() Order {
	return &BaseOrder{}
}

func (baseOrder BaseOrder) GetNegotiationID() NegotiationID {
	return baseOrder.NegotiationID
}

func (baseOrder *BaseOrder) SetNegotiationID(negotiationID NegotiationID) {
	baseOrder.NegotiationID = negotiationID
}

func (baseOrder BaseOrder) GetAssetPegWallet() AssetPegWallet {
	return baseOrder.AssetPegWallet
}

func (baseOrder *BaseOrder) SetAssetPegWallet(assetPegWallet AssetPegWallet) {
	baseOrder.AssetPegWallet = assetPegWallet
}

func (baseOrder BaseOrder) GetFiatPegWallet() FiatPegWallet {
	return baseOrder.FiatPegWallet
}

func (baseOrder *BaseOrder) SetFiatPegWallet(fiatPegWallet FiatPegWallet) {
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
