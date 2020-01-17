package keeper

import (
	"strconv"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/modules/orders"
	reputationTypes "github.com/persistenceOne/persistenceSDK/modules/reputation/internal/types"
	"github.com/persistenceOne/persistenceSDK/types"
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

func (k Keeper) encodeAccountReputation(accountReputation types.AccountReputation) []byte {
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(accountReputation)
	if err != nil {
		panic(err)
	}
	return bz
}

func (k Keeper) decodeAccountReputation(bz []byte) (accountReputation types.AccountReputation) {
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &accountReputation)
	return
}

func (k Keeper) GetAccountReputation(ctx cTypes.Context, addr cTypes.AccAddress) types.AccountReputation {
	store := ctx.KVStore(k.key)
	bz := store.Get(AccountStoreKey(addr))
	if bz == nil {
		accountReputation := types.NewAccountReputation()
		accountReputation.SetAddress(addr)
		k.SetAccountReputation(ctx, accountReputation)
		bz = store.Get(AccountStoreKey(addr))
	}
	accountReputation := k.decodeAccountReputation(bz)
	return accountReputation
}

func (keeper Keeper) GetReputations(ctx cTypes.Context) []types.BaseAccountReputation {
	var reputations []types.BaseAccountReputation

	store := ctx.KVStore(keeper.key)
	iterator := cTypes.KVStorePrefixIterator(store, []byte("address:"))

	for ; iterator.Valid(); iterator.Next() {
		var reputation types.BaseAccountReputation

		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &reputation)
		reputations = append(reputations, reputation)
	}

	return reputations
}

func (k Keeper) SetAccountReputation(ctx cTypes.Context, accountReputation types.AccountReputation) {
	addr := accountReputation.GetAddress()
	store := ctx.KVStore(k.key)
	bz := k.encodeAccountReputation(accountReputation)
	store.Set(AccountStoreKey(addr), bz)
}

func (k Keeper) GetBaseReputationDetails(ctx cTypes.Context, addr cTypes.AccAddress) (cTypes.AccAddress, types.TransactionFeedback, types.TraderFeedbackHistory, int64) {
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

func (k Keeper) SetFeedback(ctx cTypes.Context, addr cTypes.AccAddress, traderFeedback types.TraderFeedback) cTypes.Error {
	accountReputation := k.GetAccountReputation(ctx, addr)
	err := accountReputation.AddTraderFeedback(traderFeedback)
	if err != nil {
		return err
	}
	k.SetAccountReputation(ctx, accountReputation)
	return nil
}

func (k Keeper) SetBuyerRatingToFeedback(ctx cTypes.Context, submitTraderFeedback reputationTypes.SubmitTraderFeedback) cTypes.Error {

	traderFeedback := submitTraderFeedback.TraderFeedback

	err, _, _, fiatProofHash, assetProofHash := k.OrderKeeper.GetOrderDetails(ctx, traderFeedback.BuyerAddress, traderFeedback.SellerAddress, traderFeedback.PegHash)

	if fiatProofHash == "" || assetProofHash == "" {
		return types.ErrFeedbackCannotRegister("you have not completed the transaction to give feedback")
	}

	err = k.SetFeedback(ctx, traderFeedback.SellerAddress, traderFeedback)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			reputationTypes.EventTypeSetBuyerRatingToFeedback,
			cTypes.NewAttribute(reputationTypes.AttributeKeyFrom, traderFeedback.BuyerAddress.String()),
			cTypes.NewAttribute(reputationTypes.AttributeKeyTo, traderFeedback.SellerAddress.String()),
			cTypes.NewAttribute(reputationTypes.AttributeKeyPegHash, traderFeedback.PegHash.String()),
			cTypes.NewAttribute(reputationTypes.AttributeKeyRating, strconv.FormatInt(traderFeedback.Rating, 10)),
		))

	return nil
}

func (k Keeper) SetSellerRatingToFeedback(ctx cTypes.Context, submitTraderFeedback reputationTypes.SubmitTraderFeedback) cTypes.Error {

	traderFeedback := submitTraderFeedback.TraderFeedback

	err, _, _, fiatProofHash, assetProofHash := k.OrderKeeper.GetOrderDetails(ctx, traderFeedback.BuyerAddress, traderFeedback.SellerAddress, traderFeedback.PegHash)

	if fiatProofHash == "" || assetProofHash == "" {
		return types.ErrFeedbackCannotRegister("you have not completed the transaction to give feedback")
	}

	err = k.SetFeedback(ctx, traderFeedback.BuyerAddress, traderFeedback)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			reputationTypes.EventTypeSetSellerRatingToFeedback,
			cTypes.NewAttribute(reputationTypes.AttributeKeyFrom, traderFeedback.SellerAddress.String()),
			cTypes.NewAttribute(reputationTypes.AttributeKeyTo, traderFeedback.BuyerAddress.String()),
			cTypes.NewAttribute(reputationTypes.AttributeKeyPegHash, traderFeedback.PegHash.String()),
			cTypes.NewAttribute(reputationTypes.AttributeKeyRating, strconv.FormatInt(traderFeedback.Rating, 10)),
		))

	return nil
}
