package acl

import sdk "github.com/comdex-blockchain/types"

// InitACL will initialize the acl store
func InitACL(ctx sdk.Context, keeper Keeper) (err error) {
	zoneID := sdk.ZoneID([]byte("ABcd"))
	organizationID := sdk.OrganizationID([]byte("ABCD"))
	
	zoneAddress := sdk.AccAddress([]byte("zoneAddress"))
	organizationAddress := sdk.AccAddress([]byte("organizationAddress"))
	
	aclAccount := DefaultACLAccount(zoneID, organizationID)
	
	keeper.am.SetZone(ctx, zoneAddress, zoneID)
	keeper.am.SetOrganization(ctx, organizationAddress, organizationID, zoneID)
	keeper.am.SetAccount(ctx, aclAccount.GetAddress(), aclAccount)
	
	return
	
}

func DefaultACLAccount(zoneID sdk.ZoneID, organizationID sdk.OrganizationID) sdk.ACLAccount {
	return &sdk.BaseACLAccount{
		Address:        sdk.AccAddress("accAddress"),
		ZoneID:         zoneID,
		OrganizationID: organizationID,
	}
}
