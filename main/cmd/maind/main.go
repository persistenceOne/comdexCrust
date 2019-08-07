package main

import (
	"encoding/json"
	"io"
	
	"github.com/cosmos/cosmos-sdk/baseapp"
	cserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	
	"github.com/commitHub/commitBlockchain/server"
	
	"github.com/commitHub/commitBlockchain/modules/genaccounts"
	genaccscli "github.com/commitHub/commitBlockchain/modules/genaccounts/client/cli"
	genutilcli "github.com/commitHub/commitBlockchain/modules/genutil/client/cli"
	"github.com/commitHub/commitBlockchain/modules/staking"
	
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	
	"github.com/commitHub/commitBlockchain/types"
	
	"github.com/commitHub/commitBlockchain/main/app"
)

const (
	flagInvCheckPeriod = "inv-check-period"
)

var (
	invCheckPeriod uint
)

func main() {
	cobra.EnableCommandSorting = false
	
	cdc := app.MakeCodec()
	
	config := cTypes.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr, types.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr, types.Bech32PrefixConsPub)
	config.Seal()
	
	ctx := cserver.NewDefaultContext()
	
	rootCmd := &cobra.Command{
		Use:               "maind",
		Short:             "Main  Daemon (server)",
		PersistentPreRunE: cserver.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome),
		genutilcli.GenTxCmd(ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{}, genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome),
		genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		genaccscli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
	)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)
	
	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "CM", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewMainApp(logger, db, traceStore, true, invCheckPeriod, baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(cserver.FlagMinGasPrices)),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	
	if height != -1 {
		nsApp := app.NewMainApp(logger, db, traceStore, false, uint(2))
		err := nsApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	
	nsApp := app.NewMainApp(logger, db, traceStore, true, uint(2))
	
	return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
