package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/comdexCrust/modules/acl"
	"github.com/persistenceOne/comdexCrust/modules/auth/exported"
	"github.com/persistenceOne/comdexCrust/types"
)

type NegotiationKeeper interface {
	GetNegotiation(ctx cTypes.Context, id types.NegotiationID) (types.Negotiation, cTypes.Error)
}

type ACLKeeper interface {
	GetAccountACLDetails(ctx cTypes.Context, fromAddress cTypes.AccAddress) (acl.ACLAccount, cTypes.Error)
}

type AccountKeeper interface {
	GetAccount(ctx cTypes.Context, address cTypes.AccAddress) exported.Account
}
