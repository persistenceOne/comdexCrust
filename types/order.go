package types

// Order : order interface
type Order interface {
	GetNegotiationID() NegotiationID
	SetNegotiationID(NegotiationID) error

	GetAssetPegWallet() AssetPegWallet
	SetAssetPegWallet(AssetPegWallet) error

	GetFiatPegWallet() FiatPegWallet
	SetFiatPegWallet(FiatPegWallet) error

	GetFiatProofHash() string
	SetFiatProofHash(string) error

	GetAWBProofHash() string
	SetAWBProofHash(string) error
}

// BaseOrder : base order type
type BaseOrder struct {
	NegotiationID  NegotiationID  `json:"negotiationID"`
	FiatPegWallet  FiatPegWallet  `json:"fiatPegWallet"`
	AssetPegWallet AssetPegWallet `json:"assetPegWallet"`
	FiatProofHash  string         `json:"fiatProofHash"`
	AWBProofHash   string         `json:"awbProofHash"`
}

var _ Order = (*BaseOrder)(nil)

// ProtoBaseOrder : converts concrete to order in
func ProtoBaseOrder() Order {
	return &BaseOrder{}
}

//GetNegotiationID : getter
func (baseOrder BaseOrder) GetNegotiationID() NegotiationID {
	return baseOrder.NegotiationID
}

//SetNegotiationID : setter
func (baseOrder *BaseOrder) SetNegotiationID(negotiationID NegotiationID) error {
	baseOrder.NegotiationID = negotiationID
	return nil
}

//GetAssetPegWallet : getter
func (baseOrder BaseOrder) GetAssetPegWallet() AssetPegWallet {
	return baseOrder.AssetPegWallet
}

//SetAssetPegWallet : setter
func (baseOrder *BaseOrder) SetAssetPegWallet(assetPegWallet AssetPegWallet) error {
	baseOrder.AssetPegWallet = assetPegWallet
	return nil
}

//GetFiatPegWallet : getter
func (baseOrder BaseOrder) GetFiatPegWallet() FiatPegWallet {
	return baseOrder.FiatPegWallet
}

//SetFiatPegWallet : setter
func (baseOrder *BaseOrder) SetFiatPegWallet(fiatPegWallet FiatPegWallet) error {
	baseOrder.FiatPegWallet = fiatPegWallet
	return nil
}

//GetFiatProofHash : getter
func (baseOrder BaseOrder) GetFiatProofHash() string {
	return baseOrder.FiatProofHash
}

//SetFiatProofHash : setter
func (baseOrder *BaseOrder) SetFiatProofHash(fiatProofHash string) error {
	baseOrder.FiatProofHash = fiatProofHash
	return nil
}

//GetAWBProofHash : getter
func (baseOrder BaseOrder) GetAWBProofHash() string {
	return baseOrder.AWBProofHash
}

//SetAWBProofHash : setter
func (baseOrder *BaseOrder) SetAWBProofHash(awbProofHash string) error {
	baseOrder.AWBProofHash = awbProofHash
	return nil
}

//OrderDecoder : decoder function for order
type OrderDecoder func(orderBytes []byte) (Order, error)
