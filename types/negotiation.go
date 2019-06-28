package types

import (
	"encoding/hex"
	
	"github.com/tendermint/tendermint/libs/common"
)

// NegotiationID :
type NegotiationID = common.HexBytes

// Signature :
type Signature = []byte

// Negotiation : interface
type Negotiation interface {
	GetNegotiationID() NegotiationID
	SetNegotiationID(NegotiationID) error
	
	GetBuyerAddress() AccAddress
	SetBuyerAddress(AccAddress) error
	
	GetSellerAddress() AccAddress
	SetSellerAddress(AccAddress) error
	
	GetPegHash() PegHash
	SetPegHash(PegHash) error
	
	GetBid() int64
	SetBid(int64) error
	
	GetTime() int64
	SetTime(int64) error
	
	GetBuyerSignature() Signature
	SetBuyerSignature(Signature) error
	
	GetSellerSignature() Signature
	SetSellerSignature(Signature) error
	
	GetBuyerBlockHeight() int64
	SetBuyerBlockHeight(int64) error
	
	GetSellerBlockHeight() int64
	SetSellerBlockHeight(int64) error
}

// BaseNegotiation : base negotiation type
type BaseNegotiation struct {
	NegotiationID     NegotiationID `json:"negotiationID" `
	BuyerAddress      AccAddress    `json:"buyerAddress" `
	SellerAddress     AccAddress    `json:"sellerAddress" `
	PegHash           PegHash       `json:"pegHash"`
	Bid               int64         `json:"bid"`
	Time              int64         `json:"time"`
	BuyerSignature    Signature     `json:"buyerSignature"`
	SellerSignature   Signature     `json:"sellerSignature"`
	BuyerBlockHeight  int64         `json:"buyerBlockHeight"`
	SellerBlockHeight int64         `json:"sellerBlockHeight"`
}

var _ Negotiation = (*BaseNegotiation)(nil)

// ProtoBaseNegotiation : converts concrete to Negotiation interface
func ProtoBaseNegotiation() Negotiation {
	return &BaseNegotiation{}
}

// GetNegotiationID : getter
func (baseNegotiation BaseNegotiation) GetNegotiationID() NegotiationID {
	return baseNegotiation.NegotiationID
}

// SetNegotiationID : setter
func (baseNegotiation *BaseNegotiation) SetNegotiationID(negotiationID NegotiationID) error {
	baseNegotiation.NegotiationID = negotiationID
	return nil
}

// GetBuyerAddress : getter
func (baseNegotiation BaseNegotiation) GetBuyerAddress() AccAddress {
	return baseNegotiation.BuyerAddress
}

// SetBuyerAddress : setter
func (baseNegotiation *BaseNegotiation) SetBuyerAddress(buyerAddress AccAddress) error {
	baseNegotiation.BuyerAddress = buyerAddress
	return nil
}

// GetSellerAddress : getter
func (baseNegotiation BaseNegotiation) GetSellerAddress() AccAddress {
	return baseNegotiation.SellerAddress
}

// SetSellerAddress : setter
func (baseNegotiation *BaseNegotiation) SetSellerAddress(sellerAddress AccAddress) error {
	baseNegotiation.SellerAddress = sellerAddress
	return nil
}

// GetPegHash : getter
func (baseNegotiation BaseNegotiation) GetPegHash() PegHash { return baseNegotiation.PegHash }

// SetPegHash : setter
func (baseNegotiation *BaseNegotiation) SetPegHash(pegHash PegHash) error {
	baseNegotiation.PegHash = pegHash
	return nil
}

// GetBid : getter
func (baseNegotiation BaseNegotiation) GetBid() int64 { return baseNegotiation.Bid }

// SetBid : setter
func (baseNegotiation *BaseNegotiation) SetBid(bid int64) error {
	baseNegotiation.Bid = bid
	return nil
}

// GetTime : getter
func (baseNegotiation BaseNegotiation) GetTime() int64 { return baseNegotiation.Time }

// SetTime : setter
func (baseNegotiation *BaseNegotiation) SetTime(time int64) error {
	baseNegotiation.Time = time
	return nil
}

// GetBuyerSignature : getter
func (baseNegotiation BaseNegotiation) GetBuyerSignature() Signature {
	return baseNegotiation.BuyerSignature
}

// SetBuyerSignature : setter
func (baseNegotiation *BaseNegotiation) SetBuyerSignature(buyerSignature Signature) error {
	baseNegotiation.BuyerSignature = buyerSignature
	return nil
}

// GetSellerSignature : getter
func (baseNegotiation BaseNegotiation) GetSellerSignature() Signature {
	return baseNegotiation.SellerSignature
}

// SetSellerSignature : setter
func (baseNegotiation *BaseNegotiation) SetSellerSignature(sellerSignature Signature) error {
	baseNegotiation.SellerSignature = sellerSignature
	return nil
}

// GetBuyerBlockHeight get buyer signed blockheight
func (baseNegotiation BaseNegotiation) GetBuyerBlockHeight() int64 {
	return baseNegotiation.BuyerBlockHeight
}

// SetBuyerBlockHeight set buyer block height
func (baseNegotiation *BaseNegotiation) SetBuyerBlockHeight(height int64) error {
	baseNegotiation.BuyerBlockHeight = height
	return nil
}

// GetSellerBlockHeight get seller signed blockheight
func (baseNegotiation BaseNegotiation) GetSellerBlockHeight() int64 {
	return baseNegotiation.SellerBlockHeight
}

// SetSellerBlockHeight set seller block height
func (baseNegotiation *BaseNegotiation) SetSellerBlockHeight(height int64) error {
	baseNegotiation.SellerBlockHeight = height
	return nil
}

// NegotiationDecoder : decoder function for Negotiation
type NegotiationDecoder func(negotiationBytes []byte) (Negotiation, error)

// GenerateNegotiationIDBytes : creates negotiation ID bytes with buyer, seller and pegHash A
func GenerateNegotiationIDBytes(buyerAddress AccAddress, sellerAddress AccAddress, pegHash PegHash) []byte {
	return append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...)
}

// GetNegotiationIDHex : convert NegotiationID string to NegotiationID hex
func GetNegotiationIDHex(negotiationIDStr string) (negotiationID NegotiationID, err error) {
	bz, err := hex.DecodeString(negotiationIDStr)
	if err != nil {
		return nil, err
	}
	return NegotiationID(bz), nil
}

// GetNegotiationIDFromHex : convert hex to string
func GetNegotiationIDFromHex(negotiationIDStr []byte) NegotiationID {
	bz := hex.EncodeToString(negotiationIDStr)
	return NegotiationID([]byte(bz))
}
