package keeper

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	abciTypes "github.com/tendermint/tendermint/abci/types"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	"github.com/commitHub/commitBlockchain/simApp"
)

func TestBaseKeeper_DelegateCoins(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	now := tmtime.Now()
	ctx = ctx.WithBlockHeader(abciTypes.Header{Time: now})
	endTime := now.Add(24 * time.Hour)
	ak := app.AccountKeeper

	origCoins := cTypes.NewCoins(cTypes.NewInt64Coin("stake", 100))
	delCoins := cTypes.NewCoins(cTypes.NewInt64Coin("stake", 50))

	addr1 := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))
	addrModule := sdk.AccAddress([]byte("moduleAcc"))

	bacc := auth.NewBaseAccountWithAddress(addr1)
	bacc.SetCoins(origCoins)
	macc := ak.NewAccountWithAddress(ctx, addrModule) // we don't need to define an actual module account bc we just need the address for testing
	vacc := vesting.NewContinuousVestingAccount(&bacc, ctx.BlockHeader().Time.Unix(), endTime.Unix())
	acc := ak.NewAccountWithAddress(ctx, addr2)
	ak.SetAccount(ctx, vacc)
	ak.SetAccount(ctx, acc)
	ak.SetAccount(ctx, macc)
	app.BankKeeper.SetCoins(ctx, addr2, origCoins)

	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))

	// require the ability for a non-vesting account to delegate
	err := app.BankKeeper.DelegateCoins(ctx, addr2, addrModule, delCoins)
	acc = ak.GetAccount(ctx, addr2)
	macc = ak.GetAccount(ctx, addrModule)
	require.NoError(t, err)
	require.Equal(t, origCoins.Sub(delCoins), acc.GetCoins())
	require.Equal(t, delCoins, macc.GetCoins())

	// require the ability for a vesting account to delegate
	err = app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins)
	vacc = ak.GetAccount(ctx, addr1).(*vesting.ContinuousVestingAccount)
	require.NoError(t, err)
	require.Equal(t, delCoins, vacc.GetCoins())

}

