package keeper

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceSDK/codec"
	fiatFactoryTypes "github.com/persistenceOne/persistenceSDK/modules/fiatFactory/internal/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

type Keeper struct {
	storeKey      cTypes.StoreKey
	accountKeeper fiatFactoryTypes.AccountKeeper
	cdc           *codec.Codec
}

func NewKeeper(cdc *codec.Codec, storeKey cTypes.StoreKey, accountKeeper fiatFactoryTypes.AccountKeeper) Keeper {
	return Keeper{
		storeKey:      storeKey,
		accountKeeper: accountKeeper,
		cdc:           cdc,
	}
}

func (k Keeper) SetFiatPeg(ctx cTypes.Context, fiatPeg types.FiatPeg) {
	store := ctx.KVStore(k.storeKey)
	fiatPegKey := fiatFactoryTypes.FiatPegHashStoreKey(fiatPeg.GetPegHash())
	bytes := k.cdc.MustMarshalBinaryLengthPrefixed(fiatPeg)
	store.Set(fiatPegKey, bytes)
}

func (k Keeper) GetFiatPeg(ctx cTypes.Context, pegHash types.PegHash) (types.FiatPeg, cTypes.Error) {
	store := ctx.KVStore(k.storeKey)
	fiatKey := fiatFactoryTypes.FiatPegHashStoreKey(pegHash)
	data := store.Get(fiatKey)
	if data == nil {
		return nil, fiatFactoryTypes.ErrInvalidString(fiatFactoryTypes.DefaultCodeSpace, fmt.Sprintf("Fiat with pegHash %s not found", pegHash))
	}
	var fiatPeg types.FiatPeg
	k.cdc.MustUnmarshalBinaryLengthPrefixed(data, &fiatPeg)
	return fiatPeg, nil
}

func (k Keeper) IterateFiats(ctx cTypes.Context, handler func(fiatPeg types.FiatPeg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := cTypes.KVStorePrefixIterator(store, fiatFactoryTypes.PegHashKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var fiatPeg types.FiatPeg
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &fiatPeg)
		if handler(fiatPeg) {
			break
		}
	}
}

func (k Keeper) GetFiatPegs(ctx cTypes.Context) (fiatPegs []types.FiatPeg) {
	k.IterateFiats(ctx, func(fiatPeg types.FiatPeg) (stop bool) {
		fiatPegs = append(fiatPegs, fiatPeg)
		return false
	},
	)
	return fiatPegs
}

func (keeper Keeper) InstantiateAndAssignFiat(ctx cTypes.Context, issuerAddress cTypes.AccAddress,
	toAddress cTypes.AccAddress, newFiat types.FiatPeg) cTypes.Error {

	fiat, err := keeper.GetFiatPeg(ctx, newFiat.GetPegHash())
	if err != nil {
		return err
	}

	_err := newFiat.SetPegHash(fiat.GetPegHash())
	if _err != nil {
		return fiatFactoryTypes.ErrInvalidPegHash(fiatFactoryTypes.DefaultCodeSpace)
	}

	var owner types.Owner
	owner.Amount = newFiat.GetTransactionAmount()
	owner.OwnerAddress = toAddress

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			fiatFactoryTypes.EventTypeFiatFactoryAssignFiat,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("issuer", issuerAddress.String()),
			cTypes.NewAttribute("fiat", newFiat.GetPegHash().String()),
		))

	return nil
}

func (keeper Keeper) InstantiateAndRedeemFiat(ctx cTypes.Context, redeemerAddress cTypes.AccAddress,
	fiatPegWallet types.FiatPegWallet) cTypes.Error {

	oldFiatPegWallet := types.FiatPegWallet{}
	for _, fiatPeg := range fiatPegWallet {
		fiat, err := keeper.GetFiatPeg(ctx, fiatPeg.GetPegHash())
		if err != nil {
			return fiatFactoryTypes.ErrInvalidPegHash(fiatFactoryTypes.DefaultCodeSpace)
		}

		if fiat == nil {
			return cTypes.ErrInsufficientCoins(fmt.Sprintf("%s", fiat.GetPegHash()))
		}

		oldFiatPegWallet = append(oldFiatPegWallet, types.ToBaseFiatPeg(fiat))
	}

	newFiatPegWallet := types.RedeemFiatPegsFromWallet(fiatPegWallet, oldFiatPegWallet, redeemerAddress)
	if newFiatPegWallet == nil {
		return cTypes.ErrInsufficientCoins(fmt.Sprintf("%s", newFiatPegWallet[0].GetPegHash()))
	}

	for _, fiatPeg := range newFiatPegWallet {
		keeper.SetFiatPeg(ctx, &fiatPeg)
	}

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			fiatFactoryTypes.EventTypeFiatFactoryRedeemFiat,
			cTypes.NewAttribute("redeemer", redeemerAddress.String()),
		))

	return nil
}

