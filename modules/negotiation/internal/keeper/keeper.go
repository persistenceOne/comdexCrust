package keeper

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/auth"
	negotiationTypes "github.com/commitHub/commitBlockchain/modules/negotiation/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

type Keeper struct {
	storeKey      cTypes.StoreKey
	accountKeeper auth.AccountKeeper
	cdc           *codec.Codec
}

func NewKeeper(storeKey cTypes.StoreKey, ak auth.AccountKeeper, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:      storeKey,
		accountKeeper: ak,
		cdc:           cdc,
	}
}

// negotiation/{0x01}/{buyerAddress+sellerAddress+pegHash} => negotiation
func (k Keeper) SetNegotiation(ctx cTypes.Context, negotiation negotiationTypes.Negotiation) {
	store := ctx.KVStore(k.storeKey)
	
	negotiationKey := negotiationTypes.GetNegotiationKey(negotiation.GetNegotiationID())
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(negotiation)
	store.Set(negotiationKey, bz)
}

// returns negotiation by negotiationID
func (k Keeper) GetNegotiation(ctx cTypes.Context, negotiationID negotiationTypes.NegotiationID) (negotiation negotiationTypes.Negotiation, err cTypes.Error) {
	store := ctx.KVStore(k.storeKey)
	
	negotiationKey := negotiationTypes.GetNegotiationKey(negotiationID)
	bz := store.Get(negotiationKey)
	if bz == nil {
		return nil, negotiationTypes.ErrInvalidNegotiationID(negotiationTypes.DefaultCodeSpace, "negotiation not found.")
	}
	
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &negotiation)
	return negotiation, nil
}

// get all negotiations => []Negotiations from store
func (k Keeper) GetNegotiations(ctx cTypes.Context) (negotiations []negotiationTypes.Negotiation) {
	k.IterateNegotiations(ctx, func(negotiation negotiationTypes.Negotiation) (stop bool) {
		negotiations = append(negotiations, negotiation)
		return false
	},
	)
	return
}

func (k Keeper) IterateNegotiations(ctx cTypes.Context, handler func(negotiation negotiationTypes.Negotiation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	
	iterator := cTypes.KVStorePrefixIterator(store, negotiationTypes.NegotiationKey)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var negotiation negotiationTypes.Negotiation
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &negotiation)
		if handler(negotiation) {
			break
		}
	}
}

func (k Keeper) GetNegotiatorAccount(ctx cTypes.Context, address cTypes.AccAddress) auth.Account {
	account := k.accountKeeper.GetAccount(ctx, address)
	return account
}

func (k Keeper) GetNegotiationDetails(ctx cTypes.Context, buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress,
	hash types.PegHash) (negotiationTypes.Negotiation, cTypes.Error) {
	
	negotiationID := negotiationTypes.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), hash.Bytes()...))
	_negotiation, err := k.GetNegotiation(ctx, negotiationID)
	if err != nil {
		return nil, err
	}
	return _negotiation, nil
}
