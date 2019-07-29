package order

import (
	sdk "github.com/commitHub/commitBlockchain/types"
)

//Keeper : asset keeper
type Keeper struct {
	om Mapper
}

//NewKeeper : return a new keeper
func NewKeeper(om Mapper) Keeper {
	return Keeper{om: om}
}

//SendAssetsToOrder fiat pegs to order
func (keeper Keeper) SendAssetsToOrder(ctx sdk.Context, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, assetPeg sdk.AssetPeg) sdk.Error {
	negotiationID := sdk.NegotiationID(append(append(toAddress.Bytes(), fromAddress.Bytes()...), assetPeg.GetPegHash().Bytes()...))
	order := keeper.om.GetOrder(ctx, negotiationID)
	if order == nil {
		order = keeper.om.NewOrder(toAddress, fromAddress, assetPeg.GetPegHash())
	}
	order.SetAssetPegWallet(sdk.AddAssetPegToWallet(assetPeg, order.GetAssetPegWallet()))
	keeper.om.SetOrder(ctx, order)
	return nil
}

//SendFiatsToOrder fiat pegs to order
func (keeper Keeper) SendFiatsToOrder(ctx sdk.Context, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, fiatPegWallet sdk.FiatPegWallet) sdk.Error {
	negotiationID := sdk.NegotiationID(append(append(fromAddress.Bytes(), toAddress.Bytes()...), pegHash.Bytes()...))
	order := keeper.om.GetOrder(ctx, negotiationID)
	if order == nil {
		order = keeper.om.NewOrder(fromAddress, toAddress, pegHash)
	}
	order.SetFiatPegWallet(sdk.AddFiatPegToWallet(order.GetFiatPegWallet(), fiatPegWallet))
	keeper.om.SetOrder(ctx, order)
	return nil
}

//SendFiatsFromOrder fiat pegs to seller
func (keeper Keeper) SendFiatsFromOrder(ctx sdk.Context, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, fiatPegWallet sdk.FiatPegWallet) sdk.FiatPegWallet {
	negotiationID := sdk.NegotiationID(append(append(fromAddress.Bytes(), toAddress.Bytes()...), pegHash.Bytes()...))
	order := keeper.om.GetOrder(ctx, negotiationID)
	updatedFiatPegWallet := sdk.SubtractFiatPegWalletFromWallet(fiatPegWallet, order.GetFiatPegWallet())
	order.SetFiatPegWallet(updatedFiatPegWallet)
	keeper.om.SetOrder(ctx, order)
	return updatedFiatPegWallet
}

//SendAssetFromOrder asset peg to buyer
func (keeper Keeper) SendAssetFromOrder(ctx sdk.Context, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, assetPeg sdk.AssetPeg) sdk.AssetPegWallet {
	negotiationID := sdk.NegotiationID(append(append(fromAddress.Bytes(), toAddress.Bytes()...), assetPeg.GetPegHash().Bytes()...))
	order := keeper.om.GetOrder(ctx, negotiationID)
	_, updatedAssetPegWallet := sdk.SubtractAssetPegFromWallet(assetPeg.GetPegHash(), order.GetAssetPegWallet())
	order.SetAssetPegWallet(updatedAssetPegWallet)
	keeper.om.SetOrder(ctx, order)
	return updatedAssetPegWallet
}

//GetOrderDetails : get the order details
func (keeper Keeper) GetOrderDetails(ctx sdk.Context, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash) (sdk.Error, sdk.AssetPegWallet, sdk.FiatPegWallet, string, string) {
	negotiationID := sdk.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	order := keeper.om.GetOrder(ctx, negotiationID)
	if order == nil {
		return sdk.ErrInvalidAddress("Order not found!"), nil, nil, "", ""
	}
	return nil, order.GetAssetPegWallet(), order.GetFiatPegWallet(), order.GetFiatProofHash(), order.GetAWBProofHash()
}

//SetOrderFiatProofHash : Set FiatProofHash to Order
func (keeper Keeper) SetOrderFiatProofHash(ctx sdk.Context, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, fiatProofHash string) {
	negotiationID := sdk.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	order := keeper.om.GetOrder(ctx, negotiationID)
	order.SetFiatProofHash(fiatProofHash)
	keeper.om.SetOrder(ctx, order)
}

//SetOrderAWBProofHash : Set AWBProofHash to Order
func (keeper Keeper) SetOrderAWBProofHash(ctx sdk.Context, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, awbProofHash string) {
	negotiationID := sdk.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	order := keeper.om.GetOrder(ctx, negotiationID)
	order.SetAWBProofHash(awbProofHash)
	keeper.om.SetOrder(ctx, order)
}
