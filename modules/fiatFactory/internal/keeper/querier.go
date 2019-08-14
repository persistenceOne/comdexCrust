package keeper

import (
	"fmt"
	
	cTypes "github.com/cosmos/cosmos-sdk/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	
	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/types"
)

const (
	QueryFiat = "queryFiat"
)

func NewQuerier(k Keeper) cTypes.Querier {
	return func(ctx cTypes.Context, path []string, req abciTypes.RequestQuery) (res []byte, err cTypes.Error) {
		switch path[0] {
		case QueryFiat:
			return queryFiat(ctx, path[1:], k)
		default:
			return nil, cTypes.ErrUnknownRequest("unknown fiatFactory query endpoint")
		}
	}
}

func queryFiat(ctx cTypes.Context, path []string, keeper Keeper) ([]byte, cTypes.Error) {
	pegHashHex, err := types.GetAssetPegHashHex(path[1])
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to parse the pegHash %s", err))
	}
	
	fiatPeg, err := keeper.GetFiatPeg(ctx, pegHashHex)
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("%s", err))
	}
	
	res, err := codec.MarshalJSONIndent(keeper.cdc, fiatPeg)
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to marshal data %s", err.Error()))
	}
	return res, nil
}
