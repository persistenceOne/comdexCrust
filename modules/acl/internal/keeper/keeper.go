package keeper

import (
	"strings"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/codec"
	aclTypes "github.com/commitHub/commitBlockchain/modules/acl/internal/types"
)

type Keeper struct {
	storeKey      cTypes.StoreKey
	accountKeeper aclTypes.AccountKeeper
	cdc           *codec.Codec
}

func NewKeeper(storeKey cTypes.StoreKey, accountKeeper aclTypes.AccountKeeper, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:      storeKey,
		accountKeeper: accountKeeper,
		cdc:           cdc,
	}
}

// Zone

func (keeper Keeper) SetZoneAddress(ctx cTypes.Context, id aclTypes.ZoneID, address cTypes.AccAddress) cTypes.Error {
	store := ctx.KVStore(keeper.storeKey)

	getZoneKey := aclTypes.GetZoneKey(id)

	data := store.Get(getZoneKey)
	if data != nil {
		return aclTypes.ErrInvalidID(aclTypes.DefaultCodeSpace, "zone with this given id already exist")
	}

	bytes := keeper.cdc.MustMarshalBinaryLengthPrefixed(address)
	store.Set(getZoneKey, bytes)

	return nil
}

func (keeper Keeper) GetZoneAddress(ctx cTypes.Context, id aclTypes.ZoneID) (cTypes.AccAddress, cTypes.Error) {
	store := ctx.KVStore(keeper.storeKey)

	getZoneKey := aclTypes.GetZoneKey(id)

	data := store.Get(getZoneKey)
	if data == nil {
		return nil, aclTypes.ErrInvalidID(aclTypes.DefaultCodeSpace, "zone with given address doesn't exist")
	}

	var address cTypes.AccAddress
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(data, &address)

	return address, nil
}

func (keeper Keeper) GetZones(ctx cTypes.Context) []cTypes.AccAddress {
	var zones []cTypes.AccAddress

	store := ctx.KVStore(keeper.storeKey)
	iterator := cTypes.KVStorePrefixIterator(store, aclTypes.ZoneKey)

	for ; iterator.Valid(); iterator.Next() {
		var zone cTypes.AccAddress

		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &zone)
		zones = append(zones, zone)
	}

	return zones
}

// Organization

func (keeper Keeper) SetOrganization(ctx cTypes.Context, id aclTypes.OrganizationID, organization aclTypes.Organization) cTypes.Error {
	store := ctx.KVStore(keeper.storeKey)

	data := store.Get(aclTypes.GetOrganizationKey(id))
	if data != nil {
		return aclTypes.ErrInvalidID(aclTypes.DefaultCodeSpace, "organization with given id already exist")
	}

	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(organization)
	store.Set(aclTypes.GetOrganizationKey(id), bz)

	return nil
}

func (keeper Keeper) GetOrganization(ctx cTypes.Context, id aclTypes.OrganizationID) (aclTypes.Organization, cTypes.Error) {
	store := ctx.KVStore(keeper.storeKey)

	data := store.Get(aclTypes.GetOrganizationKey(id))
	if data == nil {
		return aclTypes.Organization{}, aclTypes.ErrInvalidID(aclTypes.DefaultCodeSpace, "organization with given id doesn't exist")
	}

	var organization aclTypes.Organization
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(data, &organization)

	return organization, nil
}

func (keeper Keeper) GetOrganizations(ctx cTypes.Context) []aclTypes.Organization {
	var organizations []aclTypes.Organization

	store := ctx.KVStore(keeper.storeKey)
	iterator := cTypes.KVStorePrefixIterator(store, aclTypes.OrganizationKey)

	for ; iterator.Valid(); iterator.Next() {
		var organization aclTypes.Organization

		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &organization)
		organizations = append(organizations, organization)
	}

	return organizations
}

func (keeper Keeper) GetOrganizationsByZoneID(ctx cTypes.Context, id aclTypes.ZoneID) []aclTypes.Organization {
	var organizationList []aclTypes.Organization

	organizations := keeper.GetOrganizations(ctx)

	for _, organization := range organizations {
		if strings.EqualFold(organization.ZoneID.String(), id.String()) {
			organizationList = append(organizationList, organization)
		}
	}

	return organizationList
}

// ACL

func (keeper Keeper) SetACLAccount(ctx cTypes.Context, acl aclTypes.ACLAccount) cTypes.Error {
	store := ctx.KVStore(keeper.storeKey)

	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(acl)
	store.Set(aclTypes.GetACLAccountKey(acl.GetAddress()), bz)

	return nil
}

