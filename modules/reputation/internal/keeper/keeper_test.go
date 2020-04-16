package keeper_test

import (
	"reflect"
	"testing"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	reputationTypes "github.com/persistenceOne/comdexCrust/modules/reputation/internal/types"
	"github.com/persistenceOne/comdexCrust/simApp"
	"github.com/persistenceOne/comdexCrust/types"
)

type prerequisites func()

func TestKeeper_SetSendAssetsPositiveTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))

	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set send assets positive transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetSendAssetsPositiveTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetSendAssetsNegativeTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set send assets negative transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetSendAssetsNegativeTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetSendFiatsPositiveTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))

	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set send fiats positive transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetSendFiatsPositiveTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetSendFiatsNegativeTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set send fiats negative transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetSendFiatsNegativeTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetIBCIssueAssetsPositiveTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set IBC issue asset positive  transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetIBCIssueAssetsPositiveTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetIBCIssueAssetsNegativeTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set IBC issue asset negative  transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetIBCIssueAssetsNegativeTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetIBCIssueFiatsPositiveTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set IBC issue fiat positive  transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetIBCIssueFiatsPositiveTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetIBCIssueFiatsNegativeTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set IBC issue fiat negative  transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetIBCIssueFiatsNegativeTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetBuyerExecuteOrderPositiveTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set buyer execute order positive transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetBuyerExecuteOrderPositiveTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetBuyerExecuteOrderNegativeTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set buyer execute order negative transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetBuyerExecuteOrderNegativeTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetSellerExecuteOrderPositiveTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set seller execute order positive transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetSellerExecuteOrderPositiveTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetSellerExecuteOrderNegativeTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set seller execute order negative transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetSellerExecuteOrderNegativeTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetChangeBuyerBidPositiveTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set change buyer bid positive transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetChangeBuyerBidPositiveTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetChangeBuyerBidNegativeTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set change buyer bid negative transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetChangeBuyerBidNegativeTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetChangeSellerBidPositiveTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set change seller bid positive transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetChangeSellerBidPositiveTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetChangeSellerBidNegativeTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set change seller bid negative transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetChangeSellerBidNegativeTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetConfirmBuyerBidPositiveTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set confirm buyer bid positive transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetConfirmBuyerBidPositiveTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetConfirmBuyerBidNegativeTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set confirm buyer bid negitive transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetConfirmBuyerBidNegativeTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetConfirmSellerBidPositiveTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set confirm seller bid positive transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetConfirmSellerBidPositiveTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetConfirmSellerBidNegativeTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set confirm seller bid negitive transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetConfirmSellerBidNegativeTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetNegotiationPositiveTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set negotiation positive transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetNegotiationPositiveTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetNegotiationNegativeTx(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set negotiation negative transaction.",
			args{ctx, addr1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.ReputationKeeper.SetNegotiationNegativeTx(tt.args.ctx, tt.args.addr)
		})
	}
}

