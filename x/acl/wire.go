package acl

import (
	"github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// RegisterWire : Most users shouldn't use this, but this comes handy for tests.
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(DefineACLBody{}, "comdex-blockchain/DefineACLBody", nil)
	cdc.RegisterConcrete(DefineOrganizationBody{}, "comdex-blockchain/DefineOrganizationBody", nil)
	cdc.RegisterConcrete(DefineZoneBody{}, "comdex-blockchain/DefineZoneBody", nil)
}

// RegisterACLAccount :  register acl account type and interface
func RegisterACLAccount(cdc *wire.Codec) {
	cdc.RegisterInterface((*types.ACLAccount)(nil), nil)
	cdc.RegisterConcrete(&types.BaseACLAccount{}, "comdex-blockchain/AclAccount", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	RegisterACLAccount(msgCdc)
}
