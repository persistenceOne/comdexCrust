package negotiation

//ChangeBuyerBidBody : a struct to handle rest request
type ChangeBuyerBidBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	Bid           int64  `json:"bid" valid:"required~Enter the Valid Bid,matches(^[1-9]{1}[0-9]*$)~Invalid Bid"`
	Time          int64  `json:"time" valid:"required~Enter the Time,matches(^[1-9]{1}[0-9]*$)~Enter valid Time"`
	PegHash       string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
	Gas           int64  `json:"gas"`
	GasAdjustment string `json:"gasAdjustment"`
}

//ChangeSellerBidBody : a struct to handle rest request
type ChangeSellerBidBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	Bid           int64  `json:"bid" valid:"required~Enter the Valid Bid,matches(^[1-9]{1}[0-9]*$)~Invalid Bid"`
	Time          int64  `json:"time" valid:"required~Enter the Time,matches(^[1-9]{1}[0-9]*$)~Enter valid Time"`
	PegHash       string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
	Gas           int64  `json:"gas"`
	GasAdjustment string `json:"gasAdjustment"`
}

//ConfirmBuyerBidBody : a struct to handle rest request
type ConfirmBuyerBidBody struct {
	From              string `json:"from" valid:"required~Enter the FromName"`
	Password          string `json:"password" valid:"required~Enter the Password"`
	To                string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	Bid               int64  `json:"bid" valid:"required~Enter the Valid Bid,matches(^[1-9]{1}[0-9]*$)~Invalid Bid"`
	Time              int64  `json:"time" valid:"required~Enter the Time,matches(^[1-9]{1}[0-9]*$)~Enter valid Time"`
	PegHash           string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	ChainID           string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	BuyerContractHash string `json:"buyerContractHash" valid:"required~Enter the BuyerContractHash, matches(^[a-fA-F0-9]{40}$)~Invalid BuyerContractHash"`
	AccountNumber     int64  `json:"accountNumber"`
	Sequence          int64  `json:"sequence"`
	Gas               int64  `json:"gas"`
	GasAdjustment     string `json:"gasAdjustment"`
}

//ConfirmSellerBidBody : a struct to handle rest request
type ConfirmSellerBidBody struct {
	From               string `json:"from" valid:"required~Enter the FromName"`
	Password           string `json:"password" valid:"required~Enter the Password"`
	To                 string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	Bid                int64  `json:"bid" valid:"required~Enter the Valid Bid,matches(^[1-9]{1}[0-9]*$)~Invalid Bid"`
	Time               int64  `json:"time" valid:"required~Enter the Time,matches(^[1-9]{1}[0-9]*$)~Enter valid Time"`
	PegHash            string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	ChainID            string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	SellerContractHash string `json:"sellerContractHash" valid:"required~Enter the SellerContractHash, matches(^[a-fA-F0-9]{40}$)~Invalid SellerContractHash"`
	AccountNumber      int64  `json:"accountNumber"`
	Sequence           int64  `json:"sequence"`
	Gas                int64  `json:"gas"`
	GasAdjustment      string `json:"gasAdjustment"`
}

//NegotiaitonBody : a struct to handle rest request
type NegotiaitonBody struct {
	From          string `json:"from"`
	To            string `json:"to"`
	ChainID       string `json:"chainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
	Gas           int64  `json:"gas"`
	PegHash       string `json:"pegHash"`
	GasAdjustment string `json:"gasAdjustment"`
}
