package simulation

import (
	"encoding/json"
	"math/rand"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/x/bank"
	"github.com/commitHub/commitBlockchain/x/gov"
	"github.com/commitHub/commitBlockchain/x/mock"
	"github.com/commitHub/commitBlockchain/x/mock/simulation"
	"github.com/commitHub/commitBlockchain/x/params"
	"github.com/commitHub/commitBlockchain/x/stake"
)

// TestGovWithRandomMessages
func TestGovWithRandomMessages(t *testing.T) {
	mapp := mock.NewApp()

	bank.RegisterWire(mapp.Cdc)
	gov.RegisterWire(mapp.Cdc)
	mapper := mapp.AccountMapper
	coinKeeper := bank.NewKeeper(mapper)
	stakeKey := sdk.NewKVStoreKey("stake")
	stakeKeeper := stake.NewKeeper(mapp.Cdc, stakeKey, coinKeeper, stake.DefaultCodespace)
	paramKey := sdk.NewKVStoreKey("params")
	paramKeeper := params.NewKeeper(mapp.Cdc, paramKey)
	govKey := sdk.NewKVStoreKey("gov")
	govKeeper := gov.NewKeeper(mapp.Cdc, govKey, paramKeeper.Setter(), coinKeeper, stakeKeeper, gov.DefaultCodespace)
	mapp.Router().AddRoute("gov", gov.NewHandler(govKeeper))
	mapp.SetEndBlocker(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		gov.EndBlocker(ctx, govKeeper)
		return abci.ResponseEndBlock{}
	})

	err := mapp.CompleteSetup([]*sdk.KVStoreKey{stakeKey, paramKey, govKey})
	if err != nil {
		panic(err)
	}

	appStateFn := func(r *rand.Rand, keys []crypto.PrivKey, accs []sdk.AccAddress) json.RawMessage {
		mock.RandomSetGenesis(r, mapp, accs, []string{"stake"})
		return json.RawMessage("{}")
	}

	setup := func(r *rand.Rand, privKeys []crypto.PrivKey) {
		ctx := mapp.NewContext(false, abci.Header{})
		stake.InitGenesis(ctx, stakeKeeper, stake.DefaultGenesisState())
		gov.InitGenesis(ctx, govKeeper, gov.DefaultGenesisState())
	}

	simulation.Simulate(
		t, mapp.BaseApp, appStateFn,
		[]simulation.Operation{
			SimulateMsgSubmitProposal(govKeeper, stakeKeeper),
			SimulateMsgDeposit(govKeeper, stakeKeeper),
			SimulateMsgVote(govKeeper, stakeKeeper),
		}, []simulation.RandSetup{
			setup,
		}, []simulation.Invariant{
			AllInvariants(),
		}, 10, 100,
		false,
	)
}
