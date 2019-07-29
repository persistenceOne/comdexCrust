package acl

import sdk "github.com/commitHub/commitBlockchain/types"

// InitACL will initialize the acl store
func InitACL(ctx sdk.Context, keeper Keeper) (err error) {
	zoneID := sdk.ZoneID([]byte("ABcd"))
	organizationID := sdk.OrganizationID([]byte("ABCD"))

	toAddress := sdk.AccAddress([]byte("ToAddress"))
	zoneAddress := sdk.AccAddress([]byte("zoneAddress"))
	organizationAddress := sdk.AccAddress([]byte("organizationAddress"))

	aclAccount := DefaultACLAccount(zoneID, organizationID, toAddress)

	keeper.am.SetZone(ctx, zoneAddress, zoneID)
	keeper.am.SetOrganization(ctx, organizationAddress, organizationID, zoneID)
	keeper.am.SetAccount(ctx, aclAccount.GetAddress(), aclAccount)

	return

}

func DefaultACLAccount(zoneID sdk.ZoneID, organizationID sdk.OrganizationID, address sdk.AccAddress) sdk.ACLAccount {
	return &sdk.BaseACLAccount{
		Address:        address,
		ZoneID:         zoneID,
		OrganizationID: organizationID,
		ACL: sdk.ACL{
			IssueAsset:         true,
			IssueFiat:          true,
			SendAsset:          true,
			SendFiat:           true,
			BuyerExecuteOrder:  true,
			SellerExecuteOrder: true,
			ChangeBuyerBid:     true,
			ChangeSellerBid:    true,
			ConfirmBuyerBid:    true,
			ConfirmSellerBid:   true,
			Negotiation:        true,
			RedeemFiat:         true,
			RedeemAsset:        true,
			ReleaseAsset:       true,
		},
	}
}
