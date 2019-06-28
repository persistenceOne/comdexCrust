package reputation

import (
	"strconv"
	
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/x/order"
)

// Keeper : for txns keeping
type Keeper struct {
	fm Mapper
}

// NewKeeper : keeper for mapper feedback
func NewKeeper(fm Mapper) Keeper {
	return Keeper{fm}
}

// GetBaseReputationDetails : gives feedback interface, all getters are present
func (fk Keeper) GetBaseReputationDetails(ctx sdk.Context, addr sdk.AccAddress) (sdk.AccAddress, sdk.TransactionFeedback, sdk.TraderFeedbackHistory, int64) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	return accountReputation.GetAddress(), accountReputation.GetTransactionFeedback(), accountReputation.GetTraderFeedbackHistory(), accountReputation.GetRating()
}

// SetSendAssetsPositiveTx : increases count by 1
func (fk Keeper) SetSendAssetsPositiveTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SendAssetsPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetSendAssetsNegativeTx : increases count by 1
func (fk Keeper) SetSendAssetsNegativeTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SendAssetsNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetSendFiatsPositiveTx : increases count by 1
func (fk Keeper) SetSendFiatsPositiveTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SendFiatsPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetSendFiatsNegativeTx : increases count by 1
func (fk Keeper) SetSendFiatsNegativeTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SendFiatsNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetIBCIssueAssetsPositiveTx : increases count by 1
func (fk Keeper) SetIBCIssueAssetsPositiveTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.IBCIssueAssetsPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetIBCIssueAssetsNegativeTx : increases count by 1
func (fk Keeper) SetIBCIssueAssetsNegativeTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.IBCIssueAssetsNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetIBCIssueFiatsPositiveTx : increases count by 1
func (fk Keeper) SetIBCIssueFiatsPositiveTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.IBCIssueFiatsPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetIBCIssueFiatsNegativeTx : increases count by 1
func (fk Keeper) SetIBCIssueFiatsNegativeTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.IBCIssueFiatsNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetBuyerExecuteOrderPositiveTx : increases count by 1
func (fk Keeper) SetBuyerExecuteOrderPositiveTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.BuyerExecuteOrderPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetBuyerExecuteOrderNegativeTx : increases count by 1
func (fk Keeper) SetBuyerExecuteOrderNegativeTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.BuyerExecuteOrderNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetSellerExecuteOrderPositiveTx : increases count by 1
func (fk Keeper) SetSellerExecuteOrderPositiveTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SellerExecuteOrderPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetSellerExecuteOrderNegativeTx : increases count by 1
func (fk Keeper) SetSellerExecuteOrderNegativeTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SellerExecuteOrderNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetChangeBuyerBidPositiveTx : increases count by 1
func (fk Keeper) SetChangeBuyerBidPositiveTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ChangeBuyerBidPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetChangeBuyerBidNegativeTx : increases count by 1
func (fk Keeper) SetChangeBuyerBidNegativeTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ChangeBuyerBidNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetChangeSellerBidPositiveTx : increases count by 1
func (fk Keeper) SetChangeSellerBidPositiveTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ChangeSellerBidPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetChangeSellerBidNegativeTx : increases count by 1
func (fk Keeper) SetChangeSellerBidNegativeTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ChangeSellerBidNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetConfirmBuyerBidPositiveTx : increases count by 1
func (fk Keeper) SetConfirmBuyerBidPositiveTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ConfirmBuyerBidPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetConfirmBuyerBidNegativeTx : increases count by 1
func (fk Keeper) SetConfirmBuyerBidNegativeTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ConfirmBuyerBidNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetConfirmSellerBidPositiveTx : increases count by 1
func (fk Keeper) SetConfirmSellerBidPositiveTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ConfirmSellerBidPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetConfirmSellerBidNegativeTx : increases count by 1
func (fk Keeper) SetConfirmSellerBidNegativeTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ConfirmSellerBidNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetNegotiationPositiveTx : increases count by 1
func (fk Keeper) SetNegotiationPositiveTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.NegotiationPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetNegotiationNegativeTx : increases count by 1
func (fk Keeper) SetNegotiationNegativeTx(ctx sdk.Context, addr sdk.AccAddress) {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.NegotiationNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	fk.fm.SetAccountReputation(ctx, accountReputation)
}

