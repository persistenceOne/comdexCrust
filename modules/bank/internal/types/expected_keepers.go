package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceSDK/modules/auth/exported"
)

// AccountKeeper defines the account contract that must be fulfilled when
// creating a x/bank keeper.
type AccountKeeper interface {
	NewAccountWithAddress(ctx cTypes.Context, addr cTypes.AccAddress) exported.Account

	GetAccount(ctx cTypes.Context, addr cTypes.AccAddress) exported.Account
	GetAllAccounts(ctx cTypes.Context) []exported.Account
	SetAccount(ctx cTypes.Context, acc exported.Account)

	GetNextAssetPegHash(ctx cTypes.Context) int
	GetNextFiatPegHash(ctx cTypes.Context) int

	IterateAccounts(ctx cTypes.Context, process func(exported.Account) bool)
}

type ReputationKeeper interface {
	SetSendAssetsPositiveTx(ctx cTypes.Context)
}
