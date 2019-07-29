package fiatFactory

import (
	"fmt"

	sdk "github.com/commitHub/commitBlockchain/types"
)

//Keeper : fiat keeper
type Keeper struct {
	fm FiatPegMapper
}

//NewKeeper : return a new keeper
func NewKeeper(fm FiatPegMapper) Keeper {
	return Keeper{fm: fm}
}

//*****comdex

func instantiateAndAssignFiat(ctx sdk.Context, fm FiatPegMapper, issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, newFiat sdk.FiatPeg) (sdk.FiatPeg, sdk.Tags, sdk.Error) {
	fiat := fm.GetFiatPeg(ctx, newFiat.GetPegHash())
	if fiat == nil {
		return fiat, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s", newFiat.GetPegHash()))
	}

	newFiat.SetPegHash(fiat.GetPegHash())

	var owner sdk.Owner
	owner.Amount = newFiat.GetTransactionAmount()
	owner.OwnerAddress = toAddress

	newFiat.SetOwners([]sdk.Owner{owner})
	fm.SetFiatPeg(ctx, newFiat)

	tags := sdk.NewTags("recepient", []byte(toAddress.String()))
	tags = tags.AppendTag("issuer", []byte(issuerAddress.String()))
	tags = tags.AppendTag("fiat", []byte(newFiat.GetPegHash().String()))

	return newFiat, tags, nil
}