func (keeper Keeper) SendFiatToOrder(ctx cTypes.Context, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress,
	pegHash types.PegHash, fiatPegWallet types.FiatPegWallet) cTypes.Error {

	oldFiatPegWallet := types.FiatPegWallet{}

	for _, fiatPeg := range fiatPegWallet {
		fiat, err := keeper.GetFiatPeg(ctx, fiatPeg.GetPegHash())
		if err != nil {
			return fiatFactoryTypes.ErrInvalidPegHash(fiatFactoryTypes.DefaultCodeSpace)
		}

		if fiat == nil {
			return cTypes.ErrInsufficientCoins(fmt.Sprintf("%s", fiatPeg.GetPegHash()))
		}

		oldFiatPegWallet = append(oldFiatPegWallet, types.ToBaseFiatPeg(fiat))
	}

	newFiatPegWallet := types.TransferFiatPegsToWallet(fiatPegWallet, oldFiatPegWallet, fromAddress, toAddress)
	if newFiatPegWallet == nil {
		return cTypes.ErrInsufficientCoins(fmt.Sprintf("%s", pegHash))
	}

	for _, fiatPeg := range newFiatPegWallet {
		keeper.SetFiatPeg(ctx, &fiatPeg)
	}

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			fiatFactoryTypes.EventTypeFiatFactorySendFiat,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("sender", fromAddress.String()),
			cTypes.NewAttribute("fiat", pegHash.String()),
		))

	return nil
}

func (keeper Keeper) SendFiatFromOrder(ctx cTypes.Context, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress,
	pegHash types.PegHash, fiatPegWallet types.FiatPegWallet) cTypes.Error {

	negotiationId := types.NegotiationID(append(append(fromAddress.Bytes(), toAddress.Bytes()...), pegHash.Bytes()...))
	keeper.SendFiats(fiatPegWallet, ctx, pegHash, cTypes.AccAddress(negotiationId), toAddress)

	for _, fiatPeg := range fiatPegWallet {
		fiat, err := keeper.GetFiatPeg(ctx, fiatPeg.GetPegHash())
		if err != nil {
			return cTypes.ErrInternal("No fiat found")
		}

		oldFiatPeg := types.ToBaseFiatPeg(fiat)
		for _, owner := range oldFiatPeg.Owners {
			negotiationId := types.NegotiationID(append(append(fromAddress.Bytes(), toAddress.Bytes()...), pegHash.Bytes()...))
			if owner.OwnerAddress.String() == cTypes.AccAddress(negotiationId).String() && owner.Amount != 0 {
				fiat.SetTransactionAmount(owner.Amount)
				fiatPegWallet = types.FiatPegWallet{types.ToBaseFiatPeg(fiat)}
				keeper.SendFiats(fiatPegWallet, ctx, pegHash, cTypes.AccAddress(negotiationId), fromAddress)
			}
		}
	}

	return nil
}

func (keeper Keeper) SendFiats(fiatPegWallet types.FiatPegWallet, ctx cTypes.Context,
	pegHash types.PegHash, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress) cTypes.Error {

	oldFiatPegWallet := types.FiatPegWallet{}
	for _, fiatPeg := range fiatPegWallet {
		fiat, err := keeper.GetFiatPeg(ctx, fiatPeg.GetPegHash())
		if err != nil {
			return fiatFactoryTypes.ErrInvalidPegHash(fiatFactoryTypes.DefaultCodeSpace)
		}
		if fiat == nil {
			return cTypes.ErrInsufficientCoins(fmt.Sprintf("%s", fiatPeg.GetPegHash()))
		}
		oldFiatPegWallet = append(oldFiatPegWallet, types.ToBaseFiatPeg(fiat))
	}

	newFiatPegWallet := types.TransferFiatPegsToWallet(fiatPegWallet, oldFiatPegWallet, fromAddress, toAddress)
	if newFiatPegWallet == nil {
		return cTypes.ErrInsufficientCoins(fmt.Sprintf("%s", pegHash))
	}

	for _, fiatPeg := range newFiatPegWallet {
		keeper.SetFiatPeg(ctx, &fiatPeg)
	}

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			fiatFactoryTypes.EventTypeFiatFactoryExecuteFiat,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("sender", fromAddress.String()),
			cTypes.NewAttribute("fiat", pegHash.String()),
		))

	return nil
}
