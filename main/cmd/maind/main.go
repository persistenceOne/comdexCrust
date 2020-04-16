package main

import (
	"encoding/json"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	cServer "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/persistenceOne/comdexCrust/main/app"
	"github.com/persistenceOne/comdexCrust/modules/genaccounts"
	genAccsCli "github.com/persistenceOne/comdexCrust/modules/genaccounts/client/cli"
	genUtilCli "github.com/persistenceOne/comdexCrust/modules/genutil/client/cli"
	"github.com/persistenceOne/comdexCrust/modules/staking"
	"github.com/persistenceOne/comdexCrust/server"
	"github.com/persistenceOne/comdexCrust/types"
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

	ctx := cServer.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "maind",
		Short:             "Main  Daemon (server)",
		PersistentPreRunE: cServer.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		genUtilCli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		genUtilCli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome),
		genUtilCli.GenTxCmd(ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{}, genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome),
		genUtilCli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		genAccsCli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
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

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abciTypes.Application {
	return app.NewMainApp(logger, db, traceStore, true, invCheckPeriod, baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(cServer.FlagMinGasPrices)),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmTypes.GenesisValidator, error) {

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