func issueFiatsToWallets(ctx sdk.Context, fm FiatPegMapper, issueFiats []IssueFiat) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()

	for _, req := range issueFiats {
		_, tags, err := instantiateAndAssignFiat(ctx, fm, req.IssuerAddress, req.ToAddress, req.FiatPeg)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

//IssueFiatsToWallets haddles a list of IssueFiat messages
func (keeper Keeper) IssueFiatsToWallets(ctx sdk.Context, issueFiats []IssueFiat) (sdk.Tags, sdk.Error) {
	return issueFiatsToWallets(ctx, keeper.fm, issueFiats)
}

//##### Redeem Fiat

func redeemFiatFromWallet(ctx sdk.Context, fm FiatPegMapper, redeemerAddress sdk.AccAddress, amount int64, fiatPegWallet sdk.FiatPegWallet) (sdk.Tags, sdk.Error) {
	oldFiatPegWallet := sdk.FiatPegWallet{}
	for _, fiatPeg := range fiatPegWallet {
		fiat := fm.GetFiatPeg(ctx, fiatPeg.GetPegHash())
		if fiat == nil {
			return nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s", fiat.GetPegHash()))
		}
		oldFiatPegWallet = append(oldFiatPegWallet, sdk.ToBaseFiatPeg(fiat))
	}
	newFiatPegWallet := sdk.RedeemFiatPegsFromWallet(fiatPegWallet, oldFiatPegWallet, redeemerAddress)
	if newFiatPegWallet == nil {
		return nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s", newFiatPegWallet[0].GetPegHash()))
	}
	for _, fiatPeg := range newFiatPegWallet {
		fm.SetFiatPeg(ctx, &fiatPeg)
	}

	tags := sdk.NewTags("redeemer", []byte(redeemerAddress.String()))
	return tags, nil
}

func redeemFiatsFromWallets(ctx sdk.Context, fm FiatPegMapper, redeemFiats []RedeemFiat) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()

	for _, req := range redeemFiats {
		tags, err := redeemFiatFromWallet(ctx, fm, req.RedeemerAddress, req.Amount, req.FiatPegWallet)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

//RedeemFiatsFromWallets handles a list of Redeem Fiat messages
func (keeper Keeper) RedeemFiatsFromWallets(ctx sdk.Context, redeemFiats []RedeemFiat) (sdk.Tags, sdk.Error) {
	return redeemFiatsFromWallets(ctx, keeper.fm, redeemFiats)
}

//***** Send Fiats
func sendFiats(fiatPegWallet sdk.FiatPegWallet, fm FiatPegMapper, ctx sdk.Context, pegHash sdk.PegHash, fromAddress sdk.AccAddress, toAddress sdk.AccAddress) (sdk.Tags, sdk.Error) {
	oldFiatPegWallet := sdk.FiatPegWallet{}
	for _, fiatPeg := range fiatPegWallet {
		fiat := fm.GetFiatPeg(ctx, fiatPeg.GetPegHash())
		if fiat == nil {
			return nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s", fiatPeg.GetPegHash()))
		}
		oldFiatPegWallet = append(oldFiatPegWallet, sdk.ToBaseFiatPeg(fiat))
	}
	newFiatPegWallet := sdk.TransferFiatPegsToWallet(fiatPegWallet, oldFiatPegWallet, fromAddress, toAddress)
	if newFiatPegWallet == nil {
		return nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s", pegHash))
	}
	for _, fiatPeg := range newFiatPegWallet {
		fm.SetFiatPeg(ctx, &fiatPeg)
	}

	tags := sdk.NewTags("recepient", []byte(toAddress.String()))
	tags = tags.AppendTag("sender", []byte(fromAddress.String()))
	tags = tags.AppendTag("asset", []byte(pegHash.String()))
	return tags, nil
}

func sendFiatToOrder(ctx sdk.Context, fm FiatPegMapper, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, fiatPegWallet sdk.FiatPegWallet) (sdk.Tags, sdk.Error) {
	return sendFiats(fiatPegWallet, fm, ctx, pegHash, fromAddress, sdk.AccAddress(sdk.GenerateNegotiationIDBytes(fromAddress, toAddress, pegHash)))
}

func sendFiatsToOrders(ctx sdk.Context, fm FiatPegMapper, sendFiats []SendFiat) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()

	for _, req := range sendFiats {
		tags, err := sendFiatToOrder(ctx, fm, req.FromAddress, req.ToAddress, req.PegHash, req.FiatPegWallet)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

//SendFiatsToOrders handles a list of SendFiat messages
func (keeper Keeper) SendFiatsToOrders(ctx sdk.Context, sendFiats []SendFiat) (sdk.Tags, sdk.Error) {
	return sendFiatsToOrders(ctx, keeper.fm, sendFiats)
}

//##### Send Fiats
//***** Execute Fiats
func sendFiatFromOrder(ctx sdk.Context, fm FiatPegMapper, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, fiatPegWallet sdk.FiatPegWallet) (sdk.Tags, sdk.Error) {
	tags, err := sendFiats(fiatPegWallet, fm, ctx, pegHash, sdk.AccAddress(sdk.GenerateNegotiationIDBytes(fromAddress, toAddress, pegHash)), toAddress)
	for _, fiatPeg := range fiatPegWallet {
		fiat := fm.GetFiatPeg(ctx, fiatPeg.GetPegHash())
		oldFiatPeg := sdk.ToBaseFiatPeg(fiat)
		for _, owner := range oldFiatPeg.Owners {
			if owner.OwnerAddress.String() == sdk.AccAddress(sdk.GenerateNegotiationIDBytes(fromAddress, toAddress, pegHash)).String() && owner.Amount != 0 {
				fiat.SetTransactionAmount(owner.Amount)
				fiatPegWallet = sdk.FiatPegWallet{sdk.ToBaseFiatPeg(fiat)}
				sendFiats(fiatPegWallet, fm, ctx, pegHash, sdk.AccAddress(sdk.GenerateNegotiationIDBytes(fromAddress, toAddress, pegHash)), fromAddress)
			}
		}
	}
	return tags, err
}

func executeFiatOrders(ctx sdk.Context, fm FiatPegMapper, executeFiats []SendFiat) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()

	for _, req := range executeFiats {
		tags, err := sendFiatFromOrder(ctx, fm, req.FromAddress, req.ToAddress, req.PegHash, req.FiatPegWallet)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

//ExecuteFiatOrders handles a list of ExecuteFiat messages
func (keeper Keeper) ExecuteFiatOrders(ctx sdk.Context, executeFiats []SendFiat) (sdk.Tags, sdk.Error) {
	return executeFiatOrders(ctx, keeper.fm, executeFiats)
}

//##### Execute Fiats
//#####comdex
