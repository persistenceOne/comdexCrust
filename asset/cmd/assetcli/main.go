package main

import (
	"os"

	"github.com/commitHub/commitBlockchain/asset/app"
	"github.com/commitHub/commitBlockchain/client"
	"github.com/commitHub/commitBlockchain/client/keys"
	"github.com/commitHub/commitBlockchain/client/lcd"
	"github.com/commitHub/commitBlockchain/client/rpc"
	"github.com/commitHub/commitBlockchain/client/tx"
	"github.com/commitHub/commitBlockchain/version"
	assetcmd "github.com/commitHub/commitBlockchain/x/assetFactory/client/cli"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	bankcmd "github.com/commitHub/commitBlockchain/x/bank/client/cli"
	ibccmd "github.com/commitHub/commitBlockchain/x/ibc/client/cli"
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
