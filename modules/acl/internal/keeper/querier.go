package keeper

import (
	"fmt"
	
	cTypes "github.com/cosmos/cosmos-sdk/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	
	"github.com/commitHub/commitBlockchain/codec"
	
	aclTypes "github.com/commitHub/commitBlockchain/modules/acl/internal/types"
)

const (
	QueryZone         = "queryZone"
	QueryOrganization = "queryOrganization"
	QueryACLAccount   = "queryACLAccount"
)

func NewQuerier(k Keeper) cTypes.Querier {
	return func(ctx cTypes.Context, path []string, req abciTypes.RequestQuery) (res []byte, err cTypes.Error) {
		switch path[0] {
		case QueryZone:
			return queryZone(ctx, path[1:], k)
		case QueryOrganization:
			return queryOrganization(ctx, path[1:], k)
		case QueryACLAccount:
			return queryACLAccount(ctx, path[1:], k)
		default:
			return nil, cTypes.ErrUnknownRequest("unknown negotiation query endpoint")
			
		}
		
	}
}

func queryZone(ctx cTypes.Context, path []string, k Keeper) ([]byte, cTypes.Error) {
	
	zoneID, err := aclTypes.GetZoneIDFromString(path[0])
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to parse the zoneID %s", err))
	}
	
	address, errRes := k.GetZoneAddress(ctx, zoneID)
	if errRes != nil {
		return nil, errRes
	}
	
	res, err := codec.MarshalJSONIndent(k.cdc, address)
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to marshal data %s", err.Error()))
	}
	return res, nil
}

func queryOrganization(ctx cTypes.Context, path []string, k Keeper) ([]byte, cTypes.Error) {
	
	organizationID, err := aclTypes.GetOrganizationIDFromString(path[0])
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to parse the organizationID %s", err))
	}
	
	organization, errRes := k.GetOrganization(ctx, organizationID)
	if errRes != nil {
		return nil, errRes
	}
	
	res, err := codec.MarshalJSONIndent(k.cdc, organization)
	
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to marshal data %s", err.Error()))
	}
	return res, nil
	
}

func queryACLAccount(ctx cTypes.Context, path []string, k Keeper) ([]byte, cTypes.Error) {
	
	address, err := cTypes.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to parse the acl address %s", err))
	}
	
	aclAccount, errRes := k.GetACLAccount(ctx, address)
	if errRes != nil {
		return nil, errRes
	}
	
	res, err := codec.MarshalJSONIndent(k.cdc, aclAccount)
	if err != nil {
		return nil, cTypes.ErrInternal(fmt.Sprintf("failed to marshal data %s", err.Error()))
	}
	return res, nil
	
}
