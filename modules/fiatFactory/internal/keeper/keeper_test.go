package keeper_test

import (
	"fmt"
	"reflect"
	"testing"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/comdexCrust/simApp"
	"github.com/persistenceOne/comdexCrust/types"
)

type prerequisites func()

func TestKeeper_SetFiatPeg(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	fiatPeg := types.ProtoBaseFiatPeg()
	fiatPeg.SetPegHash([]byte("30"))
	type args struct {
		ctx     cTypes.Context
		fiatPeg types.FiatPeg
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set fiat peg",
			args{ctx, fiatPeg}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.FiatKeeper.SetFiatPeg(tt.args.ctx, tt.args.fiatPeg)
		})
	}
}

func TestKeeper_GetFiatPeg(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	fiatPeg := types.ProtoBaseFiatPeg()
	fiatPeg.SetPegHash([]byte("30"))

	var fiatPeg2 types.FiatPeg

	type args struct {
		ctx     cTypes.Context
		pegHash types.PegHash
	}
	arg1 := args{ctx, fiatPeg.GetPegHash()}
	arg2 := args{ctx, []byte("31")}
	tests := []struct {
		name  string
		args  args
		pre   prerequisites
		want  types.FiatPeg
		want1 cTypes.Error
	}{
		{"Getting fiat.",
			arg1,
			func() {
				app.FiatKeeper.SetFiatPeg(ctx, fiatPeg)
			},
			fiatPeg,
			nil},
		{"Getting fiat peg which does not exist",
			arg2,
			func() {},
			fiatPeg2,
			cTypes.NewError("fiat", 602, fmt.Sprintf("Fiat with pegHash %s not found", arg2.pegHash))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			got, got1 := app.FiatKeeper.GetFiatPeg(tt.args.ctx, tt.args.pegHash)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetFiatPeg() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Keeper.GetFiatPeg() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
