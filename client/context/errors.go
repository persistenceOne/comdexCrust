package context

import (
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/pkg/errors"
)

// ErrInvalidAccount returns a standardized error reflecting that a given
// account address does not exist.
func ErrInvalidAccount(addr sdk.AccAddress) error {
	return errors.Errorf(`No account with address %s was found in the state.
Are you sure there has been a transaction involving it?`, addr)
}
