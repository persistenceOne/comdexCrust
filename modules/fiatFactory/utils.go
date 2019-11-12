package fiatFactory

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/modules/fiatFactory/internal/keeper"
	"github.com/commitHub/commitBlockchain/types"
)

func instantiateAndAssignFiat(ctx cTypes.Context, keeper Keeper, issuerAddress cTypes.AccAddress,
	toAddress cTypes.AccAddress, newFiat types.FiatPeg) cTypes.Error {

	fiat, err := keeper.GetFiatPeg(ctx, newFiat.GetPegHash())
	if err != nil {
		return err
	}

	_err := newFiat.SetPegHash(fiat.GetPegHash())
	if _err != nil {
		return ErrInvalidPegHash(DefaultCodeSpace)
	}

	var owner types.Owner
	owner.Amount = newFiat.GetTransactionAmount()
	owner.OwnerAddress = toAddress

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			EventTypeFiatFactoryAssignFiat,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("issuer", issuerAddress.String()),
			cTypes.NewAttribute("fiat", newFiat.GetPegHash().String()),
		))

	return nil
}

func instantiateAndRedeemFiat(ctx cTypes.Context, keeper Keeper, redeemerAddress cTypes.AccAddress,
	fiatPegWallet types.FiatPegWallet) cTypes.Error {

	oldFiatPegWallet := types.FiatPegWallet{}
	for _, fiatPeg := range fiatPegWallet {
		fiat, err := keeper.GetFiatPeg(ctx, fiatPeg.GetPegHash())
		if err != nil {
			return ErrInvalidPegHash(DefaultCodeSpace)
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
			EventTypeFiatFactoryRedeemFiat,
			cTypes.NewAttribute("redeemer", redeemerAddress.String()),
		))

	return nil
}

func sendFiatToOrder(ctx cTypes.Context, keeper Keeper, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress,
	pegHash types.PegHash, fiatPegWallet types.FiatPegWallet) cTypes.Error {

	oldFiatPegWallet := types.FiatPegWallet{}

	for _, fiatPeg := range fiatPegWallet {
		fiat, err := keeper.GetFiatPeg(ctx, fiatPeg.GetPegHash())
		if err != nil {
			return ErrInvalidPegHash(DefaultCodeSpace)
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
			EventTypeFiatFactorySendFiat,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("sender", fromAddress.String()),
			cTypes.NewAttribute("fiat", pegHash.String()),
		))

	return nil
}

func sendFiatFromOrder(ctx cTypes.Context, keeper Keeper, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress,
	pegHash types.PegHash, fiatPegWallet types.FiatPegWallet) cTypes.Error {

	negotiationId := types.NegotiationID(append(append(fromAddress.Bytes(), toAddress.Bytes()...), pegHash.Bytes()...))
	sendFiats(fiatPegWallet, keeper, ctx, pegHash, cTypes.AccAddress(negotiationId), toAddress)

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
				sendFiats(fiatPegWallet, keeper, ctx, pegHash, cTypes.AccAddress(negotiationId), fromAddress)
			}
		}
	}

	return nil
}

func sendFiats(fiatPegWallet types.FiatPegWallet, keeper keeper.Keeper, ctx cTypes.Context,
	pegHash types.PegHash, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress) cTypes.Error {

	oldFiatPegWallet := types.FiatPegWallet{}
	for _, fiatPeg := range fiatPegWallet {
		fiat, err := keeper.GetFiatPeg(ctx, fiatPeg.GetPegHash())
		if err != nil {
			return ErrInvalidPegHash(DefaultCodeSpace)
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
			EventTypeFiatFactoryExecuteFiat,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("sender", fromAddress.String()),
			cTypes.NewAttribute("fiat", pegHash.String()),
		))

	return nil
}
