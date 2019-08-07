package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/commitHub/commitBlockchain/modules/auth/exported"
)

// AccountKeeper defines the account contract that must be fulfilled when
// creating a x/bank keeper.
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	GetAllAccounts(ctx sdk.Context) []exported.Account
	SetAccount(ctx sdk.Context, acc exported.Account)
	
	GetNextAssetPegHash(ctx sdk.Context) int
	GetNextFiatPegHash(ctx sdk.Context) int
	
	IterateAccounts(ctx sdk.Context, process func(exported.Account) bool)
}

type ReputationKeeper interface {
	SetSendAssetsPositiveTx(ctx sdk.Context)
}
