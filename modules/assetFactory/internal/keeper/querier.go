package keeper

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/types"
)

const (
	QueryAsset = "queryAsset"
)

func NewQuerier(k Keeper) cTypes.Querier {
	return func(ctx cTypes.Context, path []string, req abciTypes.RequestQuery) (res []byte, err cTypes.Error) {
		switch path[0] {
		case QueryAsset:
			return queryAsset(ctx, path[1:], k)
		default:
			return nil, cTypes.ErrUnknownRequest("unknown assetFactory query endpoint")
		}
	}
}

func queryAsset(ctx cTypes.Context, path []string, k Keeper) ([]byte, cTypes.Error) {
	pegHashHex, err := types.GetAssetPegHashHex(path[0])
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to parse the pegHash %s", err))
	}

	assetPeg, err := k.GetAssetPeg(ctx, pegHashHex)
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("%s", err))
	}
	res, err := codec.MarshalJSONIndent(k.cdc, assetPeg)
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to marshal data %s", err.Error()))
	}
	return res, nil
}
