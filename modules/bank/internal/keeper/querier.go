package keeper

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/modules/bank/internal/types"
)

const (
	QueryBalance = "balances"
)

// NewQuerier returns a new cTypes.Keeper instance.
func NewQuerier(k Keeper) cTypes.Querier {
	return func(ctx cTypes.Context, path []string, req abciTypes.RequestQuery) ([]byte, cTypes.Error) {
		switch path[0] {
		case QueryBalance:
			return queryBalance(ctx, req, k)

		default:
			return nil, cTypes.ErrUnknownRequest("unknown bank query endpoint")
		}
	}
}

// queryBalance fetch an account's balance for the supplied height.
// Height and account address are passed as first and second path components respectively.
func queryBalance(ctx cTypes.Context, req abciTypes.RequestQuery, k Keeper) ([]byte, cTypes.Error) {
	var params types.QueryBalanceParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, k.GetCoins(ctx, params.Address))
	if err != nil {
		return nil, cTypes.ErrInternal(cTypes.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
