package genaccounts

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/commitHub/commitBlockchain/codec"
)

// initialize accounts and deliver genesis transactions
func InitGenesis(ctx sdk.Context, _ *codec.Codec, accountKeeper AccountKeeper, genesisState GenesisState) {
	genesisState.Sanitize()
	
	// load the accounts
	for _, gacc := range genesisState {
		acc := gacc.ToAccount()
		acc = accountKeeper.NewAccount(ctx, acc) // set account number
		accountKeeper.SetAccount(ctx, acc)
	}
}
