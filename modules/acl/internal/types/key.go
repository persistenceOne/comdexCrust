package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "acl"
	StoreKey     = ModuleName
	RouterKey    = StoreKey
	QuerierRoute = RouterKey
)

var (
	ZoneKey               = []byte{0x01}
	OrganizationKey       = []byte{0x02}
	ACLKey                = []byte{0x03}
	DefaultZoneID         = []byte("zone")
	DefaultOrganizationID = []byte("organization")
)

func GetZoneKey(id ZoneID) []byte {
	return append(ZoneKey, id...)
}

func GetOrganizationKey(id OrganizationID) []byte {
	return append(OrganizationKey, id...)
}

func GetACLAccountKey(address cTypes.AccAddress) []byte {
	return append(ACLKey, address.Bytes()...)
}
