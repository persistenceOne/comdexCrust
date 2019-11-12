package keeper

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/auth"
	negotiationTypes "github.com/commitHub/commitBlockchain/modules/negotiation/internal/types"
	"github.com/commitHub/commitBlockchain/modules/reputation"
	"github.com/commitHub/commitBlockchain/types"
)

type Keeper struct {
	storeKey      cTypes.StoreKey
	accountKeeper auth.AccountKeeper
	cdc           *codec.Codec
	aclKeeper     acl.Keeper
	rk            reputation.Keeper
}

func NewKeeper(storeKey cTypes.StoreKey, ak auth.AccountKeeper, aclKeeper acl.Keeper, rk reputation.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:      storeKey,
		accountKeeper: ak,
		aclKeeper:     aclKeeper,
		rk:            rk,
		cdc:           cdc,
	}
}

// negotiation/{0x01}/{buyerAddress+sellerAddress+pegHash} => negotiation
func (k Keeper) SetNegotiation(ctx cTypes.Context, negotiation types.Negotiation) {
	store := ctx.KVStore(k.storeKey)

	negotiationKey := negotiationTypes.GetNegotiationKey(negotiation.GetNegotiationID())
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(negotiation)
	store.Set(negotiationKey, bz)
}

// returns negotiation by negotiationID
func (k Keeper) GetNegotiation(ctx cTypes.Context, negotiationID types.NegotiationID) (negotiation types.Negotiation, err cTypes.Error) {
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
func (k Keeper) GetNegotiations(ctx cTypes.Context) (negotiations []types.Negotiation) {
	k.IterateNegotiations(ctx, func(negotiation types.Negotiation) (stop bool) {
		negotiations = append(negotiations, negotiation)
		return false
	},
	)
	return
}

func (k Keeper) IterateNegotiations(ctx cTypes.Context, handler func(negotiation types.Negotiation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := cTypes.KVStorePrefixIterator(store, negotiationTypes.NegotiationKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var negotiation types.Negotiation
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
	hash types.PegHash) (types.Negotiation, cTypes.Error) {

	negotiationID := types.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), hash.Bytes()...))
	_negotiation, err := k.GetNegotiation(ctx, negotiationID)
	if err != nil {
		return nil, err
	}
	return _negotiation, nil
}

func (k Keeper) ChangeNegotiationBidWithACL(ctx cTypes.Context, changeBid negotiationTypes.ChangeBid) cTypes.Error {
	aclStoreBuyer, err := k.aclKeeper.GetAccountACLDetails(ctx, changeBid.Negotiation.GetBuyerAddress())
	if err != nil {
		return cTypes.ErrInternal("Unauthorized transaction")
	}
	aclStoreBuyer.GetACL()
	accountBuyer := aclStoreBuyer.GetACL()
	aclStoreSeller, err := k.aclKeeper.GetAccountACLDetails(ctx, changeBid.Negotiation.GetSellerAddress())
	if err != nil {
		return cTypes.ErrInternal("Unauthorized transaction")
	}
	aclStoreSeller.GetACL()
	accountSeller := aclStoreSeller.GetACL()
	if !accountBuyer.Negotiation || !accountSeller.Negotiation {
		return cTypes.ErrInternal("Unauthorized transaction")
	}

	err = CheckTakerAddress(ctx, k, changeBid.Negotiation.GetSellerAddress(), changeBid.Negotiation.GetBuyerAddress(), changeBid.Negotiation.GetPegHash())
	if err != nil {
		return err
	}
	err = createOrChangeNegotiationBid(ctx, k, changeBid.Negotiation)
	if err != nil {
		return err
	}

	k.rk.SetChangeBuyerBidPositiveTx(ctx, changeBid.Negotiation.GetBuyerAddress())
	k.rk.SetChangeSellerBidPositiveTx(ctx, changeBid.Negotiation.GetSellerAddress())
	return nil
}

func createOrChangeNegotiationBid(ctx cTypes.Context, negotiationKeeper Keeper, negotiation types.Negotiation) cTypes.Error {

	oldNegotiation, _ := negotiationKeeper.GetNegotiation(ctx, negotiation.GetNegotiationID())

	if oldNegotiation == nil {
		oldNegotiation = types.NewNegotiation(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash())
	}

	if oldNegotiation.GetBuyerSignature() != nil || oldNegotiation.GetSellerSignature() != nil {
		return negotiationTypes.ErrVerifySignature(negotiationTypes.DefaultCodeSpace, "Already signed. Cannot change negotiation now")
	}

	oldNegotiation.SetBid(negotiation.GetBid())
	oldNegotiation.SetTime(negotiation.GetTime())

	negotiationKeeper.SetNegotiation(ctx, oldNegotiation)

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(negotiationTypes.EventTypeChangeNegotiationBid,
			cTypes.NewAttribute(negotiationTypes.AttributeKeyNegotiationID, oldNegotiation.GetNegotiationID().String()),
			cTypes.NewAttribute(negotiationTypes.AttributeKeyBuyerAddress, oldNegotiation.GetBuyerAddress().String()),
			cTypes.NewAttribute(negotiationTypes.AttributeKeySellerAddress, oldNegotiation.GetSellerAddress().String()),
			cTypes.NewAttribute(negotiationTypes.AttributeKeyPegHash, oldNegotiation.GetPegHash().String()),
		))

	return nil
}

