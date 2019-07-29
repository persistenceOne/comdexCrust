package acl

//DefineACLBody : Request object to DefineAcl
type DefineACLBody struct {
	From               string `json:"from" valid:"required~Enter the from"`
	Password           string `json:"password" valid:"required~Enter the password"`
	ACLAddress         string `json:"aclAddress" valid:"required~Enter the aclAddress,matches(^cosmos[a-z0-9]{39}$)~aclAddress is Invalid"`
	OrganizationID     string `json:"organizationID" valid:"required~Enter the organizationID, matches(^[A-Fa-f0-9]+$)~Invalid organizationID,length(2|40)~OrganizationID length should be 2 to 40"`
	ZoneID             string `json:"zoneID"  valid:"required~Enter the zoneID, matches(^[A-Fa-f0-9]+$)~Invalid zoneID,length(2|40)~ZoneID length should be 2 to 40"`
	ChainID            string `json:"chainID" valid:"required~Enter the chainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid ChainID"`
	AccountNumber      int64  `json:"accountNumber" `
	Sequence           int64  `json:"sequence"`
	Gas                int64  `json:"gas"`
	GasAdjustment      string `json:"gasAdjustment"`
	IssueAsset         string `json:"issueAsset"  valid:"required~Enter the issueAsset, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid issueAsset"`
	IssueFiat          string `json:"issueFiat"  valid:"required~Enter the issueFiat, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid issueFiat"`
	SendAsset          string `json:"sendAsset"  valid:"required~Enter the sendAsset, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid sendAsset"`
	SendFiat           string `json:"sendFiat"  valid:"required~Enter the sendFiat, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid sendFiat"`
	BuyerExecuteOrder  string `json:"buyerExecuteOrder"  valid:"required~Enter the buyerExecuteOrder, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid buyerExecuteOrder"`
	SellerExecuteOrder string `json:"sellerExecuteOrder"  valid:"required~Enter the sellerExecuteOrder, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid sellerExecuteOrder"`
	ChangeBuyerBid     string `json:"changeBuyerBid"  valid:"required~Enter the changeBuyerBid, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid changeBuyerBid"`
	ChangeSellerBid    string `json:"changeSellerBid"  valid:"required~Enter the changeSellerBid, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid changeSellerBid"`
	ConfirmBuyerBid    string `json:"confirmBuyerBid"  valid:"required~Enter the confirmBuyerBid, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid confirmBuyerBid"`
	ConfirmSellerBid   string `json:"confirmSellerBid"  valid:"required~Enter the confirmSellerBid, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid confirmSellerBid"`
	Negotiation        string `json:"negotiation"  valid:"required~Enter the negotiation, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid negotiation"`
	RedeemAsset        string `json:"redeemAsset"  valid:"required~Enter the redeemAsset, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid redeemAsset"`
	RedeemFiat         string `json:"redeemFiat"  valid:"required~Enter the redeemFiat, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid redeemFiat"`
	ReleaseAsset       string `json:"releaseAsset"  valid:"required~Enter the releaseAsset, matches(^(true|TRUE|True|false|FALSE|False)*$)~Invalid releaseAsset"`
}

//DefineOrganizationBody : request object to define organization
type DefineOrganizationBody struct {
	From           string `json:"from" valid:"required~Enter the from"`
	To             string `json:"to" valid:"required~Enter the to Address,matches(^cosmos[a-z0-9]{39}$)~to Address is Invalid"`
	Password       string `json:"password"  valid:"required~Enter the password"`
	OrganizationID string `json:"organizationID" valid:"required~Enter the organizationID, matches(^[A-Fa-f0-9]+$)~Invalid OrganizationID,length(2|40)~OrganizationID length should be 2 to 40"`
	ChainID        string `json:"chainID" valid:"required~Enter the chainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid chainID"`
	AccountNumber  int64  `json:"accountNumber"`
	Sequence       int64  `json:"sequence"`
	Gas            int64  `json:"gas"`
	GasAdjustment  string `json:"gasAdjustment"`
	ZoneID         string `json:"zoneID" valid:"required~Enter the zoneID, matches(^[A-Fa-f0-9]+$)~Invalid zoneID,length(2|40)~ZoneID length should be 2 to 40"`
}

//DefineZoneBody : request object to add zone
type DefineZoneBody struct {
	From          string `json:"from" valid:"required~Enter the from"`
	To            string `json:"to" valid:"required~Enter the toAddress,matches(^cosmos[a-z0-9]{39}$)~toAddress is Invalid"`
	Password      string `json:"password"  valid:"required~Enter the password"`
	ZoneID        string `json:"zoneID" valid:"required~Enter the zoneID, matches(^[A-Fa-f0-9]+$)~Invalid zoneID,length(2|40)~ZoneID length should be 2 to 40"`
	ChainID       string `json:"chainID" valid:"required~Enter the chainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid chainID"`
	AccountNumber int64  `json:"accountNumber"`
	Sequence      int64  `json:"sequence"`
	Gas           int64  `json:"gas"`
	GasAdjustment string `json:"gasAdjustment"`
}
