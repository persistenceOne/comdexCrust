package keeper

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/comdexCrust/modules/bank/internal/types"
)

// register bank invariants
func RegisterInvariants(ir cTypes.InvariantRegistry, ak types.AccountKeeper) {
	ir.RegisterRoute(types.ModuleName, "nonnegative-outstanding",
		NonnegativeBalanceInvariant(ak))
}

// NonnegativeBalanceInvariant checks that all accounts in the application have non-negative balances
func NonnegativeBalanceInvariant(ak types.AccountKeeper) cTypes.Invariant {
	return func(ctx cTypes.Context) (string, bool) {
		var msg string
		var count int

		accts := ak.GetAllAccounts(ctx)
		for _, acc := range accts {
			coins := acc.GetCoins()
			if coins.IsAnyNegative() {
				count++
				msg += fmt.Sprintf("\t%s has a negative denomination of %s\n",
					acc.GetAddress().String(),
					coins.String())
			}
		}
		broken := count != 0

		return cTypes.FormatInvariant(types.ModuleName, "nonnegative-outstanding",
			fmt.Sprintf("amount of negative accounts found %d\n%s", count, msg), broken)
	}
}