func TestKeeper_SetFeedback(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))
	pegHash := types.PegHash([]byte("30"))

	traderFeedback := types.NewTraderFeedback(buyer, seller, pegHash, 50)
	type args struct {
		ctx            cTypes.Context
		addr           cTypes.AccAddress
		traderFeedback types.TraderFeedback
	}
	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"Buyer is setting feedback for the trade.",
			args{ctx, buyer, traderFeedback},
			nil},
		{"Buyer is setting feedback for the trade again.",
			args{ctx, buyer, traderFeedback},
			types.ErrFeedbackCannotRegister("types")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.ReputationKeeper.SetFeedback(tt.args.ctx, tt.args.addr, tt.args.traderFeedback); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.SetFeedback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_SetBuyerRatingToFeedback(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))
	pegHash := types.PegHash([]byte("30"))

	traderFeedback := types.NewTraderFeedback(buyer, seller, pegHash, 50)
	reputationMsg := reputationTypes.NewSubmitTraderFeedback(traderFeedback)

	type args struct {
		ctx                  cTypes.Context
		submitTraderFeedback reputationTypes.SubmitTraderFeedback
	}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"Buyer setting rating for seller before trade is complete.",
			args{ctx, reputationMsg},
			func() {},
			types.ErrFeedbackCannotRegister("you have not completed the transaction to give feedback")},
		{"Buyer is setting feedback for the trade.",
			args{ctx, reputationMsg},
			func() {
				order := app.OrderKeeper.NewOrder(buyer, seller, pegHash)
				order.SetAWBProofHash("awbProofHash")
				order.SetFiatProofHash("fiatProofHash")
				app.OrderKeeper.SetOrder(ctx, order)
			},
			nil},
		{"Buyer is setting feedback for the trade again.",
			args{ctx, reputationMsg},
			func() {},
			types.ErrFeedbackCannotRegister("types")},
	}
	for _, tt := range tests {
		tt.pre()
		t.Run(tt.name, func(t *testing.T) {
			if got := app.ReputationKeeper.SetBuyerRatingToFeedback(tt.args.ctx, tt.args.submitTraderFeedback); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.SetBuyerRatingToFeedback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_SetSellerRatingToFeedback(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))
	pegHash := types.PegHash([]byte("30"))

	traderFeedback := types.NewTraderFeedback(buyer, seller, pegHash, 50)
	reputationMsg := reputationTypes.NewSubmitTraderFeedback(traderFeedback)

	type args struct {
		ctx                  cTypes.Context
		submitTraderFeedback reputationTypes.SubmitTraderFeedback
	}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"Seller setting rating for buyer before trade is complete.",
			args{ctx, reputationMsg},
			func() {},
			types.ErrFeedbackCannotRegister("you have not completed the transaction to give feedback")},
		{"Seller is setting feedback for the trade.",
			args{ctx, reputationMsg},
			func() {
				order := app.OrderKeeper.NewOrder(buyer, seller, pegHash)
				order.SetAWBProofHash("awbProofHash")
				order.SetFiatProofHash("fiatProofHash")
				app.OrderKeeper.SetOrder(ctx, order)
			},
			nil},
		{"Seller is setting feedback for the trade again.",
			args{ctx, reputationMsg},
			func() {},
			types.ErrFeedbackCannotRegister("types")},
	}
	for _, tt := range tests {
		tt.pre()
		t.Run(tt.name, func(t *testing.T) {
			if got := app.ReputationKeeper.SetSellerRatingToFeedback(tt.args.ctx, tt.args.submitTraderFeedback); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.SetSellerRatingToFeedback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetReputations(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	var accountReputation types.BaseAccountReputation
	app.ReputationKeeper.SetAccountReputation(ctx, &accountReputation)

	var reputations []types.BaseAccountReputation

	type args struct {
		ctx cTypes.Context
	}
	tests := []struct {
		name string
		args args
		want []types.BaseAccountReputation
	}{
		{"Get accountReputation.",
			args{ctx},
			append(reputations, accountReputation)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.ReputationKeeper.GetReputations(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetReputations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetBaseReputationDetails(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))
	pegHash := types.PegHash([]byte("30"))
	traderFeedback := types.NewTraderFeedback(buyer, seller, pegHash, 50)

	accountReputation := types.NewBaseAccountReputation()
	accountReputation.SetAddress(buyer)
	accountReputation.SetTraderFeedbackHistory([]types.TraderFeedback{traderFeedback})
	app.ReputationKeeper.SetAccountReputation(ctx, &accountReputation)

	type args struct {
		ctx  cTypes.Context
		addr cTypes.AccAddress
	}
	tests := []struct {
		name  string
		args  args
		want  cTypes.AccAddress
		want1 types.TransactionFeedback
		want2 types.TraderFeedbackHistory
		want3 int64
	}{
		{"Get accountReputation details.",
			args{ctx, buyer},
			accountReputation.GetAddress(),
			accountReputation.GetTransactionFeedback(),
			accountReputation.GetTraderFeedbackHistory(),
			accountReputation.GetRating()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := app.ReputationKeeper.GetBaseReputationDetails(tt.args.ctx, tt.args.addr)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetBaseReputationDetails() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Keeper.GetBaseReputationDetails() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("Keeper.GetBaseReputationDetails() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("Keeper.GetBaseReputationDetails() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}
