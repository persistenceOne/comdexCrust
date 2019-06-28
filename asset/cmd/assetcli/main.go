package main

import (
	"os"
	
	"github.com/comdex-blockchain/asset/app"
	"github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/keys"
	"github.com/comdex-blockchain/client/lcd"
	"github.com/comdex-blockchain/client/rpc"
	"github.com/comdex-blockchain/client/tx"
	"github.com/comdex-blockchain/version"
	assetcmd "github.com/comdex-blockchain/x/assetFactory/client/cli"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	bankcmd "github.com/comdex-blockchain/x/bank/client/cli"
	ibccmd "github.com/comdex-blockchain/x/ibc/client/cli"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "assetcli",
		Short: "Asset light-client",
	}
)

func main() {
	cobra.EnableCommandSorting = false
	
	cdc := app.MakeCodec()
	
	rpc.AddCommands(rootCmd)
	rootCmd.AddCommand(client.LineBreak)
	tx.AddCommands(rootCmd, cdc)
	rootCmd.AddCommand(client.LineBreak)
	
	rootCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("acc", cdc, authcmd.GetAccountDecoder(cdc)),
			assetcmd.GetAssetCmd("asset", cdc, assetcmd.GetAssetPegDecoder(cdc)),
		)...)
	
	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
			assetcmd.IssueAssetCmd(cdc),
			assetcmd.RedeemAssetCmd(cdc),
			assetcmd.SendAssetCmd(cdc),
			assetcmd.ExecuteAssetCmd(cdc),
		)...)
	
	ibcCmd := &cobra.Command{
		Use:   "ibc",
		Short: "Inter-Blockchain Communication subcommands",
	}
	ibcCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCTransferCmd(cdc),
			ibccmd.IBCRelayCmd(cdc),
		)...)
	rootCmd.AddCommand(ibcCmd)
	
	rootCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)
	
	executor := cli.PrepareMainCmd(rootCmd, "CA", os.ExpandEnv("$HOME/.assetcli"))
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
