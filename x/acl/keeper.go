package acl

import (
	sdk "github.com/commitHub/commitBlockchain/types"
)

//Keeper : acl keeper
type Keeper struct {
	am Mapper
}

//NewKeeper : return a new keeper
func NewKeeper(am Mapper) Keeper {
	return Keeper{am: am}
}

//DefineZoneAddress : Define Zone Address
func (keeper Keeper) DefineZoneAddress(ctx sdk.Context, to sdk.AccAddress, zoneID sdk.ZoneID) (sdk.Tags, sdk.Error) {
	err := keeper.am.SetZone(ctx, to, zoneID)
	if err != nil {
		return nil, err
	}
	tags := sdk.NewTags("zoneAddress", []byte(to.String()))
	tags = tags.AppendTag("zoneID", []byte(zoneID.String()))
	return tags, nil
}

//DefineOrganizationAddress : Define Organization Address
func (keeper Keeper) DefineOrganizationAddress(ctx sdk.Context, to sdk.AccAddress, organizationID sdk.OrganizationID, zoneID sdk.ZoneID) (sdk.Tags, sdk.Error) {
	err := keeper.am.SetOrganization(ctx, to, organizationID, zoneID)
	if err != nil {
		return nil, err
	}
	tags := sdk.NewTags("organizationAddress", []byte(to.String()))
	tags = tags.AppendTag("organizationID", []byte(organizationID.String()))
	return tags, nil
}

//DefineACLAccount : Define ACL Account
func (keeper Keeper) DefineACLAccount(ctx sdk.Context, to sdk.AccAddress, aclAccount sdk.ACLAccount) (sdk.Tags, sdk.Error) {
	err := keeper.am.SetAccount(ctx, to, aclAccount)
	if err != nil {
		return nil, err
	}
	tags := sdk.NewTags("aclAccountAddress", []byte(to.String()))
	return tags, nil
}

//GetZoneDetails :  get zone account address details if they exist
func (keeper Keeper) GetZoneDetails(ctx sdk.Context, zoneID sdk.ZoneID) (sdk.AccAddress, sdk.Error) {
	zoneAccAddress := keeper.am.GetZone(ctx, zoneID)
	if zoneAccAddress == nil {
		return nil, sdk.ErrInvalidAddress("Account ACL details not defined yet.")
	}
	return zoneAccAddress, nil
}

//GetOrganizationDetails :  get organization address details if they exist
func (keeper Keeper) GetOrganizationDetails(ctx sdk.Context, organizationID sdk.OrganizationID) (sdk.Organization, sdk.Error) {
	organization := keeper.am.GetOrganization(ctx, organizationID)
	if &organization == nil {
		return sdk.Organization{}, sdk.ErrInvalidAddress("Account ACL details not defined yet.")
	}
	return organization, nil
}

//GetAccountACLDetails :  get account acl details if they exist
func (keeper Keeper) GetAccountACLDetails(ctx sdk.Context, accAddress sdk.AccAddress) (sdk.ACLAccount, sdk.Error) {
	aclAccount := keeper.am.GetAccount(ctx, accAddress)
	if aclAccount == nil {
		return nil, sdk.ErrInvalidAddress("Account ACL details not defined yet.")
	}
	return aclAccount, nil
}

//CheckZoneAndGetACL : check if the from address is the zone address of the to address and returns back its acl details
func (keeper Keeper) CheckZoneAndGetACL(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) (sdk.ACL, sdk.Error) {
	aclAccount, err := keeper.GetAccountACLDetails(ctx, to)
	if err != nil {
		return sdk.ACL{}, sdk.ErrInternal("To account acl not defined.")
	}
	zoneAddress, err := keeper.GetZoneDetails(ctx, aclAccount.GetZoneID())
	if err != nil {
		return sdk.ACL{}, sdk.ErrInternal("To account zone not found.")
	}
	if zoneAddress.String() != from.String() {
		return sdk.ACL{}, sdk.ErrInternal("Unauthorised transaction.")
	}
	return aclAccount.GetACL(), nil
}

//CheckOrganizationAndGetACL : check if the from address is the organization address of the to address and returns back its acl details
func (keeper Keeper) CheckOrganizationAndGetACL(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) (sdk.ACL, sdk.Error) {
	aclAccount, err := keeper.GetAccountACLDetails(ctx, to)
	if err != nil {
		return sdk.ACL{}, sdk.ErrInternal("To account acl not defined.")
	}
	organizationAddress, err := keeper.GetOrganizationDetails(ctx, aclAccount.GetOrganizationID())
	if err != nil {
		return sdk.ACL{}, sdk.ErrInternal("To account organization not found.")
	}
	if organizationAddress.Address.String() != from.String() {
		return sdk.ACL{}, sdk.ErrInternal("Unauthorised transaction.")
	}
	return aclAccount.GetACL(), nil
}

//CheckIfZoneAccount : check if address belongs to zone account
func (keeper Keeper) CheckIfZoneAccount(ctx sdk.Context, zoneID sdk.ZoneID, accAddress sdk.AccAddress) bool {
	zoneAddress, err := keeper.GetZoneDetails(ctx, zoneID)
	if err != nil {
		return false
	}
	if zoneAddress.String() != accAddress.String() {
		return false
	}
	return true
}

//CheckIfOrganizationAccount : check if address belongs to organization account
func (keeper Keeper) CheckIfOrganizationAccount(ctx sdk.Context, organizationID sdk.OrganizationID, accAddress sdk.AccAddress) bool {
	organizationAddress, err := keeper.GetOrganizationDetails(ctx, organizationID)
	if err != nil {
		return false
	}
	if organizationAddress.Address.String() != accAddress.String() {
		return false
	}
	return true
}
