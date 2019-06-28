package acl

import (
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// Mapper :  Acl mapper object
type Mapper struct {
	key   sdk.StoreKey
	proto func() sdk.ACLAccount
	cdc   *wire.Codec
}

// NewACLMapper :  return Acl Mapper object
func NewACLMapper(cdc *wire.Codec, key sdk.StoreKey, proto func() sdk.ACLAccount) Mapper {
	return Mapper{
		key:   key,
		proto: proto,
		cdc:   cdc,
	}
}

// AccountStoreKey : create the key for acl account
func AccountStoreKey(addr sdk.AccAddress) []byte {
	return append([]byte("address:"), addr.Bytes()...)
}

// ZoneStoreKey : create the key for zone
func ZoneStoreKey(zoneID sdk.ZoneID) []byte {
	return append([]byte("zoneID:"), zoneID.Bytes()...)
}

// OrganizationStoreKey : create the key for organization
func OrganizationStoreKey(organizationID sdk.OrganizationID) []byte {
	return append([]byte("organizationID:"), organizationID.Bytes()...)
}

// GetAccount : Get Acl Account
func (am Mapper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) sdk.ACLAccount {
	store := ctx.KVStore(am.key)
	bz := store.Get(AccountStoreKey(addr))
	if bz == nil {
		return nil
	}
	acc := am.decodeAccount(bz)
	return acc
}

// GetZone : Get Acl zone
func (am Mapper) GetZone(ctx sdk.Context, zoneID sdk.ZoneID) sdk.AccAddress {
	store := ctx.KVStore(am.key)
	bz := store.Get(ZoneStoreKey(zoneID))
	if bz == nil {
		return nil
	}
	accAddress := sdk.AccAddress(bz)
	return accAddress
}

// GetOrganization : Get Acl organization
func (am Mapper) GetOrganization(ctx sdk.Context, organizationID sdk.OrganizationID) sdk.Organization {
	store := ctx.KVStore(am.key)
	bz := store.Get(OrganizationStoreKey(organizationID))
	if bz == nil {
		return sdk.Organization{}
	}
	org := am.decodeOrganization(bz)
	return org
}

// SetAccount : Set the Acl Account
func (am Mapper) SetAccount(ctx sdk.Context, Addr sdk.AccAddress, acl sdk.ACLAccount) sdk.Error {
	store := ctx.KVStore(am.key)
	bz := am.encodeAccount(acl)
	store.Set(AccountStoreKey(Addr), bz)
	return nil
}

// SetZone : Set the Acl zone
func (am Mapper) SetZone(ctx sdk.Context, accAddress sdk.AccAddress, zoneID sdk.ZoneID) sdk.Error {
	store := ctx.KVStore(am.key)
	bz := accAddress.Bytes()
	store.Set(ZoneStoreKey(zoneID), bz)
	return nil
}

// SetOrganization : Set the Acl organization
func (am Mapper) SetOrganization(ctx sdk.Context, accAddress sdk.AccAddress, organizationID sdk.OrganizationID, zoneID sdk.ZoneID) sdk.Error {
	store := ctx.KVStore(am.key)
	org := sdk.Organization{
		Address: accAddress,
		ZoneID:  zoneID,
	}
	bz := am.encodeOrganization(org)
	store.Set(OrganizationStoreKey(organizationID), bz)
	return nil
}

// IterateAccounts : iterate over account in kv store and add accounts
func (am Mapper) IterateAccounts(ctx sdk.Context, process func(sdk.ACLAccount) (stop bool)) {
	store := ctx.KVStore(am.key)
	iter := sdk.KVStorePrefixIterator(store, []byte("address:"))
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc := am.decodeAccount(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}

// IterateZones : iterate over zones in kv store and process them
func (am Mapper) IterateZones(ctx sdk.Context, process func(sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(am.key)
	iter := sdk.KVStorePrefixIterator(store, []byte("zoneID:"))
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		accAddress := sdk.AccAddress(val)
		if process(accAddress) {
			return
		}
		iter.Next()
	}
}

// IterateOrganizations : iterate over organizations in kv store and process them
func (am Mapper) IterateOrganizations(ctx sdk.Context, process func(sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(am.key)
	iter := sdk.KVStorePrefixIterator(store, []byte("organizationID:"))
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		accAddress := sdk.AccAddress(val)
		if process(accAddress) {
			return
		}
		iter.Next()
	}
}

// encodeAccount : encode and return []byte of Acl Account
func (am Mapper) encodeAccount(acc sdk.ACLAccount) []byte {
	bz, err := am.cdc.MarshalBinaryBare(acc)
	if err != nil {
		panic(err)
	}
	return bz
}

// decodeAccount : decode and return Acl Account type of []bytes
func (am Mapper) decodeAccount(bz []byte) (acc sdk.ACLAccount) {
	err := am.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}

func (am Mapper) encodeOrganization(organization sdk.Organization) []byte {
	bz, err := am.cdc.MarshalBinaryBare(organization)
	if err != nil {
		panic(err)
	}
	return bz
}
func (am Mapper) decodeOrganization(bz []byte) (organization sdk.Organization) {
	err := am.cdc.UnmarshalBinaryBare(bz, &organization)
	if err != nil {
		panic(err)
	}
	return
}
