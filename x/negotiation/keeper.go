package negotiation

import (
	"fmt"
	
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/x/acl"
	"github.com/comdex-blockchain/x/auth"
	"github.com/comdex-blockchain/x/reputation"
	crypto2 "github.com/tendermint/tendermint/crypto"
)

// Keeper : asset keeper
type Keeper struct {
	nm Mapper
	am auth.AccountMapper
}

// NewKeeper : return a new keeper
func NewKeeper(nm Mapper, am auth.AccountMapper) Keeper {
	return Keeper{nm: nm, am: am}
}

// GetNegotiation fiat pegs to order
func (keeper Keeper) GetNegotiation(ctx sdk.Context, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash) (sdk.Error, sdk.Negotiation) {
	negotiationID := sdk.NegotiationID(append(append(fromAddress.Bytes(), toAddress.Bytes()...), pegHash.Bytes()...))
	negotiation := keeper.nm.GetNegotiation(ctx, negotiationID)
	if negotiation == nil {
		return sdk.ErrInvalidAddress("Negotiation not found!"), nil
	}
	return nil, negotiation
}

func createOrChangeNegotiationBid(ctx sdk.Context, nm Mapper, negotiation sdk.Negotiation) (sdk.Negotiation, sdk.Tags, sdk.Error) {
	oldNegotiation := nm.GetNegotiation(ctx, negotiation.GetNegotiationID())
	if oldNegotiation == nil {
		oldNegotiation = nm.NewNegotiation(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash())
	}
	
	oldNegotiation.SetBid(negotiation.GetBid())
	oldNegotiation.SetTime(negotiation.GetTime())
	nm.SetNegotiation(ctx, oldNegotiation)
	
	tags := sdk.NewTags("negotiation_id", []byte(oldNegotiation.GetNegotiationID().String()))
	tags = tags.AppendTag("buyer", []byte(negotiation.GetBuyerAddress().String()))
	tags = tags.AppendTag("seller", []byte(negotiation.GetSellerAddress().String()))
	tags = tags.AppendTag("asset", []byte(negotiation.GetPegHash().String()))
	return oldNegotiation, tags, nil
}

func changeNegotiationBids(ctx sdk.Context, nm Mapper, changeBids []ChangeBid, ak acl.Keeper, reputationKeeper reputation.Keeper, am auth.AccountMapper) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()
	for _, req := range changeBids {
		aclStoreBuyer, err := ak.GetAccountACLDetails(ctx, req.Negotiation.GetBuyerAddress())
		if err != nil {
			return nil, sdk.ErrInternal("Unauthorized transaction")
		}
		aclStoreBuyer.GetACL()
		accountBuyer := aclStoreBuyer.GetACL()
		aclStoreSeller, err := ak.GetAccountACLDetails(ctx, req.Negotiation.GetSellerAddress())
		if err != nil {
			return nil, sdk.ErrInternal("Unauthorized transaction")
		}
		aclStoreSeller.GetACL()
		accountSeller := aclStoreSeller.GetACL()
		if !accountBuyer.Negotiation || !accountSeller.Negotiation {
			return nil, sdk.ErrInternal("Unauthorized transaction")
		}
		
		err = CheckTakerAddress(ctx, am, req.Negotiation.GetSellerAddress(), req.Negotiation.GetBuyerAddress(), req.Negotiation.GetPegHash())
		if err != nil {
			return nil, err
		}
		_, tags, err := createOrChangeNegotiationBid(ctx, nm, req.Negotiation)
		if err != nil {
			return nil, err
		}
		
		allTags = allTags.AppendTags(tags)
		reputationKeeper.SetChangeBuyerBidPositiveTx(ctx, req.Negotiation.GetBuyerAddress())
		reputationKeeper.SetChangeSellerBidPositiveTx(ctx, req.Negotiation.GetSellerAddress())
	}
	return allTags, nil
}

// ChangeNegotiationBids haddles a list of ChangeBid messages
func (keeper Keeper) ChangeNegotiationBids(ctx sdk.Context, changeBids []ChangeBid, ak acl.Keeper, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error) {
	return changeNegotiationBids(ctx, keeper.nm, changeBids, ak, reputationKeeper, keeper.am)
}

