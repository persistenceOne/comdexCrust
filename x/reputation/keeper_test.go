package reputation

import (
	"testing"

	sdk "github.com/commitHub/commitBlockchain/types"
)

var (
	sellerAddress = sdk.AccAddress("SellerAddress")
	buyerAddress  = sdk.AccAddress("BuyerAddress")
)

func Test_GetBaseReputationDetails(t *testing.T) {
	k.GetBaseReputationDetails(ctx, buyerAddress)
	k.SetSendAssetsPositiveTx(ctx, buyerAddress)
	k.SetSendAssetsNegativeTx(ctx, buyerAddress)
	k.SetSendFiatsPositiveTx(ctx, buyerAddress)
	k.SetSendFiatsNegativeTx(ctx, buyerAddress)
	k.SetIBCIssueAssetsPositiveTx(ctx, buyerAddress)
	k.SetIBCIssueAssetsNegativeTx(ctx, buyerAddress)
	k.SetIBCIssueFiatsPositiveTx(ctx, buyerAddress)
	k.SetIBCIssueFiatsNegativeTx(ctx, buyerAddress)
	k.SetBuyerExecuteOrderPositiveTx(ctx, buyerAddress)
	k.SetBuyerExecuteOrderNegativeTx(ctx, buyerAddress)
	k.SetSellerExecuteOrderPositiveTx(ctx, buyerAddress)
	k.SetSellerExecuteOrderNegativeTx(ctx, buyerAddress)
	k.SetChangeBuyerBidPositiveTx(ctx, buyerAddress)
	k.SetChangeBuyerBidNegativeTx(ctx, buyerAddress)
	k.SetChangeSellerBidPositiveTx(ctx, buyerAddress)
	k.SetChangeSellerBidNegativeTx(ctx, buyerAddress)
	k.SetConfirmBuyerBidPositiveTx(ctx, buyerAddress)
	k.SetConfirmBuyerBidNegativeTx(ctx, buyerAddress)
	k.SetConfirmSellerBidPositiveTx(ctx, buyerAddress)
	k.SetConfirmSellerBidNegativeTx(ctx, buyerAddress)
	k.SetNegotiationPositiveTx(ctx, buyerAddress)
	k.SetNegotiationNegativeTx(ctx, buyerAddress)
	
}
