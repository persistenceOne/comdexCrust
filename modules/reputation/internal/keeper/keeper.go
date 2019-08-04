package keeper

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/codec"

	"github.com/commitHub/commitBlockchain/modules/orders"
	reputationTypes "github.com/commitHub/commitBlockchain/modules/reputation/internal/types"
)

type Keeper struct {
	key         cTypes.StoreKey
	cdc         *codec.Codec
	OrderKeeper orders.Keeper
}

func NewKeeper(cdc *codec.Codec, key cTypes.StoreKey, orderKeeper orders.Keeper) Keeper {
	return Keeper{
		key:         key,
		cdc:         cdc,
		OrderKeeper: orderKeeper,
	}
}

func AccountStoreKey(addr cTypes.AccAddress) []byte {
	return append([]byte("address:"), addr.Bytes()...)
}

func (k Keeper) encodeAccountReputation(accountReputation reputationTypes.AccountReputation) []byte {
	bz, err := k.cdc.MarshalBinaryBare(accountReputation)
	if err != nil {
		panic(err)
	}
	return bz
}

func (k Keeper) decodeAccountReputation(bz []byte) (accountReputation reputationTypes.AccountReputation) {
	err := k.cdc.UnmarshalBinaryBare(bz, &accountReputation)
	if err != nil {
		panic(err)
	}
	return
}

func (k Keeper) GetAccountReputation(ctx cTypes.Context, addr cTypes.AccAddress) reputationTypes.AccountReputation {
	store := ctx.KVStore(k.key)
	bz := store.Get(AccountStoreKey(addr))
	if bz == nil {
		accountReputation := reputationTypes.NewAccountReputation()
		accountReputation.SetAddress(addr)
		k.SetAccountReputation(ctx, accountReputation)
		bz = store.Get(AccountStoreKey(addr))
	}
	accountReputation := k.decodeAccountReputation(bz)
	return accountReputation
}

func (k Keeper) SetAccountReputation(ctx cTypes.Context, accountReputation reputationTypes.AccountReputation) {
	addr := accountReputation.GetAddress()
	store := ctx.KVStore(k.key)
	bz := k.encodeAccountReputation(accountReputation)
	store.Set(AccountStoreKey(addr), bz)
}

func (k Keeper) GetBaseReputationDetails(ctx cTypes.Context, addr cTypes.AccAddress) (cTypes.AccAddress, reputationTypes.TransactionFeedback, reputationTypes.TraderFeedbackHistory, int64) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	return accountReputation.GetAddress(), accountReputation.GetTransactionFeedback(), accountReputation.GetTraderFeedbackHistory(), accountReputation.GetRating()
}

func (k Keeper) SetSendAssetsPositiveTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SendAssetsPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetSendAssetsNegativeTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SendAssetsNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetSendFiatsPositiveTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SendFiatsPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetSendFiatsNegativeTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SendFiatsNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetIBCIssueAssetsPositiveTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.IBCIssueAssetsPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetIBCIssueAssetsNegativeTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.IBCIssueAssetsNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetIBCIssueFiatsPositiveTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.IBCIssueFiatsPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetIBCIssueFiatsNegativeTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.IBCIssueFiatsNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetBuyerExecuteOrderPositiveTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.BuyerExecuteOrderPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetBuyerExecuteOrderNegativeTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.BuyerExecuteOrderNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetSellerExecuteOrderPositiveTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SellerExecuteOrderPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetSellerExecuteOrderNegativeTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.SellerExecuteOrderNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetChangeBuyerBidPositiveTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ChangeBuyerBidPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

// SetChangeBuyerBidNegativeTx : increases count by 1
func (k Keeper) SetChangeBuyerBidNegativeTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ChangeBuyerBidNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetChangeSellerBidPositiveTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ChangeSellerBidPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetChangeSellerBidNegativeTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ChangeSellerBidNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetConfirmBuyerBidPositiveTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ConfirmBuyerBidPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetConfirmBuyerBidNegativeTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ConfirmBuyerBidNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetConfirmSellerBidPositiveTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ConfirmSellerBidPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetConfirmSellerBidNegativeTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.ConfirmSellerBidNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetNegotiationPositiveTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.NegotiationPositiveTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetNegotiationNegativeTx(ctx cTypes.Context, addr cTypes.AccAddress) {
	accountReputation := k.GetAccountReputation(ctx, addr)
	transactionFeedback := accountReputation.GetTransactionFeedback()
	transactionFeedback.NegotiationNegativeTx++
	_ = accountReputation.SetTransactionFeedback(transactionFeedback)
	k.SetAccountReputation(ctx, accountReputation)
}

func (k Keeper) SetFeedback(ctx cTypes.Context, addr cTypes.AccAddress, traderFeedback reputationTypes.TraderFeedback) cTypes.Error {
	accountReputation := k.GetAccountReputation(ctx, addr)
	err := accountReputation.AddTraderFeedback(traderFeedback)
	if err != nil {
		return err
	}
	k.SetAccountReputation(ctx, accountReputation)
	return nil
}