func (keeper Keeper) GetACLAccount(ctx cTypes.Context, address cTypes.AccAddress) (aclTypes.ACLAccount, cTypes.Error) {
	store := ctx.KVStore(keeper.storeKey)

	data := store.Get(aclTypes.GetACLAccountKey(address))
	if data == nil {
		return nil, aclTypes.ErrInvalidID(aclTypes.DefaultCodeSpace, "acl for this account not defined")
	}
	var acl aclTypes.ACLAccount
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(data, &acl)

	return acl, nil
}

func (keeper Keeper) GetACLAccounts(ctx cTypes.Context) []aclTypes.ACLAccount {
	var accounts []aclTypes.ACLAccount

	store := ctx.KVStore(keeper.storeKey)
	iterator := cTypes.KVStorePrefixIterator(store, aclTypes.ACLKey)

	for ; iterator.Valid(); iterator.Next() {
		var account aclTypes.ACLAccount

		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &account)
		accounts = append(accounts, account)
	}

	return accounts
}

func (keeper Keeper) CheckValidZoneAddress(ctx cTypes.Context, id aclTypes.ZoneID, address cTypes.AccAddress) bool {
	zoneAddress, err := keeper.GetZoneAddress(ctx, id)
	if err != nil {
		return false
	}

	if !strings.EqualFold(zoneAddress.String(), address.String()) {
		return false
	}

	return true
}

func (keeper Keeper) CheckValidOrganizationAddress(ctx cTypes.Context, zoneID aclTypes.ZoneID, organizationID aclTypes.OrganizationID, address cTypes.AccAddress) bool {

	organization, err := keeper.GetOrganization(ctx, organizationID)
	if err != nil {
		return false
	}

	if !(strings.EqualFold(organization.ZoneID.String(), zoneID.String()) && strings.EqualFold(organization.Address.String(), address.String())) {
		return false
	}

	return true
}

func (keeper Keeper) CheckValidGenesisAddress(ctx cTypes.Context, address cTypes.AccAddress) bool {
	k := keeper
	account := k.accountKeeper.GetAccount(ctx, address)
	accountNumber := account.GetAccountNumber()

	if accountNumber == uint64(0) {
		return true
	}

	return false
}
func (keeper Keeper) GetAccountACLDetails(ctx cTypes.Context, address cTypes.AccAddress) (aclTypes.ACLAccount, cTypes.Error) {
	aclAccount, err := keeper.GetACLAccount(ctx, address)
	if err != nil {
		return nil, err
	}
	return aclAccount, nil
}

// CheckZoneAndGetACL : check if the from address is the zone address of the to address and returns back its acl details
func (keeper Keeper) CheckZoneAndGetACL(ctx cTypes.Context, from cTypes.AccAddress, to cTypes.AccAddress) (aclTypes.ACL, cTypes.Error) {
	aclAccount, err := keeper.GetAccountACLDetails(ctx, to)
	if err != nil {
		return aclTypes.ACL{}, cTypes.ErrInternal("To account acl not defined.")
	}
	zoneAddress, err := keeper.GetZoneAddress(ctx, aclAccount.GetZoneID())
	if err != nil {
		return aclTypes.ACL{}, cTypes.ErrInternal("To account zone not found.")
	}
	if zoneAddress.String() != from.String() {
		return aclTypes.ACL{}, cTypes.ErrInternal("Unauthorised transaction.")
	}
	return aclAccount.GetACL(), nil
}

func (keeper Keeper) DefineZoneAddress(ctx cTypes.Context, to cTypes.AccAddress, zoneID aclTypes.ZoneID) cTypes.Error {

	err := keeper.SetZoneAddress(ctx, zoneID, to)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			aclTypes.EventTypeDefineZone,
			cTypes.NewAttribute("zoneAddress", to.String()),
			cTypes.NewAttribute("zoneID", zoneID.String()),
		))

	return nil
}

// DefineOrganizationAddress : Define Organization Address
func (keeper Keeper) DefineOrganizationAddress(ctx cTypes.Context, to cTypes.AccAddress,
	organizationID aclTypes.OrganizationID, zoneID aclTypes.ZoneID) cTypes.Error {

	organization := aclTypes.Organization{
		Address: to,
		ZoneID:  zoneID,
	}
	err := keeper.SetOrganization(ctx, organizationID, organization)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			aclTypes.EventTypeDefineOrganization,
			cTypes.NewAttribute("organizationAddress", to.String()),
			cTypes.NewAttribute("organizationID", organizationID.String()),
		))
	return nil
}

func (keeper Keeper) DefineACLAccount(ctx cTypes.Context, toAddress cTypes.AccAddress, aclAccount aclTypes.ACLAccount) cTypes.Error {

	err := keeper.SetACLAccount(ctx, aclAccount)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(cTypes.NewEvent(
		aclTypes.EventTypeDefineACL,
		cTypes.NewAttribute("aclAccountAddress", toAddress.String()),
	))

	return nil
}
