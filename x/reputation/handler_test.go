package reputation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/commitHub/commitBlockchain/store"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/order"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setup() (*wire.Codec, sdk.Context, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()

	orderKey := sdk.NewKVStoreKey("orderKey")
	reputationKey := sdk.NewKVStoreKey("reputation")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(orderKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(reputationKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	cdc := wire.NewCodec()
	order.RegisterOrder(cdc)
	sdk.RegisterWire(cdc)
	RegisterWire(cdc)
	RegisterReputation(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	return cdc, ctx, orderKey, reputationKey
}

var (
	cdc, ctx, orderKey, reputationKey = setup()
	orderMapper                       = order.NewMapper(cdc, orderKey, sdk.ProtoBaseOrder)
	orderKeeper                       = order.NewKeeper(orderMapper)

	reputationMapper = NewMapper(cdc, reputationKey, sdk.ProtoBaseAccountReputation)
	k                = NewKeeper(reputationMapper)
)

func TestNewFeedbackHandler(t *testing.T) {

	toTest := NewFeedbackHandler(k, orderKeeper)

	var buyer = []sdk.AccAddress{
		sdk.AccAddress("buyer1"),
		sdk.AccAddress(""),
		sdk.AccAddress(nil),
		sdk.AccAddress("buyer2"),
	}
	var seller = []sdk.AccAddress{
		sdk.AccAddress("seller1"),
		sdk.AccAddress(""),
		sdk.AccAddress(nil),
		sdk.AccAddress("seller2"),
	}

	var peghash = []sdk.PegHash{
		sdk.PegHash("31"),
		sdk.PegHash("32"),
		sdk.PegHash("33"),
	}

	order := orderMapper.NewOrder(buyer[0], seller[0], peghash[0])
	order2 := orderMapper.NewOrder(buyer[3], seller[3], peghash[0])

	orderMapper.SetOrder(ctx, order)
	orderMapper.SetOrder(ctx, order2)

	orderKeeper.SetOrderAWBProofHash(ctx, buyer[0], seller[0], peghash[0], "awb1")
	orderKeeper.SetOrderFiatProofHash(ctx, buyer[0], seller[0], peghash[0], "fph1")

	var msgFeedback = []sdk.Msg{
		NewMsgBuyerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[0], seller[0], peghash[0], 55))}),
		NewMsgBuyerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[0], seller[0], peghash[0], 56))}),
		NewMsgBuyerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[1], seller[1], peghash[1], 55))}),
		NewMsgBuyerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[0], seller[2], peghash[2], 55))}),
		NewMsgSellerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[0], seller[0], peghash[0], 55))}),
		NewMsgSellerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[0], seller[0], peghash[0], 56))}),
		NewMsgSellerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[1], seller[1], peghash[1], 55))}),
		NewMsgSellerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[0], seller[2], peghash[2], 55))}),
		NewMsgBuyerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[3], seller[3], peghash[0], 56))}),
		NewMsgSellerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[3], seller[3], peghash[0], 56))}),
		NewMsgSellerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(nil, seller[3], peghash[0], 56))}),
		NewMsgSellerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[3], nil, peghash[0], 56))}),
		NewMsgSellerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[3], seller[3], nil, 56))}),
		NewMsgSellerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[3], seller[3], peghash[0], 105))}),
		NewMsgBuyerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(nil, seller[3], peghash[0], 56))}),
		NewMsgBuyerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[3], nil, peghash[0], 56))}),
		NewMsgBuyerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[3], seller[3], nil, 56))}),
		NewMsgBuyerFeedbacks([]SubmitTraderFeedback{NewSubmitTraderFeedback(sdk.NewTraderFeedback(buyer[3], seller[3], peghash[0], 105))}),
	}
	buyerMsg := msgFeedback[0]
	sellerMsg := msgFeedback[4]

	require.Equal(t, buyerMsg.Type(), "reputation")
	require.Equal(t, sellerMsg.Type(), "reputation")
	require.Nil(t, buyerMsg.ValidateBasic())
	require.Nil(t, sellerMsg.ValidateBasic())
	require.NotNil(t, msgFeedback[10].ValidateBasic())
	require.NotNil(t, msgFeedback[11].ValidateBasic())
	require.NotNil(t, msgFeedback[12].ValidateBasic())
	require.NotNil(t, msgFeedback[13].ValidateBasic())
	require.NotNil(t, msgFeedback[14].ValidateBasic())
	require.NotNil(t, msgFeedback[15].ValidateBasic())
	require.NotNil(t, msgFeedback[16].ValidateBasic())
	require.NotNil(t, msgFeedback[17].ValidateBasic())
	require.NotNil(t, buyerMsg.GetSignBytes())
	require.NotNil(t, sellerMsg.GetSignBytes())
	require.Equal(t, buyerMsg.GetSigners()[0], buyer[0])
	require.Equal(t, sellerMsg.GetSigners()[0], seller[0])
	require.Equal(t, BuildBuyerFeedbackMsg(buyer[0], seller[0], peghash[0], 55).GetSigners(), buyerMsg.GetSigners())
	require.Equal(t, BuildSellerFeedbackMsg(buyer[0], seller[0], peghash[0], 55).GetSigners(), sellerMsg.GetSigners())
	require.Nil(t, toTest(ctx, BuildBuyerFeedbackMsg(buyer[0], seller[0], peghash[0], 55)).Tags)

	tests1 := []struct {
		result bool
		msg    sdk.Msg
	}{
		{true, msgFeedback[0]},
		{true, msgFeedback[1]},
		{false, msgFeedback[2]},
		{false, msgFeedback[3]},
		{true, msgFeedback[4]},
		{true, msgFeedback[5]},
		{false, msgFeedback[6]},
		{false, msgFeedback[7]},
		{false, msgFeedback[8]},
		{false, msgFeedback[9]},
	}
	for _, tt := range tests1 {
		if tt.result == true {
			tags := toTest(ctx, tt.msg).Tags
			fmt.Println(string(tags[1].Value))
			require.NotNil(t, tags)
		} else {
			fmt.Println(toTest(ctx, tt.msg).Tags)
			require.Nil(t, toTest(ctx, tt.msg).Tags)
		}
	}
}
