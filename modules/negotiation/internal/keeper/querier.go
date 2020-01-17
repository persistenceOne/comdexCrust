package keeper

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/persistenceOne/persistenceSDK/codec"
	negotiationTypes "github.com/persistenceOne/persistenceSDK/modules/negotiation/internal/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

const (
	QueryNegotiation = "queryNegotiation"
)

type QueryNegotiationParams struct {
	BuyerAddress  cTypes.AccAddress
	SellerAddress cTypes.AccAddress
	PegHash       types.PegHash
}

func NewQueryNegotiationParams(buyerAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash) QueryNegotiationParams {
	return QueryNegotiationParams{
		BuyerAddress:  buyerAddress,
		SellerAddress: sellerAddress,
		PegHash:       pegHash,
	}
}

func NewQuerier(k Keeper) cTypes.Querier {
	return func(ctx cTypes.Context, path []string, req abciTypes.RequestQuery) (res []byte, err cTypes.Error) {
		switch path[0] {
		case QueryNegotiation:
			return queryNegotiation(ctx, path[1:], k)
		default:
			return nil, cTypes.ErrUnknownRequest("unknown negotiation query endpoint")
		}
	}
}

// query Negotiation handler

func queryNegotiation(ctx cTypes.Context, path []string, k Keeper) ([]byte, cTypes.Error) {

	negotiationKey, err := types.GetNegotiationIDFromString(path[0])
	if err != nil {
		return nil, negotiationTypes.ErrInvalidNegotiationID(negotiationTypes.DefaultCodeSpace, fmt.Sprintf("negotiation with %s "+
			" not found", err.Error()))
	}
	negotiation, err := k.GetNegotiation(ctx, negotiationKey)
	if err != nil {
		return nil, negotiationTypes.ErrInvalidNegotiationID(negotiationTypes.DefaultCodeSpace, fmt.Sprintf("negotiation with %s "+
			" not found", negotiationKey.String()))
	}

	res, errRes := codec.MarshalJSONIndent(negotiationTypes.ModuleCdc, negotiation)
	if errRes != nil {
		return nil, cTypes.ErrInternal(cTypes.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}

	return res, nil
}
