package keeper_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/comdexCrust/simApp"
	"github.com/persistenceOne/comdexCrust/types"
)

type prerequisites func()

func TestKeeper_SetAssetPeg(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	assetPeg := types.NewBaseAssetPegWithPegHash([]byte("30"))
	type args struct {
		ctx      cTypes.Context
		assetPeg types.AssetPeg
	}
	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"Set asset peg in assetFactory.",
			args{ctx, &assetPeg},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.AssetKeeper.SetAssetPeg(tt.args.ctx, tt.args.assetPeg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.SetAssetPeg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetAssetPeg(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	assetPeg := types.ProtoBaseAssetPeg()
	assetPeg.SetPegHash([]byte("30"))

	var assetPeg2 types.AssetPeg

	type args struct {
		ctx     cTypes.Context
		peghash types.PegHash
	}
	arg1 := args{ctx, assetPeg.GetPegHash()}
	arg2 := args{ctx, []byte("31")}
	tests := []struct {
		name  string
		args  args
		pre   prerequisites
		want  types.AssetPeg
		want1 cTypes.Error
	}{
		{"Getting asset.",
			arg1,
			func() {
				app.AssetKeeper.SetAssetPeg(ctx, assetPeg)
			},
			assetPeg,
			nil},
		{"Getting asset peg which does not exist",
			arg2,
			func() {},
			assetPeg2,
			cTypes.NewError("asset", 203, fmt.Sprintf("Asset with pegHash %s not found", arg2.peghash))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			got, got1 := app.AssetKeeper.GetAssetPeg(tt.args.ctx, tt.args.peghash)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetAssetPeg() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Keeper.GetAssetPeg() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestKeeper_IterateAssets(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	setAssetPeg := types.ProtoBaseAssetPeg()
	setAssetPeg.SetPegHash([]byte("30"))

	type args struct {
		ctx     cTypes.Context
		handler func(assetPeg types.AssetPeg) (stop bool)
	}
	tests := []struct {
		name string
		args args
		pre  prerequisites
	}{
		{"", args{ctx, func(assetPeg types.AssetPeg) (stop bool) {
			if bytes.Compare(assetPeg.GetPegHash(), setAssetPeg.GetPegHash()) == 0 {
				return true
			}
			return false
		}},
			func() {
				app.AssetKeeper.SetAssetPeg(ctx, setAssetPeg)
			}},
	}
	for _, tt := range tests {
		tt.pre()
		t.Run(tt.name, func(t *testing.T) {
			app.AssetKeeper.IterateAssets(tt.args.ctx, tt.args.handler)
		})
	}
}
