package types

import (
	"bytes"
	"sort"
)

// TransactionFeedback : type
type TransactionFeedback struct {
	SendAssetsPositiveTx int64 `json:"sendAssetsPositiveTx"`
	SendAssetsNegativeTx int64 `json:"sendAssetsNegativeTx"`
	
	SendFiatsPositiveTx int64 `json:"sendFiatsPositiveTx"`
	SendFiatsNegativeTx int64 `json:"sendFiatsNegativeTx"`
	
	IBCIssueAssetsPositiveTx int64 `json:"ibcIssueAssetsPositiveTx"`
	IBCIssueAssetsNegativeTx int64 `json:"ibcIssueAssetsNegativeTx"`
	
	IBCIssueFiatsPositiveTx int64 `json:"ibcIssueFiatsPositiveTx"`
	IBCIssueFiatsNegativeTx int64 `json:"ibcIssueFiatsNegativeTx"`
	
	BuyerExecuteOrderPositiveTx int64 `json:"buyerExecuteOrderPositiveTx"`
	BuyerExecuteOrderNegativeTx int64 `json:"buyerExecuteOrderNegativeTx"`
	
	SellerExecuteOrderPositiveTx int64 `json:"sellerExecuteOrderPositiveTx"`
	SellerExecuteOrderNegativeTx int64 `json:"sellerExecuteOrderNegativeTx"`
	
	ChangeBuyerBidPositiveTx int64 `json:"changeBuyerBidPositiveTx"`
	ChangeBuyerBidNegativeTx int64 `json:"changeBuyerBidNegativeTx"`
	
	ChangeSellerBidPositiveTx int64 `json:"changeSellerBidPositiveTx"`
	ChangeSellerBidNegativeTx int64 `json:"changeSellerBidNegativeTx"`
	
	ConfirmBuyerBidPositiveTx int64 `json:"confirmBuyerBidPositiveTx"`
	ConfirmBuyerBidNegativeTx int64 `json:"confirmBuyerBidNegativeTx"`
	
	ConfirmSellerBidPositiveTx int64 `json:"confirmSellerBidPositiveTx"`
	ConfirmSellerBidNegativeTx int64 `json:"confirmSellerBidNegativeTx"`
	
	NegotiationPositiveTx int64 `json:"negotiationPositiveTx"`
	NegotiationNegativeTx int64 `json:"negotiationNegativeTx"`
}

// TraderFeedback : traders  traderFeedback this msg
type TraderFeedback struct {
	BuyerAddress  AccAddress `json:"buyerAddress"`
	SellerAddress AccAddress `json:"sellerAddress"`
	PegHash       PegHash    `json:"pegHash"`
	Rating        int64      `json:"rating"`
}

// NewTraderFeedback : create new Rating
func NewTraderFeedback(buyerAddress AccAddress, sellerAddress AccAddress, pegHash PegHash, rating int64) TraderFeedback {
	return TraderFeedback{BuyerAddress: buyerAddress,
		SellerAddress: sellerAddress,
		PegHash:       pegHash,
		Rating:        rating}
}

// GenerateNegotiationID : generates negotiationID from  traderFeedback struct
func (traderFeedback TraderFeedback) GenerateNegotiationID() []byte {
	return GenerateNegotiationIDBytes(traderFeedback.BuyerAddress, traderFeedback.SellerAddress, traderFeedback.PegHash)
}

// TraderFeedbackHistory : A array of  traderFeedbacks
type TraderFeedbackHistory []TraderFeedback

// Sort interface

func (traderFeedbackHistory TraderFeedbackHistory) Len() int { return len(traderFeedbackHistory) }

func (traderFeedbackHistory TraderFeedbackHistory) Less(i, j int) bool {
	return bytes.Compare(traderFeedbackHistory[i].GenerateNegotiationID(), traderFeedbackHistory[j].GenerateNegotiationID()) < 0
}

func (traderFeedbackHistory TraderFeedbackHistory) Swap(i, j int) {
	traderFeedbackHistory[i], traderFeedbackHistory[j] = traderFeedbackHistory[j], traderFeedbackHistory[i]
}

var _ sort.Interface = TraderFeedbackHistory{}

// Sort is a helper function to sort the set of  traderFeedbacks inplace
func (traderFeedbackHistory TraderFeedbackHistory) Sort() TraderFeedbackHistory {
	sort.Sort(traderFeedbackHistory)
	return traderFeedbackHistory
}

