package types

import (
	"encoding/hex"
	"fmt"
	
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/common"
	
	"github.com/commitHub/commitBlockchain/types"
)

type NegotiationID = common.HexBytes

type Signature []byte

type Negotiation interface {
	GetNegotiationID() NegotiationID
	SetNegotiationID(NegotiationID) error
	
	GetBuyerAddress() cTypes.AccAddress
	SetBuyerAddress(cTypes.AccAddress) error
	
	GetSellerAddress() cTypes.AccAddress
	SetSellerAddress(cTypes.AccAddress) error
	
	GetPegHash() types.PegHash
	SetPegHash(types.PegHash) error
	
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
	
	GetBuyerContractHash() string
	SetBuyerContractHash(string) error
	
	GetSellerContractHash() string
	SetSellerContractHash(string) error
}

var _ Negotiation = (*BaseNegotiation)(nil)

// BaseNegotiation : base negotiation type
type BaseNegotiation struct {
	NegotiationID      NegotiationID     `json:"negotiationID" `
	BuyerAddress       cTypes.AccAddress `json:"buyerAddress" `
	SellerAddress      cTypes.AccAddress `json:"sellerAddress" `
	PegHash            types.PegHash     `json:"pegHash"`
	Bid                int64             `json:"bid"`
	Time               int64             `json:"time"`
	BuyerSignature     Signature         `json:"buyerSignature"`
	SellerSignature    Signature         `json:"sellerSignature"`
	BuyerBlockHeight   int64             `json:"buyerBlockHeight"`
	SellerBlockHeight  int64             `json:"sellerBlockHeight"`
	BuyerContractHash  string            `json:"buyerContractHash"`
	SellerContractHash string            `json:"sellerContractHash"`
}

func (negotiation BaseNegotiation) String() string {
	return fmt.Sprintf(`Negotiation:
NegotiationID: %s,
BuyerAddress: %s,
SellerAddress:%s,
PegHash: %s,
Bid: %d,
Time: %d,
BuyerSignature: %s,
SellerSignature: %s,
BuyerBlockHeight: %d,
SellerBlockHeight: %d,
`, negotiation.NegotiationID.String(), negotiation.BuyerAddress.String(), negotiation.SellerAddress.String(), negotiation.PegHash.String(), negotiation.Bid, negotiation.Time,
		negotiation.BuyerSignature.String(), negotiation.SellerSignature.String(), negotiation.BuyerBlockHeight, negotiation.SellerBlockHeight)
}

// ProtoBaseNegotiation : converts concrete to Negotiation interface
func ProtoBaseNegotiation() Negotiation {
	return &BaseNegotiation{}
}

func NewBaseNegotiation(id NegotiationID, buyerAddress, sellerAddress cTypes.AccAddress, bid, time int64, buyerSignature, sellerSignature Signature, buyerBlockHeight, sellerBlockHeight int64, buyerContractHash, sellerContractHash string) Negotiation {
	return &BaseNegotiation{
		NegotiationID:      id,
		BuyerAddress:       buyerAddress,
		SellerAddress:      sellerAddress,
		Bid:                bid,
		Time:               time,
		BuyerSignature:     buyerSignature,
		SellerSignature:    sellerSignature,
		BuyerBlockHeight:   buyerBlockHeight,
		SellerBlockHeight:  sellerBlockHeight,
		BuyerContractHash:  buyerContractHash,
		SellerContractHash: sellerContractHash,
	}
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
func (baseNegotiation BaseNegotiation) GetBuyerAddress() cTypes.AccAddress {
	return baseNegotiation.BuyerAddress
}

// SetBuyerAddress : setter
func (baseNegotiation *BaseNegotiation) SetBuyerAddress(buyerAddress cTypes.AccAddress) error {
	baseNegotiation.BuyerAddress = buyerAddress
	return nil
}

// GetSellerAddress : getter
func (baseNegotiation BaseNegotiation) GetSellerAddress() cTypes.AccAddress {
	return baseNegotiation.SellerAddress
}

// SetSellerAddress : setter
func (baseNegotiation *BaseNegotiation) SetSellerAddress(sellerAddress cTypes.AccAddress) error {
	baseNegotiation.SellerAddress = sellerAddress
	return nil
}

// GetPegHash : getter
func (baseNegotiation BaseNegotiation) GetPegHash() types.PegHash { return baseNegotiation.PegHash }

// SetPegHash : setter
func (baseNegotiation *BaseNegotiation) SetPegHash(pegHash types.PegHash) error {
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

func (baseNegotiation BaseNegotiation) GetBuyerContractHash() string {
	return baseNegotiation.BuyerContractHash
}

// SetBuyerContractHash : set SetBuyerNegotiation
func (baseNegotiation *BaseNegotiation) SetBuyerContractHash(buyerContractHash string) error {
	baseNegotiation.BuyerContractHash = buyerContractHash
	return nil
}

// GetSellerContractHash : Get SellerCOntractHash
func (baseNegotiation BaseNegotiation) GetSellerContractHash() string {
	return baseNegotiation.SellerContractHash
}

// SetSellerContractHash : set sellercontracthash
func (baseNegotiation *BaseNegotiation) SetSellerContractHash(sellerContractHash string) error {
	baseNegotiation.SellerContractHash = sellerContractHash
	return nil
}

func (sign Signature) String() string {
	return string(sign)
}

func NewNegotiation(buyerAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash) Negotiation {
	return &BaseNegotiation{
		NegotiationID: NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...)),
		BuyerAddress:  buyerAddress,
		SellerAddress: sellerAddress,
		PegHash:       pegHash,
	}
}

// SignNegotiationBody :
type SignNegotiationBody struct {
	BuyerAddress  cTypes.AccAddress `json:"buyerAddress"`
	SellerAddress cTypes.AccAddress `json:"sellerAddress"`
	PegHash       types.PegHash     `json:"pegHash"`
	Bid           int64             `json:"bid"`
	Time          int64             `json:"time"`
}

// NewSignNegotiationBody :
func NewSignNegotiationBody(buyerAddress, sellerAddress cTypes.AccAddress, peghash types.PegHash, bid, time int64) *SignNegotiationBody {
	return &SignNegotiationBody{
		BuyerAddress:  buyerAddress,
		SellerAddress: sellerAddress,
		PegHash:       peghash,
		Bid:           bid,
		Time:          time,
	}
}

// GetSignBytes :
func (bytes SignNegotiationBody) GetSignBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(bytes)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetNegotiationIDHex : convert NegotiationID string to NegotiationID hex
func GetNegotiationIDFromString(negotiationIDStr string) (negotiationID NegotiationID, err error) {
	bz, err := hex.DecodeString(negotiationIDStr)
	if err != nil {
		return nil, err
	}
	return NegotiationID(bz), nil
}
