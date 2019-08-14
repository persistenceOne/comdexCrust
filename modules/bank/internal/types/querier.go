package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryBalanceParams struct {
	Address cTypes.AccAddress
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryBalanceParams(addr cTypes.AccAddress) QueryBalanceParams {
	return QueryBalanceParams{Address: addr}
}