// Search : searches if the element is in the array returns length of array if element is not found
func (traderFeedbackHistory TraderFeedbackHistory) Search(incomingFeedback TraderFeedback) int {
	index := sort.Search(traderFeedbackHistory.Len(), func(i int) bool {
		return bytes.Compare(traderFeedbackHistory[i].GenerateNegotiationID(), incomingFeedback.GenerateNegotiationID()) != -1
	})
	return index
}

// AccountReputation : implements basefeedback
type AccountReputation interface {
	GetAddress() AccAddress
	SetAddress(AccAddress) error
	
	GetTransactionFeedback() TransactionFeedback
	SetTransactionFeedback(TransactionFeedback) error
	
	GetTraderFeedbackHistory() TraderFeedbackHistory
	SetTraderFeedbackHistory(TraderFeedbackHistory) error
	
	AddTraderFeedback(TraderFeedback) Error
	
	GetRating() int64
}

// NewAccountReputation := creates new  traderFeedback
func NewAccountReputation() AccountReputation {
	baseAccountReputation := NewBaseAccountReputation()
	return &baseAccountReputation
}

var _ AccountReputation = (*BaseAccountReputation)(nil)

// BaseAccountReputation : base  account reputation
type BaseAccountReputation struct {
	Address               AccAddress            `json:"address"`
	TransactionFeedback   TransactionFeedback   `json:"transactionFeedback"`
	TraderFeedbackHistory TraderFeedbackHistory `json:"traderFeedbackHistory"`
}

// NewBaseAccountReputation : creates new
func NewBaseAccountReputation() BaseAccountReputation {
	return BaseAccountReputation{
		TransactionFeedback:   TransactionFeedback{},
		TraderFeedbackHistory: TraderFeedbackHistory{},
	}
}

// ProtoBaseAccountReputation : converts concrete to order in
func ProtoBaseAccountReputation() AccountReputation {
	return &BaseAccountReputation{}
}

// GetAddress : gets
func (baseAccountReputation BaseAccountReputation) GetAddress() AccAddress {
	return baseAccountReputation.Address
}

// SetAddress : sets
func (baseAccountReputation *BaseAccountReputation) SetAddress(addr AccAddress) error {
	baseAccountReputation.Address = addr
	return nil
}

// GetTransactionFeedback : gets
func (baseAccountReputation BaseAccountReputation) GetTransactionFeedback() TransactionFeedback {
	return baseAccountReputation.TransactionFeedback
}

// SetTransactionFeedback : sets
func (baseAccountReputation *BaseAccountReputation) SetTransactionFeedback(transactionFeedback TransactionFeedback) error {
	baseAccountReputation.TransactionFeedback = transactionFeedback
	return nil
}

// GetTraderFeedbackHistory : gets
func (baseAccountReputation BaseAccountReputation) GetTraderFeedbackHistory() TraderFeedbackHistory {
	return baseAccountReputation.TraderFeedbackHistory
}

// SetTraderFeedbackHistory : sets
func (baseAccountReputation *BaseAccountReputation) SetTraderFeedbackHistory(traderFeedbacks TraderFeedbackHistory) error {
	baseAccountReputation.TraderFeedbackHistory = traderFeedbacks
	return nil
}

// hasTraderFeedback : gets rating
func (baseAccountReputation BaseAccountReputation) hasTraderFeedback(traderFeedback TraderFeedback) bool {
	index := baseAccountReputation.TraderFeedbackHistory.Search(traderFeedback)
	return index < baseAccountReputation.TraderFeedbackHistory.Len() && bytes.Compare(baseAccountReputation.TraderFeedbackHistory[index].GenerateNegotiationID(), traderFeedback.GenerateNegotiationID()) != 0
}

// AddTraderFeedback : sets rating
func (baseAccountReputation *BaseAccountReputation) AddTraderFeedback(traderFeedback TraderFeedback) Error {
	ok := baseAccountReputation.hasTraderFeedback(traderFeedback)
	if !ok {
		baseAccountReputation.TraderFeedbackHistory = append(baseAccountReputation.TraderFeedbackHistory, traderFeedback)
		baseAccountReputation.TraderFeedbackHistory = baseAccountReputation.TraderFeedbackHistory.Sort()
		return nil
	}
	return ErrFeedbackCannotRegister("You have already given a  traderFeedback for this transaction")
}

// GetRating : gets a rating of that account
func (baseAccountReputation BaseAccountReputation) GetRating() int64 { return 100 }

// ReputationDecoder : decoder function for Reputation
type ReputationDecoder func(reputationBytes []byte) (AccountReputation, error)
