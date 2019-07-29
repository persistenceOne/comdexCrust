package reputation

//SubmitBuyerFeedbackBody : request msg
type SubmitBuyerFeedbackBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	PegHash       string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	Rating        int64  `json:"rating" valid:"required~Enter the Rating,matches(^[1-9][0-9]?$|^100$)~invalid Rating"`
	GasAdjustment string `json:"gasAdjustment"`
	Gas           int64  `json:"gas"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
}

//SubmitSellerFeedbackBody : request msg
type SubmitSellerFeedbackBody struct {
	From          string `json:"from" valid:"required~Enter the FromName"`
	Password      string `json:"password" valid:"required~Enter the Password"`
	To            string `json:"to" valid:"required~Enter the ToAddress,matches(^cosmos[a-z0-9]{39}$)~ToAddress is Invalid"`
	PegHash       string `json:"pegHash" valid:"required~Enter the PegHash,matches(^[A-F0-9]+$)~Invalid PegHash,length(2|40)~PegHash length between 2-40"`
	Rating        int64  `json:"rating" valid:"required~Enter the Rating,matches(^[1-9][0-9]?$|^100$)~invalid Rating"`
	GasAdjustment string `json:"gasAdjustment"`
	Gas           int64  `json:"gas"`
	ChainID       string `json:"chainID" valid:"required~Enter the ChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
}
