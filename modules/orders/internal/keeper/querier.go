package keeper

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/types"
)

const (
	QueryOrder = "orderQuery"
)

func NewQuerier(k Keeper) cTypes.Querier {
	return func(ctx cTypes.Context, path []string, req abciTypes.RequestQuery) (res []byte, err cTypes.Error) {
		switch path[0] {
		case QueryOrder:
			return queryOrder(ctx, path[1:], k)
		default:
			return nil, cTypes.ErrUnknownRequest("unknown negotiation query endpoint")
		}
	}
}

func queryOrder(ctx cTypes.Context, path []string, k Keeper) ([]byte, cTypes.Error) {

	negotiationID, err := types.GetNegotiationIDFromString(path[0])
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to parse the negotiationID %s", err))
	}

	order := k.GetOrder(ctx, negotiationID)
	if order == nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf(" order not found"))
	}

	res, err := codec.MarshalJSONIndent(k.cdc, order)

	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to marshal data %s", err.Error()))
	}
	return res, nil
}
