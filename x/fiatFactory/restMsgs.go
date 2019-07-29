package fiatFactory

//IssueFiatBody : a struct for implementation of issue fiat in fiat chain
type IssueFiatBody struct {
	From              string `json:"from" valid:"required~Enter the FromName"`
	To                string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	GasAdjustment     string `json:"gasAdjustment"`
	PegHash           string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	TransactionID     string `json:"transactionID" valid:"required~Enter the TransactionID"`
	TransactionAmount int64  `json:"transactionAmount" valid:"required~Enter the TransactionAmount,matches(^[1-9]{1}[0-9]*$)~Invalid TransactionAmount"`
	ChainID           string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber     int64  `json:"accountNumber"`
	Password          string `json:"password" valid:"required~Enter the Password"`
	Sequence          int64  `json:"sequence"`
	Gas               int64  `json:"gas"`
}

//SendFiatBody : a struct for implementation of send fiat in fiat chain
type SendFiatBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	OwnerAddress  string `json:"ownerAddress" valid:"required~Enter the OwnerAddress,matches(^cosmos[a-z0-9]{39}$)~OwnerAddress is Invalid"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	AssetPegHash  string `json:"assetPegHash" valid:"required~Enter the AssetPegHash,matches(^[A-F0-9]+$)~Invalid AssetPegHash,length(2|40)~AssetPegHash length between 2-40"`
	FiatPegHash   string `json:"fiatPegHash" valid:"required~Enter the FiatPegHash,matches(^[A-F0-9]+$)~Invalid FiatPegHash,length(2|40)~FiatPegHash length between 2-40"`
	GasAdjustment string `json:"gasAdjustment"`
	Gas           int64  `json:"gas"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
	Amount        int64  `json:"amount" valid:"required~Enter the Valid Amount,matches(^[1-9]{1}[0-9]*$)~Invalid Amount"`
}

//ExecuteFiatBody : a struct for implementation of exectute fiat order in fiat chain
type ExecuteFiatBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	OwnerAddress  string `json:"ownerAddress" valid:"required~Enter the OwnerAddress,matches(^cosmos[a-z0-9]{39}$)~OwnerAddress is Invalid"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	AssetPegHash  string `json:"assetPegHash" valid:"required~Enter the AssetPegHash,matches(^[A-F0-9]+$)~Invalid AssetPegHash,length(2|40)~AssetPegHash length between 2-40"`
	FiatPegHash   string `json:"fiatPegHash" valid:"required~Enter the FiatPegHash,matches(^[A-F0-9]+$)~Invalid FiatPegHash,length(2|40)~FiatPegHash length between 2-40"`
	GasAdjustment string `json:"gasAdjustment"`
	Gas           int64  `json:"gas"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
	Amount        int64  `json:"amount" valid:"required~Enter the Valid Amount,matches(^[1-9]{1}[0-9]*$)~Invalid Amount"`
}

//RedeemFiatBody implement struct for redeem fiat
type RedeemFiatBody struct {
	From           string `json:"from" valid:"required~Enter the FromName"`
	Password       string `json:"password" valid:"required~Enter the Password"`
	OwnerAddress   string `json:"ownerAddress" valid:"required~Enter the OwnerAddress,matches(^cosmos[a-z0-9]{39}$)~OwnerAddress is Invalid"`
	PegHash        string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	GasAdjustment  string `json:"gasAdjustment"`
	Gas            int64  `json:"gas"`
	ChainID        string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber  int64  `json:"accountNumber"`
	Sequence       int64  `json:"sequence"`
	RedeemedAmount int64  `json:"redeemedAmount" valid:"required~Enter the Valid Amount,matches(^[1-9]{1}[0-9]*$)~Invalid Amount"`
}
