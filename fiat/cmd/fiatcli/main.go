package main

import (
	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/commitHub/commitBlockchain/client"
	"github.com/commitHub/commitBlockchain/client/keys"
	"github.com/commitHub/commitBlockchain/client/rpc"
	"github.com/commitHub/commitBlockchain/client/tx"
	"github.com/commitHub/commitBlockchain/version"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	bankcmd "github.com/commitHub/commitBlockchain/x/bank/client/cli"
	fiatcmd "github.com/commitHub/commitBlockchain/x/fiatFactory/client/cli"
	ibccmd "github.com/commitHub/commitBlockchain/x/ibc/client/cli"

	"github.com/commitHub/commitBlockchain/client/lcd"
	"github.com/commitHub/commitBlockchain/fiat/app"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "fiatcli",
		Short: "Fiat Chain light-client",
	}
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	// add standard rpc commands
	rpc.AddCommands(rootCmd)

	//Add state commands
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint state querying subcommands",
	}
	tendermintCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(),
	)
	tx.AddCommands(tendermintCmd, cdc)

	//Add IBC commands
	ibcCmd := &cobra.Command{
		Use:   "ibc",
		Short: "Inter-Blockchain Communication subcommands",
	}
	ibcCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCTransferCmd(cdc),
			ibccmd.IBCRelayCmd(cdc),
		)...)

	rootCmd.AddCommand(
		ibcCmd,
		lcd.ServeCommand(cdc),
	)

	//Add auth and bank commands
	rootCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("acc", cdc, authcmd.GetAccountDecoder(cdc)),
			fiatcmd.GetFiatCmd("fiat", cdc, fiatcmd.GetFiatPegDecoder(cdc)),
		)...)
	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
			fiatcmd.IssueFiatCmd(cdc),
			fiatcmd.SendFiatCmd(cdc),
			fiatcmd.ExecuteFiatCmd(cdc),
		)...)

	// add proxy, version and key info
	rootCmd.AddCommand(
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "CF", app.DefaultCLIHome)
	executor.Execute()
}
