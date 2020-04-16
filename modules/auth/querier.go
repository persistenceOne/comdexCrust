package auth

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/modules/auth/types"
)

// creates a querier for auth REST endpoints
func NewQuerier(keeper AccountKeeper) cTypes.Querier {
	return func(ctx cTypes.Context, path []string, req abciTypes.RequestQuery) ([]byte, cTypes.Error) {
		switch path[0] {
		case types.QueryAccount:
			return queryAccount(ctx, req, keeper)
		default:
			return nil, cTypes.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

func queryAccount(ctx cTypes.Context, req abciTypes.RequestQuery, keeper AccountKeeper) ([]byte, cTypes.Error) {
	var params types.QueryAccountParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	account := keeper.GetAccount(ctx, params.Address)
	if account == nil {
		return nil, cTypes.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", params.Address))
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, account)
	if err != nil {
		return nil, cTypes.ErrInternal(cTypes.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
