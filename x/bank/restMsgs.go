package bank

import (
	sdk "github.com/comdex-blockchain/types"
)

// BuyerExecuteOrderBody : request body for buyer execute order
type BuyerExecuteOrderBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	BuyerAddress  string `json:"buyerAddress" valid:"required~Enter the BuyerAddress,matches(^cosmos[a-z0-9]{39}$)~BuyerAddress is Invalid"` // Check Match(BuyerAddr,SellerAddr)
	SellerAddress string `json:"sellerAddress" valid:"required~Enter the SellerAddress,matches(^cosmos[a-z0-9]{39}$)~SellerAddress is Invalid"`
	PegHash       string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	FiatProofHash string `json:"fiatProofHash" valid:"required~Mandatory parameter FiatProofHash missing,matches(^[A-Za-z0-9]+$)~Invalid fiatProofHash,length(2|40)~fiatProofHash length must be between 2-40"`
	GasAdjustment string `json:"gasAdjustment"`
	Gas           int64  `json:"gas"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
}

// SellerExecuteOrderBody : request body for seller execute order
type SellerExecuteOrderBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	BuyerAddress  string `json:"buyerAddress" valid:"required~Enter the BuyerAddress,matches(^cosmos[a-z0-9]{39}$)~BuyerAddress is Invalid"` // Check Match(BuyerAddr,SellerAddr)
	SellerAddress string `json:"sellerAddress" valid:"required~Enter the SellerAddress,matches(^cosmos[a-z0-9]{39}$)~SellerAddress is Invalid"`
	AWBProofHash  string `json:"awbProofHash" valid:"required~Mandatory parameter awbProofHash missing,matches(^[A-Za-z0-9]+$)~Invalid awbProofHash,length(2|40)~awbProofHash length must be between 2-40"`
	PegHash       string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	GasAdjustment string `json:"gasAdjustment"`
	Gas           int64  `json:"gas"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
}

// IssueAssetBody : request for issue asset rest
type IssueAssetBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	To            string `json:"to" valid:"matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	GasAdjustment string `json:"gasAdjustment"`
	DocumentHash  string `json:"documentHash" valid:"required~Enter the DocumentHash"`
	AssetType     string `json:"assetType" valid:"required~Enter the assetType,matches(^[A-Za-z ]*$)~Invalid AssetType"`
	AssetPrice    int64  `json:"assetPrice" valid:"required~Enter the assetPrice,matches(^[1-9]{1}[0-9]*$)~Invalid assetPrice"`
	QuantityUnit  string `json:"quantityUnit" valid:"required~Enter the QuantityUnit,matches(^[A-Za-z]*$)~Invalid QuantityUnit"`
	AssetQuantity int64  `json:"assetQuantity" valid:"required~Enter the AssetQuantity,matches(^[1-9]{1}[0-9]*$)~Invalid AssetQuantity"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	Sequence      int64  `json:"sequence"`
	Gas           int64  `json:"gas"`
	Moderated     bool   `json:"moderated"`
	TakerAddress  string `json:"takerAddress" valid:"matches(^cosmos[a-z0-9]{39}$)~TakerAddress is Invalid"`
}

// IssueFiatBody : request for issue fiat rest
type IssueFiatBody struct {
	From              string `json:"from" valid:"required~Enter the FromName"`
	To                string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	GasAdjustment     string `json:"gasAdjustment"`
	TransactionID     string `json:"transactionID" valid:"required~Enter the TransactionID,  matches(^[A-Z0-9]+$)~transactionID is Invalid,length(2|40)~TransactionID length should be 2 to 40"`
	TransactionAmount int64  `json:"transactionAmount" valid:"required~Enter the TransactionAmount,matches(^[1-9]{1}[0-9]*$)~Invalid TransactionAmount"`
	ChainID           string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber     int64  `json:"accountNumber"`
	Password          string `json:"password" valid:"required~Enter the Password"`
	Sequence          int64  `json:"sequence"`
	Gas               int64  `json:"gas"`
}

// SendAssetBody : body to sendasset request
type SendAssetBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	PegHash       string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	GasAdjustment string `json:"gasAdjustment"`
	Gas           int64  `json:"gas"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
}

// SendFiatBody : body for send fiat request
type SendFiatBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	PegHash       string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	GasAdjustment string `json:"gasAdjustment"`
	Gas           int64  `json:"gas"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
	Amount        int64  `json:"amount" valid:"required~Enter the Valid Amount,matches(^[1-9]{1}[0-9]*$)~Invalid Amount"`
}

// SendTxBody is a body for sendTx
type SendTxBody struct {
	Amount        sdk.Coins `json:"amount"`
	From          string    `json:"from"`
	To            string    `json:"to"`
	Password      string    `json:"password"`
	ChainID       string    `json:"chainID"`
	AccountNumber int64     `json:"accountNumber"`
	Sequence      int64     `json:"sequence"`
	Gas           int64     `json:"gas"`
	GasAdjustment string    `json:"gasAdjustment"`
}

// RedeemAssetBody : request for redeem asset rest
type RedeemAssetBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	GasAdjustment string `json:"gasAdjustment"`
	PegHash       string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	Sequence      int64  `json:"sequence"`
	Gas           int64  `json:"gas"`
}

// RedeemFiatBody : request for redeem fiat rest
type RedeemFiatBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	GasAdjustment string `json:"gasAdjustment"`
	Gas           int64  `json:"gas"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
	RedeemAmount  int64  `json:"redeemAmount" valid:"required~Enter the Valid Amount,matches(^[1-9]{1}[0-9]*$)~Invalid Amount"`
}

// ReleaseAssetBody : request for issue asset rest
type ReleaseAssetBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	PegHash       string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	GasAdjustment string `json:"gasAdjustment"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	Sequence      int64  `json:"sequence"`
	Gas           int64  `json:"gas"`
}
