package keeper

import (
	"fmt"
	
	cTypes "github.com/cosmos/cosmos-sdk/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	
	"github.com/commitHub/commitBlockchain/codec"
)

const (
	QueryReputation = "reputationQuery"
)

func NewQuerier(k Keeper) cTypes.Querier {
	return func(ctx cTypes.Context, path []string, req abciTypes.RequestQuery) (res []byte, err cTypes.Error) {
		switch path[0] {
		case QueryReputation:
			return queryReputation(ctx, path[1:], k)
		default:
			return nil, cTypes.ErrUnknownRequest("unknown negotiation query endpoint")
		}
	}
}

func queryReputation(ctx cTypes.Context, path []string, k Keeper) ([]byte, cTypes.Error) {
	
	address, err := cTypes.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to parse the address %s", err))
	}
	account := k.GetAccountReputation(ctx, address)
	
	res, err := codec.MarshalJSONIndent(k.cdc, account)
	
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to marshal data %s", err.Error()))
	}
	return res, nil
}
