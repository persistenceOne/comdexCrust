package acl

import (
	"github.com/persistenceOne/persistenceSDK/modules/acl/internal/keeper"
	"github.com/persistenceOne/persistenceSDK/modules/acl/internal/types"
)

const (
	ModuleName   = types.ModuleName
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey

	DefaultCodeSpace = types.DefaultCodeSpace
)

type (
	GenesisState = types.GenesisState
	Keeper       = keeper.Keeper

	AccountKeeper = types.AccountKeeper

	MsgDefineZones         = types.MsgDefineZones
	MsgDefineACLs          = types.MsgDefineACLs
	MsgDefineOrganizations = types.MsgDefineOrganizations

	DefineZone         = types.DefineZone
	DefineOrganization = types.DefineOrganization
	DefineACL          = types.DefineACL
	Organization       = types.Organization

	ACLAccount     = types.ACLAccount
	BaseACLAccount = types.BaseACLAccount
	ZoneID         = types.ZoneID
	OrganizationID = types.OrganizationID

	ACL = types.ACL
)

var (
	DefaultZoneID         = types.DefaultZoneID
	DefaultOrganizationID = types.DefaultOrganizationID

	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	ErrInvalidAddress = types.ErrInvalidAddress
	ErrNoInputs       = types.ErrNoInputs

	GetACLAccountKey   = types.GetACLAccountKey
	GetOrganizationKey = types.GetOrganizationKey
	GetZoneKey         = types.GetZoneKey

	GetZoneIDFromString         = types.GetZoneIDFromString
	GetOrganizationIDFromString = types.GetOrganizationIDFromString

	EventTypeDefineACL          = types.EventTypeDefineACL
	EventTypeDefineOrganization = types.EventTypeDefineOrganization
	EventTypeDefineZone         = types.EventTypeDefineZone

	AttributeKeyZoneID              = types.AttributeKeyZoneID
	AttributeKeyZoneAddress         = types.AttributeKeyZoneAddress
	AttributeKeyOrganizationID      = types.AttributeKeyOrganizationID
	AttributeKeyOrganizationAddress = types.AttributeKeyOrganizationAddress
	AttributeACLAccountAddress      = types.AttributeACLAccountAddress

	NewOrganization = types.NewOrganization
	NewKeeper       = keeper.NewKeeper

	NewQuerier = keeper.NewQuerier
)
