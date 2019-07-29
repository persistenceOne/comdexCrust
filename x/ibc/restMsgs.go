package ibc

//IssueAssetBody : is a msg body for kafka msg
type IssueAssetBody struct {
	// Fees             sdk.Coin  `json="fees"`
	From               string `json:"from" valid:"required~Enter the FromName"`
	To                 string `json:"to" valid:"matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	Password           string `json:"password" valid:"required~Enter the Password"`
	SourceChainID      string `json:"sourceChainID" valid:"required~Enter the SourceChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid SourceChainID"`
	DestinationChainID string `json:"destinationChainID" valid:"required~Enter the DestinationChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid DestinationChainID"`
	AccountNumber      int64  `json:"accountNumber"`
	Sequence           int64  `json:"sequence"`
	Gas                int64  `json:"gas"`
	GasAdjustment      string `json:"gasAdjustment"`
	DocumentHash       string `json:"documentHash" valid:"required~Enter the documentHash"`
	AssetQuantity      int64  `json:"assetQuantity" valid:"required~Enter the AssetQuantity,matches(^[1-9]{1}[0-9]*$)~Invalid AssetQuantity"`
	AssetType          string `json:"assetType" valid:"required~Enter the assetType,matches(^[A-Za-z]*$)~Invalid assetType"`
	AssetPrice         int64  `json:"assetPrice" valid:"required~Enter the assetPrice,matches(^[1-9]{1}[0-9]*$)~Invalid assetPrice"`
	QuantityUnit       string `json:"quantityUnit" valid:"required~Enter the QuantityUnit,matches(^[A-Za-z]*$)~Invalid QuantityUnit"`
	Moderated          bool   `json:"moderated"`
	TakerAddress       string `json:"takerAddress" valid:"matches(^cosmos[a-z0-9]{39}$)~TakerAddress is Invalid"`
}

//IssueFiatBody : is a msg body for kafka msg
type IssueFiatBody struct {
	From               string `json:"from" valid:"required~Enter the FromName"`
	To                 string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	Password           string `json:"password" valid:"required~Enter the Password"`
	PegHash            string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	SourceChainID      string `json:"sourceChainID" valid:"required~Enter the SourceChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid SourceChainID"`
	DestinationChainID string `json:"destinationChainID" valid:"required~Enter the DestinationChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid DestinationChainID"`
	AccountNumber      int64  `json:"accountNumber"`
	Sequence           int64  `json:"sequence"`
	Gas                int64  `json:"gas"`
	GasAdjustment      string `json:"gasAdjustment"`
	TransactionID      string `json:"transactionID" valid:"required~Enter the TransactionID"`
	TransactionAmount  int64  `json:"transactionAmount" valid:"required~Enter the TransactionAmount,matches(^[1-9]{1}[0-9]*$)~Invalid TransactionAmount"`
}

//RedeemAssetBody : is a msg body for kafka msg
type RedeemAssetBody struct {
	From               string `json:"from" valid:"required~Enter the FromName"`
	To                 string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	Password           string `json:"password" valid:"required~Enter the Password"`
	SourceChainID      string `json:"sourceChainID" valid:"required~Enter the SourceChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid SourceChainID"`
	DestinationChainID string `json:"destinationChainID" valid:"required~Enter the DestinationChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid DestinationChainID"`
	AccountNumber      int64  `json:"accountNumber"`
	Sequence           int64  `json:"sequence"`
	Gas                int64  `json:"gas"`
	GasAdjustment      string `json:"gasAdjustment"`
	PegHash            string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
}

//RedeemFiatBody : is a msg body for kafka msg
type RedeemFiatBody struct {
	From               string `json:"from" valid:"required~Enter the FromName"`
	Password           string `json:"password" valid:"required~Enter the Password"`
	To                 string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	SourceChainID      string `json:"sourceChainID" valid:"required~Enter the SourceChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid SourceChainID"`
	DestinationChainID string `json:"destinationChainID" valid:"required~Enter the DestinationChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid DestinationChainID"`
	PegHash            string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	GasAdjustment      string `json:"gasAdjustment"`
	Gas                int64  `json:"gas"`
	AccountNumber      int64  `json:"accountNumber"`
	Sequence           int64  `json:"sequence"`
	RedeemAmount       int64  `json:"redeemAmount" valid:"required~Enter the Valid Amount,matches(^[1-9]{1}[0-9]*$)~Invalid Amount"`
}

//SendAssetBody : body to sendasset request
type SendAssetBody struct {
	From               string `json:"from" valid:"required~Enter the FromName"`
	Password           string `json:"password" valid:"required~Enter the Password"`
	To                 string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	SourceChainID      string `json:"sourceChainID" valid:"required~Enter the SourceChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid SourceChainID"`
	DestinationChainID string `json:"destinationChainID" valid:"required~Enter the DestinationChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid DestinationChainID"`
	PegHash            string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	GasAdjustment      string `json:"gasAdjustment"`
	Gas                int64  `json:"gas"`
	AccountNumber      int64  `json:"accountNumber"`
	Sequence           int64  `json:"sequence"`
}

// SendFiatBody :
type SendFiatBody struct {
	From               string `json:"from" valid:"required~Enter the FromName"`
	Password           string `json:"password" valid:"required~Enter the Password"`
	To                 string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	PegHash            string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	GasAdjustment      string `json:"gasAdjustment"`
	Gas                int64  `json:"gas"`
	SourceChainID      string `json:"sourceChainID" valid:"required~Enter the sourceChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid sourceChainID"`
	DestinationChainID string `json:"destinationChainID" valid:"required~Enter the destinationChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid destinationChainID"`
	AccountNumber      int64  `json:"accountNumber"`
	Sequence           int64  `json:"sequence"`
	Amount             int64  `json:"amount" valid:"required~Enter the Valid Amount,matches(^[1-9]{1}[0-9]*$)~Invalid Amount"`
}

// BuyerExecuteOrderBody :
type BuyerExecuteOrderBody struct {
	From               string `json:"from" valid:"required~Enter the FromName"`
	Password           string `json:"password" valid:"required~Enter the Password"`
	BuyerAddress       string `json:"buyerAddress" valid:"required~Enter the BuyerAddress,matches(^cosmos[a-z0-9]{39}$)~BuyerAddress is Invalid"` //Check Match(BuyerAddr,SellerAddr)
	SellerAddress      string `json:"sellerAddress" valid:"required~Enter the SellerAddress,matches(^cosmos[a-z0-9]{39}$)~SellerAddress is Invalid"`
	PegHash            string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	FiatProofHash      string `json:"fiatProofHash" valid:"required~Mandatory parameter FiatProofHash missing,matches(^[A-Za-z0-9]+$)~Invalid fiatProofHash,length(2|40)~fiatProofHash length must be between 2-40"`
	GasAdjustment      string `json:"gasAdjustment"`
	Gas                int64  `json:"gas"`
	AccountNumber      int64  `json:"accountNumber"`
	Sequence           int64  `json:"sequence"`
	SourceChainID      string `json:"sourceChainID" valid:"required~Enter the sourceChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid sourceChainID"`
	DestinationChainID string `json:"destinationChainID" valid:"required~Enter the destinationChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid destinationChainID"`
}

// SellerExecuteOrderBody :
type SellerExecuteOrderBody struct {
	From               string `json:"from" valid:"required~Enter the FromName"`
	Password           string `json:"password" valid:"required~Enter the Password"`
	BuyerAddress       string `json:"buyerAddress" valid:"required~Enter the BuyerAddress,matches(^cosmos[a-z0-9]{39}$)~BuyerAddress is Invalid"` //Check Match(BuyerAddr,SellerAddr)
	SellerAddress      string `json:"sellerAddress" valid:"required~Enter the SellerAddress,matches(^cosmos[a-z0-9]{39}$)~SellerAddress is Invalid"`
	AWBProofHash       string `json:"awbProofHash" valid:"required~Mandatory parameter awbProofHash missing,matches(^[A-Za-z0-9]+$)~Invalid awbProofHash,length(2|40)~awbProofHash length must be between 2-40"`
	PegHash            string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	GasAdjustment      string `json:"gasAdjustment"`
	Gas                int64  `json:"gas"`
	AccountNumber      int64  `json:"accountNumber"`
	Sequence           int64  `json:"sequence"`
	SourceChainID      string `json:"sourceChainID" valid:"required~Enter the sourceChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid sourceChainID"`
	DestinationChainID string `json:"destinationChainID" valid:"required~Enter the destinationChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid destinationChainID"`
}
