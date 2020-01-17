package keeper_test

import (
	"reflect"
	"testing"

	"github.com/persistenceOne/persistenceSDK/simApp"
	"github.com/persistenceOne/persistenceSDK/types"
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

type prerequisites func()

func TestKeeper_SendAssetsToOrder(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))
	assetPeg := types.BaseAssetPeg{
		PegHash:       []byte("30"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}
	type args struct {
		ctx         cTypes.Context
		fromAddress cTypes.AccAddress
		toAddress   cTypes.AccAddress
		assetPeg    types.AssetPeg
	}
	arg1 := args{ctx, buyer, seller, &assetPeg}
	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"Creating new order and sending asset to it.",
			arg1,
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.OrderKeeper.SendAssetsToOrder(tt.args.ctx, tt.args.fromAddress, tt.args.toAddress, tt.args.assetPeg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.SendAssetsToOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_SendFiatsToOrder(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))

	fiatPeg := types.BaseFiatPeg{
		PegHash:           []byte("30"),
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}
	type args struct {
		ctx           cTypes.Context
		fromAddress   cTypes.AccAddress
		toAddress     cTypes.AccAddress
		pegHash       types.PegHash
		fiatPegWallet types.FiatPegWallet
	}
	arg1 := args{ctx, buyer, seller, []byte("30"), types.FiatPegWallet{fiatPeg}}
	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"Creating new order and sending fiat to it.",
			arg1,
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.OrderKeeper.SendFiatsToOrder(tt.args.ctx, tt.args.fromAddress, tt.args.toAddress, tt.args.pegHash, tt.args.fiatPegWallet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.SendFiatsToOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetOrderDetails(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))

	assetPeg := types.BaseAssetPeg{
		PegHash:       []byte("30"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}

	fiatPeg := types.BaseFiatPeg{
		PegHash:           []byte("30"),
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}

	type args struct {
		ctx           cTypes.Context
		buyerAddress  cTypes.AccAddress
		sellerAddress cTypes.AccAddress
		pegHash       types.PegHash
	}
	arg1 := args{ctx, seller, buyer, assetPeg.GetPegHash()}
	arg2 := args{ctx, buyer, seller, assetPeg.GetPegHash()}
	tests := []struct {
		name  string
		args  args
		pre   prerequisites
		want  cTypes.Error
		want1 types.AssetPegWallet
		want2 types.FiatPegWallet
		want3 string
		want4 string
	}{
		{"Order does not exist.",
			arg1,
			func() {

			},
			cTypes.ErrInvalidAddress("Order not found!"),
			nil,
			nil,
			"",
			"",
		},
		{"Create new order with all details and then get the order details",
			arg2,
			func() {
				order := app.OrderKeeper.NewOrder(buyer, seller, assetPeg.GetPegHash())
				order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
				app.OrderKeeper.SetOrder(ctx, order)
				app.OrderKeeper.SetOrderFiatProofHash(ctx, buyer, seller, assetPeg.GetPegHash(), "fiatProofHash")
				app.OrderKeeper.SetOrderAWBProofHash(ctx, buyer, seller, assetPeg.GetPegHash(), "awbProofHash")
			},
			nil,
			types.AssetPegWallet{assetPeg},
			types.FiatPegWallet{fiatPeg},
			"fiatProofHash",
			"awbProofHash",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			got, got1, got2, got3, got4 := app.OrderKeeper.GetOrderDetails(tt.args.ctx, tt.args.buyerAddress, tt.args.sellerAddress, tt.args.pegHash)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetOrderDetails() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Keeper.GetOrderDetails() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("Keeper.GetOrderDetails() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("Keeper.GetOrderDetails() got3 = %v, want %v", got3, tt.want3)
			}
			if got4 != tt.want4 {
				t.Errorf("Keeper.GetOrderDetails() got4 = %v, want %v", got4, tt.want4)
			}
		})
	}
}

func TestKeeper_IterateOrders(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))

	assetPeg := types.BaseAssetPeg{
		PegHash:       []byte("30"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}

	fiatPeg := types.BaseFiatPeg{
		PegHash:           []byte("30"),
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}

	order := app.OrderKeeper.NewOrder(buyer, seller, assetPeg.GetPegHash())
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	type args struct {
		ctx     cTypes.Context
		process func(types.Order) (stop bool)
	}
	arg := args{ctx, func(types.Order) bool {
		if order.GetAWBProofHash() == "" {
			return true
		}
		return true
	}}
	tests := []struct {
		name string
		args args
	}{
		{"Iterate orders.", arg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.OrderKeeper.IterateOrders(tt.args.ctx, tt.args.process)
		})
	}
}

func TestKeeper_SendAssetFromOrder(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))

	assetPeg := types.BaseAssetPeg{
		PegHash:       []byte("30"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}

	order := app.OrderKeeper.NewOrder(buyer, seller, assetPeg.GetPegHash())
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	type args struct {
		ctx         cTypes.Context
		fromAddress cTypes.AccAddress
		toAddress   cTypes.AccAddress
		assetPeg    types.AssetPeg
	}
	arg1 := args{ctx, buyer, seller, &assetPeg}
	tests := []struct {
		name string
		args args
		want types.AssetPegWallet
	}{
		{"Remove Asset from order",
			arg1,
			types.AssetPegWallet{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.OrderKeeper.SendAssetFromOrder(tt.args.ctx, tt.args.fromAddress, tt.args.toAddress, tt.args.assetPeg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.SendAssetFromOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_SendFiatsFromOrder(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))

	assetPeg := types.BaseAssetPeg{
		PegHash:       []byte("30"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}

	fiatPeg := types.BaseFiatPeg{
		PegHash:           []byte("30"),
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}

	order := app.OrderKeeper.NewOrder(buyer, seller, assetPeg.GetPegHash())
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	type args struct {
		ctx           cTypes.Context
		fromAddress   cTypes.AccAddress
		toAddress     cTypes.AccAddress
		pegHash       types.PegHash
		fiatPegWallet types.FiatPegWallet
	}
	arg1 := args{ctx, buyer, seller, assetPeg.GetPegHash(), types.FiatPegWallet{fiatPeg}}
	tests := []struct {
		name string
		args args
		want types.FiatPegWallet
	}{
		{"Send Fiats from wallet.",
			arg1,
			types.FiatPegWallet{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.OrderKeeper.SendFiatsFromOrder(tt.args.ctx, tt.args.fromAddress, tt.args.toAddress, tt.args.pegHash, tt.args.fiatPegWallet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.SendFiatsFromOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