func confirmNegotiationBid(ctx sdk.Context, nm Mapper, am auth.AccountMapper, negotiation sdk.Negotiation) (sdk.Negotiation, sdk.Tags, sdk.Error) {
	oldNegotiation := nm.GetNegotiation(ctx, negotiation.GetNegotiationID())
	if oldNegotiation == nil {
		oldNegotiation = nm.NewNegotiation(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash())
		oldNegotiation.SetBid(negotiation.GetBid())
	}
	if oldNegotiation.GetSellerSignature() != nil && oldNegotiation.GetBuyerSignature() != nil {
		return nil, nil, sdk.ErrUnauthorized("Already Exist the Signatures")
	}
	
	if oldNegotiation.GetBid() != negotiation.GetBid() {
		return nil, nil, sdk.ErrInternal("Buyer and Seller must confirm with same bid amount")
	}
	oldNegotiation.SetTime(negotiation.GetTime())
	if negotiation.GetSellerSignature() != nil {
		oldNegotiation.SetSellerBlockHeight(ctx.BlockHeight())
		SignBytes := NewSignNegotiationBody(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash(), negotiation.GetBid(), negotiation.GetTime()).GetSignBytes()
		account := am.GetAccount(ctx, negotiation.GetSellerAddress())
		if !VerifySignature(account.GetPubKey(), negotiation.GetSellerSignature(), SignBytes) {
			return nil, nil, sdk.ErrInternal("Seller signature verification failed")
		}
		oldNegotiation.SetSellerSignature(negotiation.GetSellerSignature())
	}
	
	if negotiation.GetBuyerSignature() != nil {
		oldNegotiation.SetBuyerBlockHeight(ctx.BlockHeight())
		SignBytes := NewSignNegotiationBody(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash(), negotiation.GetBid(), negotiation.GetTime()).GetSignBytes()
		account := am.GetAccount(ctx, negotiation.GetBuyerAddress())
		if !VerifySignature(account.GetPubKey(), negotiation.GetBuyerSignature(), SignBytes) {
			return nil, nil, sdk.ErrInternal("Buyer Signature verification failed")
		}
		oldNegotiation.SetBuyerSignature(negotiation.GetBuyerSignature())
	}
	
	nm.SetNegotiation(ctx, oldNegotiation)
	tags := sdk.NewTags("negotiation_id", []byte(oldNegotiation.GetNegotiationID().String()))
	tags = tags.AppendTag("buyer", []byte(negotiation.GetBuyerAddress().String()))
	tags = tags.AppendTag("seller", []byte(negotiation.GetSellerAddress().String()))
	tags = tags.AppendTag("asset", []byte(negotiation.GetPegHash().String()))
	return oldNegotiation, tags, nil
}

func confirmNegotiationBids(ctx sdk.Context, nm Mapper, am auth.AccountMapper, confirmBids []ConfirmBid, ak acl.Keeper, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()
	for _, req := range confirmBids {
		aclStoreBuyer, err := ak.GetAccountACLDetails(ctx, req.Negotiation.GetBuyerAddress())
		if err != nil {
			return nil, sdk.ErrInternal("Unauthorized transaction")
		}
		aclStoreBuyer.GetACL()
		accountBuyer := aclStoreBuyer.GetACL()
		aclStoreSeller, err := ak.GetAccountACLDetails(ctx, req.Negotiation.GetSellerAddress())
		if err != nil {
			return nil, sdk.ErrInternal("Unauthorized transaction")
		}
		aclStoreSeller.GetACL()
		accountSeller := aclStoreSeller.GetACL()
		if !accountBuyer.ConfirmBuyerBid || !accountSeller.ConfirmSellerBid {
			return nil, sdk.ErrInternal("Unauthorized transaction")
		}
		err = CheckTakerAddress(ctx, am, req.Negotiation.GetSellerAddress(), req.Negotiation.GetBuyerAddress(), req.Negotiation.GetPegHash())
		if err != nil {
			return nil, err
		}
		_, tags, err := confirmNegotiationBid(ctx, nm, am, req.Negotiation)
		if err != nil {
			return nil, err
		}
		
		allTags = allTags.AppendTags(tags)
		reputationKeeper.SetConfirmBuyerBidPositiveTx(ctx, req.Negotiation.GetBuyerAddress())
		reputationKeeper.SetConfirmSellerBidPositiveTx(ctx, req.Negotiation.GetSellerAddress())
		
	}
	return allTags, nil
}

// ConfirmNegotiationBids haddles a list of ChangeBid messages
func (keeper Keeper) ConfirmNegotiationBids(ctx sdk.Context, confirmBids []ConfirmBid, ak acl.Keeper, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error) {
	return confirmNegotiationBids(ctx, keeper.nm, keeper.am, confirmBids, ak, reputationKeeper)
}

// #####comdex

// VerifySignature : vrifies the signature
func VerifySignature(pubkey crypto2.PubKey, signature sdk.Signature, signBytes []byte) bool {
	if !pubkey.VerifyBytes(signBytes, signature) {
		return false
	}
	return true
}

// CheckTakerAddress :
func CheckTakerAddress(ctx sdk.Context, am auth.AccountMapper, sellerAddress sdk.AccAddress, buyerAddress sdk.AccAddress, pegHash sdk.PegHash) sdk.Error {
	account := am.GetAccount(ctx, sellerAddress)
	assetPegWallet := account.GetAssetPegWallet()
	assetPeg := sdk.BaseAssetPeg{}
	i := assetPegWallet.SearchAssetPeg(pegHash)
	if i < len(assetPegWallet) && assetPegWallet[i].GetPegHash().String() == pegHash.String() {
		assetPeg = assetPegWallet[i]
	} else {
		return nil
	}
	takerAddress := assetPeg.GetTakerAddress()
	if takerAddress != nil && takerAddress.String() != buyerAddress.String() {
		return sdk.ErrInternal(fmt.Sprintf("Transaction is not permitted with %s", buyerAddress.String()))
	}
	return nil
}
