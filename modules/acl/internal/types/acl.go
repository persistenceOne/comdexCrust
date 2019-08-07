package types

import (
	"encoding/hex"
	"fmt"
	
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/common"
)

type ZoneID = common.HexBytes

func GetZoneIDFromString(zoneID string) (ZoneID, error) {
	bz, err := hex.DecodeString(zoneID)
	if err != nil {
		return nil, err
	}
	return ZoneID(bz), nil
}

type OrganizationID = common.HexBytes

type Organization struct {
	Address cTypes.AccAddress `json:"address"`
	ZoneID  ZoneID            `json:"zoneID"`
}

func NewOrganization(address cTypes.AccAddress, id ZoneID) Organization {
	return Organization{
		Address: address,
		ZoneID:  id,
	}
}

func (org Organization) String() string {
	return fmt.Sprintf(`
Address: %s
ZoneID: %s
`, org.Address.String(), org.ZoneID.String())
}

func GetOrganizationIDFromString(organizationID string) (OrganizationID, error) {
	bz, err := hex.DecodeString(organizationID)
	if err != nil {
		return nil, err
	}
	return OrganizationID(bz), nil
}

type ACL struct {
	IssueAsset         bool `json:"issueAsset" valid:"required~Mandatory parameter issueAssets missing"`
	IssueFiat          bool `json:"issueFiat" valid:"required~Mandatory parameter issueFiats missing"`
	SendAsset          bool `json:"sendAsset" valid:"required~Mandatory parameter sendAssets missing"`
	SendFiat           bool `json:"sendFiat" valid:"required~Mandatory parameter sendFiats missing"`
	BuyerExecuteOrder  bool `json:"buyerExecuteOrder" valid:"required~Mandatory parameter buyerExecuteOrder missing"`
	SellerExecuteOrder bool `json:"sellerExecuteOrder" valid:"required~Mandatory parameter sellerExecuteOrder missing"`
	ChangeBuyerBid     bool `json:"changeBuyerBid" valid:"required~Mandatory parameter changeBuyerBid missing"`
	ChangeSellerBid    bool `json:"changeSellerBid" valid:"required~Mandatory parameter changeSellerBid missing"`
	ConfirmBuyerBid    bool `json:"confirmBuyerBid" valid:"required~Mandatory parameter confirmBuyerBid missing"`
	ConfirmSellerBid   bool `json:"confirmSellerBid" valid:"required~Mandatory parameter confirmSellerBid missing"`
	Negotiation        bool `json:"negotiation" valid:"required~Mandatory parameter negotiation missing"`
	RedeemFiat         bool `json:"redeemFiat" valid:"required~Mandatory parameter redeemFiat missing"`
	RedeemAsset        bool `json:"redeemAsset" valid:"required~Mandatory parameter redeemAsset missing"`
	ReleaseAsset       bool `json:"releaseAsset" valid:"required~Mandatory parameter releaseAsset missing"`
}

func (acl ACL) String() string {
	return fmt.Sprintf(`
IssueAsset: %t
IssueFiat: %t
SendAsset: %t
SendFiat: %t
BuyerExecuteOrder: %t
SellerExecuteOrder: %t
ChangeBuyerBid: %t
ChangeSellerBid: %t
ConfirmBuyerBid: %t
ConfirmSellerBid: %t
Negotiation: %t
RedeemFiat: %t
RedeemAsset: %t
ReleaseAsset: %t
`, acl.IssueAsset, acl.IssueFiat, acl.SendAsset, acl.SendFiat, acl.BuyerExecuteOrder, acl.SellerExecuteOrder, acl.ChangeBuyerBid, acl.ChangeSellerBid, acl.ConfirmBuyerBid, acl.ConfirmSellerBid,
		acl.Negotiation, acl.RedeemFiat, acl.RedeemAsset, acl.ReleaseAsset)
}

type ACLAccount interface {
	GetAddress() cTypes.AccAddress
	SetAddress(address cTypes.AccAddress) error
	
	GetZoneID() ZoneID
	SetZoneID(id ZoneID) error
	
	GetOrganizationID() OrganizationID
	SetOrganizationID(id OrganizationID) error
	
	GetACL() ACL
	SetACL(acl ACL) error
}

// BaseACLAccount : Acl account type
type BaseACLAccount struct {
	Address        cTypes.AccAddress `json:"address" valid:"required~Mandatory Parameter Address missing,matches(^[A-F0-9]+$)~Parameter Address is Invalid,length(2|40)~ToAddress length between 2-40"`
	ZoneID         ZoneID            `json:"zoneID" valid:"required~matches(^[A-F0-9]+$)~Invalid TOAddress,length(2|40)~ToAddress length between 2-40"`
	OrganizationID OrganizationID    `json:"organizationID" valid:"required~matches(^[A-F0-9]+$)~Invalid TOAddress,length(2|40)~ToAddress length between 2-40"`
	ACL            ACL               `json:"acl"`
}

var _ ACLAccount = (*BaseACLAccount)(nil)

// GetAddress : getter
func (baseACLAccount BaseACLAccount) GetAddress() cTypes.AccAddress {
	return baseACLAccount.Address
}

// SetAddress : setter
func (baseACLAccount *BaseACLAccount) SetAddress(address cTypes.AccAddress) error {
	baseACLAccount.Address = address
	return nil
}

// GetOrganizationID : getter
func (baseACLAccount BaseACLAccount) GetOrganizationID() OrganizationID {
	return baseACLAccount.OrganizationID
}

// SetOrganizationID : setter
func (baseACLAccount *BaseACLAccount) SetOrganizationID(organizationID OrganizationID) error {
	baseACLAccount.OrganizationID = organizationID
	return nil
}

// GetACL : getter
func (baseACLAccount BaseACLAccount) GetACL() ACL {
	return baseACLAccount.ACL
}

// SetACL : setter
func (baseACLAccount *BaseACLAccount) SetACL(acl ACL) error {
	baseACLAccount.ACL = acl
	return nil
}

// GetZoneID : getter
func (baseACLAccount BaseACLAccount) GetZoneID() ZoneID {
	return baseACLAccount.ZoneID
}

// SetZoneID : setter
func (baseACLAccount *BaseACLAccount) SetZoneID(zoneID ZoneID) error {
	baseACLAccount.ZoneID = zoneID
	return nil
}

func (baseACLAccount BaseACLAccount) String() string {
	return fmt.Sprintf(`
Address: %s
ZoneID: %s
OrganizationID: %s
ACL: %s
`, baseACLAccount.GetAddress().String(), baseACLAccount.ZoneID.String(), baseACLAccount.OrganizationID.String(), baseACLAccount.ACL.String())
}

// ACLAccountDecoder : decoder function for acl account
type ACLAccountDecoder func(aclbytes []byte) (ACLAccount, error)
type OrgDecoder func(orgBytes []byte) (Organization, error)

// ProtoBaseACLAccount : prototype of acl account
func ProtoBaseACLAccount() ACLAccount {
	return &BaseACLAccount{}
}
