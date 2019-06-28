package main

import (
	"encoding/json"
	"io"
	
	"github.com/comdex-blockchain/baseapp"
	
	"github.com/comdex-blockchain/main/app"
	"github.com/comdex-blockchain/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	
	rootCmd := &cobra.Command{
		Use:               "maind",
		Short:             "Main Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	
	server.AddCommands(ctx, cdc, rootCmd, app.MainAppInit(),
		server.ConstructAppCreator(newApp, "main"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "main"))
	
	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "CM", app.DefaultNodeHome)
	
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	return app.NewMainApp(logger, db, storeTracer, baseapp.SetPruning(viper.GetString("pruning")))
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewMainApp(logger, db, storeTracer)
	return bapp.ExportAppStateAndValidators()
}
