package main

import (
	"encoding/json"
	"io"
	
	"github.com/spf13/viper"
	
	"github.com/spf13/cobra"
	
	"github.com/comdex-blockchain/baseapp"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	
	"github.com/comdex-blockchain/fiat/app"
	"github.com/comdex-blockchain/server"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	
	rootCmd := &cobra.Command{
		Use:               "fiatd",
		Short:             "Fiat Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	
	server.AddCommands(ctx, cdc, rootCmd, app.FiatAppInit(),
		server.ConstructAppCreator(newApp, "fiat"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "fiat"))
	
	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "CF", app.DefaultNodeHome)
	
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	return app.NewFiatApp(logger, db, storeTracer, baseapp.SetPruning(viper.GetString("pruning")))
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewFiatApp(logger, db, storeTracer)
	return bapp.ExportAppStateAndValidators()
}
