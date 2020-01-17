package exported

import (
	"time"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/persistenceOne/persistenceSDK/types"
)

// Account is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements Account.
type Account interface {
	GetAddress() cTypes.AccAddress
	SetAddress(cTypes.AccAddress) error // errors if already set.

	GetPubKey() crypto.PubKey // can return nil.
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error

	GetCoins() cTypes.Coins
	SetCoins(cTypes.Coins) error

	// Calculates the amount of coins that can be sent to other accounts given
	// the current time.
	SpendableCoins(blockTime time.Time) cTypes.Coins

	// Ensure that account implements stringer
	String() string

	GetAssetPegWallet() types.AssetPegWallet
	SetAssetPegWallet(types.AssetPegWallet) error

	GetFiatPegWallet() types.FiatPegWallet
	SetFiatPegWallet(types.FiatPegWallet) error
}

// VestingAccount defines an account type that vests coins via a vesting schedule.
type VestingAccount interface {
	Account

	// Delegation and undelegation accounting that returns the resulting base
	// coins amount.
	TrackDelegation(blockTime time.Time, amount cTypes.Coins)
	TrackUndelegation(amount cTypes.Coins)

	GetVestedCoins(blockTime time.Time) cTypes.Coins
	GetVestingCoins(blockTime time.Time) cTypes.Coins

	GetStartTime() int64
	GetEndTime() int64

	GetOriginalVesting() cTypes.Coins
	GetDelegatedFree() cTypes.Coins
	GetDelegatedVesting() cTypes.Coins
}
