package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/commitHub/commitBlockchain/modules/auth/exported"
	
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/negotiation"
)

type NegotiationKeeper interface {
	GetNegotiation(ctx cTypes.Context, id negotiation.NegotiationID) (negotiation.Negotiation, cTypes.Error)
}

type ACLKeeper interface {
	GetAccountACLDetails(ctx cTypes.Context, fromAddress cTypes.AccAddress) (acl.ACLAccount, cTypes.Error)
}

type AccountKeeper interface {
	GetAccount(ctx cTypes.Context, address cTypes.AccAddress) exported.Account
}
