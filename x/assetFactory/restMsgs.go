package assetFactory

// IssueAssetBody : request for issue asset rest
type IssueAssetBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	GasAdjustment string `json:"gasAdjustment"`
	DocumentHash  string `json:"documentHash" valid:"required~Enter the DocumentHash"`
	PegHash       string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	AssetType     string `json:"assetType" valid:"required~Enter the assetType,matches(^[A-Za-z]*$)~Invalid assetType"`
	AssetPrice    int64  `json:"assetPrice" valid:"required~Enter the assetPrice,matches(^[1-9]{1}[0-9]*$)~Invalid assetPrice"`
	QuantityUnit  string `json:"quantityUnit" valid:"required~Enter the QuantityUnit,matches(^[A-Za-z]*$)~Invalid QuantityUnit"`
	AssetQuantity int64  `json:"assetQuantity" valid:"required~Enter the AssetQuantity,matches(^[1-9]{1}[0-9]*$)~Invalid AssetQuantity"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	Sequence      int64  `json:"sequence"`
	Gas           int64  `json:"gas"`
}

// ExecuteAssetBody : request for execute asset rest
type ExecuteAssetBody struct {
	Owner         string `json:"owner" valid:"required~Enter the OwnerAddress,matches(^cosmos[a-z0-9]{39}$)~OwnerAddress is Invalid"`
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

// RedeemAssetBody : request for redeem asset rest
type RedeemAssetBody struct {
	Owner         string `json:"owner" valid:"required~Enter the OwnerAddress,matches(^cosmos[a-z0-9]{39}$)~OwnerAddress is Invalid"`
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

// SendAssetBody : request for send asset rest
type SendAssetBody struct {
	Owner         string `json:"owner" valid:"required~Enter the OwnerAddress,matches(^cosmos[a-z0-9]{39}$)~OwnerAddress is Invalid"`
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