// SetFeedback : adds new rating
func (fk Keeper) SetFeedback(ctx sdk.Context, addr sdk.AccAddress, traderFeedback sdk.TraderFeedback) sdk.Error {
	accountReputation := fk.fm.GetAccountReputation(ctx, addr)
	err := accountReputation.AddTraderFeedback(traderFeedback)
	if err != nil {
		return err
	}
	fk.fm.SetAccountReputation(ctx, accountReputation)
	return nil
}

// SetBuyerRatingToFeedback : handler calls this method
func (fk Keeper) SetBuyerRatingToFeedback(ctx sdk.Context, orderkeeper order.Keeper, msgFeedback MsgBuyerFeedbacks) (sdk.Tags, sdk.Error) {
	
	allTags := sdk.EmptyTags()
	tags := sdk.EmptyTags()
	for _, submitTraderFeedback := range msgFeedback.SubmitTraderFeedbacks {
		traderFeedback := submitTraderFeedback.TraderFeedback
		err, _, _, fiatProofHash, awbProofHash := orderkeeper.GetOrderDetails(ctx, traderFeedback.BuyerAddress, traderFeedback.SellerAddress, traderFeedback.PegHash)
		if err != nil {
			return nil, err
		}
		
		if fiatProofHash == "" || awbProofHash == "" {
			return nil, sdk.ErrFeedbackCannotRegister("you have not completed the transaction to give feedback")
		}
		
		err = fk.SetFeedback(ctx, traderFeedback.SellerAddress, traderFeedback)
		if err != nil {
			return nil, err
		}
		tags = sdk.NewTags("from", []byte(traderFeedback.BuyerAddress.String()))
		tags = tags.AppendTag("to", []byte(traderFeedback.SellerAddress.String()))
		tags = tags.AppendTag("peghash", []byte(traderFeedback.PegHash))
		tags = tags.AppendTag("rating", []byte(strconv.FormatInt(traderFeedback.Rating, 10)))
	}
	allTags = allTags.AppendTags(tags)
	return allTags, nil
}

// SetSellerRatingToFeedback : handler calls this method
func (fk Keeper) SetSellerRatingToFeedback(ctx sdk.Context, orderkeeper order.Keeper, msgFeedback MsgSellerFeedbacks) (sdk.Tags, sdk.Error) {
	
	allTags := sdk.EmptyTags()
	tags := sdk.EmptyTags()
	for _, submitTraderFeedback := range msgFeedback.SubmitTraderFeedbacks {
		traderFeedback := submitTraderFeedback.TraderFeedback
		err, _, _, fiatProofHash, awbProofHash := orderkeeper.GetOrderDetails(ctx, traderFeedback.BuyerAddress, traderFeedback.SellerAddress, traderFeedback.PegHash)
		if err != nil {
			return nil, err
		}
		
		if fiatProofHash == "" || awbProofHash == "" {
			return nil, sdk.ErrFeedbackCannotRegister("you have not completed the transaction to give feedback")
		}
		
		err = fk.SetFeedback(ctx, traderFeedback.BuyerAddress, traderFeedback)
		if err != nil {
			return nil, err
		}
		tags = sdk.NewTags("from", []byte(traderFeedback.SellerAddress.String()))
		tags = tags.AppendTag("to", []byte(traderFeedback.BuyerAddress.String()))
		tags = tags.AppendTag("peghash", []byte(traderFeedback.PegHash))
		tags = tags.AppendTag("rating", []byte(strconv.FormatInt(traderFeedback.Rating, 10)))
	}
	allTags = allTags.AppendTags(tags)
	return allTags, nil
}
