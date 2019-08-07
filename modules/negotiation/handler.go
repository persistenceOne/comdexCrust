package negotiation

import (
	"fmt"
	
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

func NewHandler(k Keeper) cTypes.Handler { // TODO AclKeeper, ReputationKeeper
	return func(ctx cTypes.Context, msg cTypes.Msg) cTypes.Result {
		ctx = ctx.WithEventManager(cTypes.NewEventManager())
		
		switch msg := msg.(type) {
		case MsgChangeBuyerBids:
			return handleMsgChangeBids(ctx, k, msg.ChangeBids)
		case MsgChangeSellerBids:
			return handleMsgChangeBids(ctx, k, msg.ChangeBids)
		case MsgConfirmSellerBids:
			return handleMsgConfirmBids(ctx, k, msg.ConfirmBids)
		case MsgConfirmBuyerBids:
			return handleMsgConfirmBids(ctx, k, msg.ConfirmBids)
		
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return cTypes.ErrUnknownRequest(errMsg).Result()
			
		}
		
	}
}

func handleMsgChangeBids(ctx cTypes.Context, negotiationKeeper Keeper, changeBids []ChangeBid) cTypes.Result {
	
	for _, changeBid := range changeBids {
		err := changeNegotiationBidWithACL(ctx, negotiationKeeper, changeBid)
		if err != nil {
			return err.Result()
		}
	}
	
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func changeNegotiationBidWithACL(ctx cTypes.Context, negotiationKeeper Keeper, changeBid ChangeBid) cTypes.Error {
	// TODO  ACLImplementation
	err := createOrChangeNegotiationBid(ctx, negotiationKeeper, changeBid.Negotiation)
	if err != nil {
		return err
	}
	return nil
}

func createOrChangeNegotiationBid(ctx cTypes.Context, negotiationKeeper Keeper, negotiation Negotiation) cTypes.Error {
	
	oldNegotiation, _ := negotiationKeeper.GetNegotiation(ctx, negotiation.GetNegotiationID())
	
	if oldNegotiation == nil {
		oldNegotiation = NewNegotiation(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash())
	}
	
	if oldNegotiation.GetBuyerSignature() != nil || oldNegotiation.GetSellerSignature() != nil {
		return ErrCodeVerifySignature(DefaultCodeSpace, "Already signed. Cannot change negotiation now")
	}
	
	oldNegotiation.SetBid(negotiation.GetBid())
	oldNegotiation.SetTime(negotiation.GetTime())
	
	negotiationKeeper.SetNegotiation(ctx, oldNegotiation)
	
	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(EventTypeChangeNegotiationBid,
			cTypes.NewAttribute(AttributeKeyNegotiationID, oldNegotiation.GetNegotiationID().String()),
			cTypes.NewAttribute(AttributeKeyBuyerAddress, oldNegotiation.GetBuyerAddress().String()),
			cTypes.NewAttribute(AttributeKeySellerAddress, oldNegotiation.GetSellerAddress().String()),
			cTypes.NewAttribute(AttributeKeyPegHash, oldNegotiation.GetPegHash().String()),
		))
	
	return nil
}

func handleMsgConfirmBids(ctx cTypes.Context, negotitationKeeper Keeper, confirmBids []ConfirmBid) cTypes.Result {
	for _, confirmBid := range confirmBids {
		err := confirmNegotiationBidWithACL(ctx, negotitationKeeper, confirmBid)
		if err != nil {
			return err.Result()
		}
	}
	
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
	
}

func confirmNegotiationBidWithACL(ctx cTypes.Context, negotiationKeeper Keeper, confirmBid ConfirmBid) cTypes.Error {
	// TODO ACLImplementation
	
	err := confirmNegotiationBid(ctx, negotiationKeeper, confirmBid.Negotiation)
	if err != nil {
		return err
	}
	return nil
}

func confirmNegotiationBid(ctx cTypes.Context, negotiationKeeper Keeper, negotiation Negotiation) cTypes.Error {
	oldNegotiation, _ := negotiationKeeper.GetNegotiation(ctx, negotiation.GetNegotiationID())
	if oldNegotiation == nil {
		oldNegotiation = NewNegotiation(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash())
		oldNegotiation.SetBid(negotiation.GetBid())
	}
	
	if oldNegotiation.GetSellerSignature() != nil && oldNegotiation.GetBuyerSignature() != nil {
		return ErrCodeVerifySignature(DefaultCodeSpace, "Already Exist the signatures")
	}
	
	if oldNegotiation.GetBid() != negotiation.GetBid() {
		return ErrCodeInvalidBid(DefaultCodeSpace, "Buyer and Seller must confirm with same bid amount")
	}
	
	oldNegotiation.SetTime(negotiation.GetTime())
	
	if negotiation.GetSellerSignature() != nil {
		oldNegotiation.SetSellerBlockHeight(ctx.BlockHeight())
		oldNegotiation.SetSellerContractHash(negotiation.GetSellerContractHash())
		signBytes := NewSignNegotiationBody(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash(), negotiation.GetBid(), negotiation.GetTime()).GetSignBytes()
		account := negotiationKeeper.GetNegotiatorAccount(ctx, negotiation.GetSellerAddress())
		
		if !VerifySignature(account.GetPubKey(), negotiation.GetSellerSignature(), signBytes) {
			return ErrCodeVerifySignature(DefaultCodeSpace, "Seller signature verification failed")
		}
		oldNegotiation.SetSellerSignature(negotiation.GetSellerSignature())
	}
	
	if negotiation.GetBuyerSignature() != nil {
		oldNegotiation.SetBuyerBlockHeight(ctx.BlockHeight())
		oldNegotiation.SetBuyerContractHash(negotiation.GetBuyerContractHash())
		signBytes := NewSignNegotiationBody(negotiation.GetBuyerAddress(), negotiation.GetSellerAddress(), negotiation.GetPegHash(), negotiation.GetBid(), negotiation.GetTime()).GetSignBytes()
		account := negotiationKeeper.GetNegotiatorAccount(ctx, negotiation.GetBuyerAddress())
		
		if !VerifySignature(account.GetPubKey(), negotiation.GetBuyerSignature(), signBytes) {
			return ErrCodeVerifySignature(DefaultCodeSpace, "Seller signature verification failed")
		}
		oldNegotiation.SetBuyerSignature(negotiation.GetBuyerSignature())
	}
	
	negotiationKeeper.SetNegotiation(ctx, oldNegotiation)
	
	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			EventTypeConfirmNegotiationBid,
			cTypes.NewAttribute(AttributeKeyNegotiationID, oldNegotiation.GetNegotiationID().String()),
			cTypes.NewAttribute(AttributeKeyBuyerAddress, oldNegotiation.GetBuyerAddress().String()),
			cTypes.NewAttribute(AttributeKeySellerAddress, oldNegotiation.GetSellerAddress().String()),
			cTypes.NewAttribute(AttributeKeyPegHash, oldNegotiation.GetPegHash().String()),
		))
	
	return nil
}

func VerifySignature(pubKey crypto.PubKey, signature Signature, signBytes []byte) bool {
	if !pubKey.VerifyBytes(signBytes, signature) {
		return false
	}
	
	return true
}