// CheckTakerAddress :
func CheckTakerAddress(ctx cTypes.Context, k Keeper, sellerAddress cTypes.AccAddress, buyerAddress cTypes.AccAddress, pegHash types.PegHash) cTypes.Error {
	account := k.accountKeeper.GetAccount(ctx, sellerAddress)
	assetPegWallet := account.GetAssetPegWallet()
	assetPeg := types.BaseAssetPeg{}
	i := assetPegWallet.SearchAssetPeg(pegHash)
	if i < len(assetPegWallet) && assetPegWallet[i].GetPegHash().String() == pegHash.String() {
		assetPeg = assetPegWallet[i]
	} else {
		return nil
	}
	takerAddress := assetPeg.GetTakerAddress()
	if takerAddress != nil && takerAddress.String() != buyerAddress.String() {
		return cTypes.ErrInternal(fmt.Sprintf("Transaction is not permitted with %s", buyerAddress.String()))
	}
	return nil
}

func (k Keeper) ConfirmNegotiationBidWithACL(ctx cTypes.Context, confirmBid negotiationTypes.ConfirmBid) cTypes.Error {
	aclStoreBuyer, err := k.aclKeeper.GetAccountACLDetails(ctx, confirmBid.Negotiation.GetBuyerAddress())
	if err != nil {
		return cTypes.ErrInternal("Unauthorized transaction")
	}
	aclStoreBuyer.GetACL()
	accountBuyer := aclStoreBuyer.GetACL()
	aclStoreSeller, err := k.aclKeeper.GetAccountACLDetails(ctx, confirmBid.Negotiation.GetSellerAddress())
	if err != nil {
		return cTypes.ErrInternal("Unauthorized transaction")
	}
	aclStoreSeller.GetACL()
	accountSeller := aclStoreSeller.GetACL()
	if !accountBuyer.ConfirmBuyerBid || !accountSeller.ConfirmSellerBid {
		return cTypes.ErrInternal("Unauthorized transaction")
	}
	err = CheckTakerAddress(ctx, k, confirmBid.Negotiation.GetSellerAddress(), confirmBid.Negotiation.GetBuyerAddress(), confirmBid.Negotiation.GetPegHash())
	if err != nil {
		return err
	}
	err = confirmNegotiationBid(ctx, k, confirmBid.Negotiation)
	if err != nil {
		return err
	}

	k.rk.SetConfirmBuyerBidPositiveTx(ctx, confirmBid.Negotiation.GetBuyerAddress())
	k.rk.SetConfirmSellerBidPositiveTx(ctx, confirmBid.Negotiation.GetSellerAddress())

	return nil
}

func confirmNegotiationBid(ctx cTypes.Context, negotiationKeeper Keeper, negotiation types.Negotiation) cTypes.Error {
	oldNegotiation, _ := negotiationKeeper.GetNegotiation(ctx, negotiation.GetNegotiationID())
	if oldNegotiation == nil {
		oldNegotiation = types.NewNegotiation(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash())
		oldNegotiation.SetBid(negotiation.GetBid())
	}

	if oldNegotiation.GetSellerSignature() != nil && oldNegotiation.GetBuyerSignature() != nil {
		return negotiationTypes.ErrVerifySignature(negotiationTypes.DefaultCodeSpace, "Already Exist the signatures")
	}

	if oldNegotiation.GetBid() != negotiation.GetBid() {
		return negotiationTypes.ErrInvalidBid(negotiationTypes.DefaultCodeSpace, "Buyer and Seller must confirm with same bid amount")
	}

	oldNegotiation.SetTime(negotiation.GetTime())

	if negotiation.GetSellerSignature() != nil {
		oldNegotiation.SetSellerBlockHeight(ctx.BlockHeight())
		oldNegotiation.SetSellerContractHash(negotiation.GetSellerContractHash())
		signBytes := negotiationTypes.NewSignNegotiationBody(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash(), negotiation.GetBid(), negotiation.GetTime()).GetSignBytes()
		account := negotiationKeeper.GetNegotiatorAccount(ctx, negotiation.GetSellerAddress())

		if !VerifySignature(account.GetPubKey(), negotiation.GetSellerSignature(), signBytes) {
			return negotiationTypes.ErrVerifySignature(negotiationTypes.DefaultCodeSpace, "Seller signature verification failed")
		}
		oldNegotiation.SetSellerSignature(negotiation.GetSellerSignature())
	}

	if negotiation.GetBuyerSignature() != nil {
		oldNegotiation.SetBuyerBlockHeight(ctx.BlockHeight())
		oldNegotiation.SetBuyerContractHash(negotiation.GetBuyerContractHash())
		signBytes := negotiationTypes.NewSignNegotiationBody(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash(), negotiation.GetBid(), negotiation.GetTime()).GetSignBytes()
		account := negotiationKeeper.GetNegotiatorAccount(ctx, negotiation.GetBuyerAddress())

		if !VerifySignature(account.GetPubKey(), negotiation.GetBuyerSignature(), signBytes) {
			return negotiationTypes.ErrVerifySignature(negotiationTypes.DefaultCodeSpace, "Seller signature verification failed")
		}
		oldNegotiation.SetBuyerSignature(negotiation.GetBuyerSignature())
	}

	negotiationKeeper.SetNegotiation(ctx, oldNegotiation)

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			negotiationTypes.EventTypeConfirmNegotiationBid,
			cTypes.NewAttribute(negotiationTypes.AttributeKeyNegotiationID, oldNegotiation.GetNegotiationID().String()),
			cTypes.NewAttribute(negotiationTypes.AttributeKeyBuyerAddress, oldNegotiation.GetBuyerAddress().String()),
			cTypes.NewAttribute(negotiationTypes.AttributeKeySellerAddress, oldNegotiation.GetSellerAddress().String()),
			cTypes.NewAttribute(negotiationTypes.AttributeKeyPegHash, oldNegotiation.GetPegHash().String()),
		))

	return nil
}
func VerifySignature(pubKey crypto.PubKey, signature types.Signature, signBytes []byte) bool {
	if !pubKey.VerifyBytes(signBytes, signature) {
		return false
	}

	return true
}
