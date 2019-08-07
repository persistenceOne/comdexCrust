package keeper

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/commitHub/commitBlockchain/codec"
	
	"github.com/commitHub/commitBlockchain/types"
	
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/negotiation"
	orderTypes "github.com/commitHub/commitBlockchain/modules/orders/internal/types"
)

type Keeper struct {
	storeKey          cTypes.StoreKey
	cdc               *codec.Codec
	NegotiationKeeper negotiation.Keeper
	ACLKeeper         acl.Keeper
	AccountKeeper     auth.AccountKeeper
}

func NewKeeper(storeKey cTypes.StoreKey, cdc *codec.Codec, negotiationKeeper negotiation.Keeper,
	aclKeeper acl.Keeper, accountKeeper auth.AccountKeeper) Keeper {
	
	return Keeper{
		storeKey:          storeKey,
		cdc:               cdc,
		NegotiationKeeper: negotiationKeeper,
		ACLKeeper:         aclKeeper,
		AccountKeeper:     accountKeeper,
	}
}

func (k Keeper) SetOrder(ctx cTypes.Context, order orderTypes.Order) {
	negotiationID := order.GetNegotiationID()
	store := ctx.KVStore(k.storeKey)
	
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(order)
	if err != nil {
		panic(err)
	}
	storeKey := orderTypes.GetOrderKey(negotiationID)
	store.Set(storeKey, bz)
	
}

func (k Keeper) GetOrder(ctx cTypes.Context, negotiationID negotiation.NegotiationID) orderTypes.Order {
	store := ctx.KVStore(k.storeKey)
	storeKey := orderTypes.GetOrderKey(negotiationID)
	bz := store.Get(storeKey)
	if bz == nil {
		return nil
	}
	var order orderTypes.Order
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &order)
	return order
}

func (k Keeper) IterateOrders(ctx cTypes.Context, process func(orderTypes.Order) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := cTypes.KVStorePrefixIterator(store, orderTypes.OrdersKey)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var order orderTypes.Order
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &order)
		if process(order) {
			break
		}
	}
}

func (k Keeper) NewOrder(buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash) orderTypes.Order {
	order := orderTypes.BaseOrder{}
	negotiationID := negotiation.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	order.SetNegotiationID(negotiationID)
	return &order
}

func (keeper Keeper) SendAssetsToOrder(ctx cTypes.Context, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress, assetPeg types.AssetPeg) cTypes.Error {
	// negotiationID := negotiation.GetOrderKey(toAddress, fromAddress, assetPeg.GetPegHash())
	negotiationID := negotiation.NegotiationID(append(append(toAddress.Bytes(), fromAddress.Bytes()...), assetPeg.GetPegHash().Bytes()...))
	order := keeper.GetOrder(ctx, negotiationID)
	if order == nil {
		order = keeper.NewOrder(toAddress, fromAddress, assetPeg.GetPegHash())
	}
	order.SetAssetPegWallet(types.AddAssetPegToWallet(assetPeg, order.GetAssetPegWallet()))
	keeper.SetOrder(ctx, order)
	return nil
}

// SendFiatsToOrder fiat pegs to order
func (keeper Keeper) SendFiatsToOrder(ctx cTypes.Context, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress, pegHash types.PegHash, fiatPegWallet types.FiatPegWallet) cTypes.Error {
	negotiationID := negotiation.NegotiationID(append(append(fromAddress.Bytes(), toAddress.Bytes()...), pegHash.Bytes()...))
	order := keeper.GetOrder(ctx, negotiationID)
	if order == nil {
		order = keeper.NewOrder(fromAddress, toAddress, pegHash)
	}
	order.SetFiatPegWallet(types.AddFiatPegToWallet(order.GetFiatPegWallet(), fiatPegWallet))
	keeper.SetOrder(ctx, order)
	return nil
}

// GetOrderDetails : get the order details
func (keeper Keeper) GetOrderDetails(ctx cTypes.Context, buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash) (cTypes.Error, types.AssetPegWallet, types.FiatPegWallet, string, string) {
	negotiationID := negotiation.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	order := keeper.GetOrder(ctx, negotiationID)
	if order == nil {
		return cTypes.ErrInvalidAddress("Order not found!"), nil, nil, "", ""
	}
	return nil, order.GetAssetPegWallet(), order.GetFiatPegWallet(), order.GetFiatProofHash(), order.GetAWBProofHash()
}

// SetOrderFiatProofHash : Set FiatProofHash to Order
func (keeper Keeper) SetOrderFiatProofHash(ctx cTypes.Context, buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash, fiatProofHash string) {
	negotiationID := negotiation.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	order := keeper.GetOrder(ctx, negotiationID)
	order.SetFiatProofHash(fiatProofHash)
	keeper.SetOrder(ctx, order)
}

// SetOrderAWBProofHash : Set AWBProofHash to Order
func (keeper Keeper) SetOrderAWBProofHash(ctx cTypes.Context, buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash, awbProofHash string) {
	negotiationID := negotiation.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	order := keeper.GetOrder(ctx, negotiationID)
	order.SetAWBProofHash(awbProofHash)
	keeper.SetOrder(ctx, order)
}

// SendAssetFromOrder asset peg to buyer
func (keeper Keeper) SendAssetFromOrder(ctx cTypes.Context, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress, assetPeg types.AssetPeg) types.AssetPegWallet {
	negotiationID := negotiation.NegotiationID(append(append(fromAddress.Bytes(), toAddress.Bytes()...), assetPeg.GetPegHash().Bytes()...))
	order := keeper.GetOrder(ctx, negotiationID)
	_, updatedAssetPegWallet := types.SubtractAssetPegFromWallet(assetPeg.GetPegHash(), order.GetAssetPegWallet())
	order.SetAssetPegWallet(updatedAssetPegWallet)
	keeper.SetOrder(ctx, order)
	return updatedAssetPegWallet
}

// SendFiatsFromOrder fiat pegs to seller
func (keeper Keeper) SendFiatsFromOrder(ctx cTypes.Context, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress, pegHash types.PegHash, fiatPegWallet types.FiatPegWallet) types.FiatPegWallet {
	negotiationID := negotiation.NegotiationID(append(append(fromAddress.Bytes(), toAddress.Bytes()...), pegHash.Bytes()...))
	order := keeper.GetOrder(ctx, negotiationID)
	updatedFiatPegWallet := types.SubtractFiatPegWalletFromWallet(fiatPegWallet, order.GetFiatPegWallet())
	order.SetFiatPegWallet(updatedFiatPegWallet)
	keeper.SetOrder(ctx, order)
	return updatedFiatPegWallet
}
