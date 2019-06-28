package main

import (
	"github.com/spf13/cobra"
	
	"github.com/tendermint/tendermint/libs/cli"
	
	"github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/keys"
	"github.com/comdex-blockchain/client/rpc"
	"github.com/comdex-blockchain/client/tx"
	"github.com/comdex-blockchain/version"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	bankcmd "github.com/comdex-blockchain/x/bank/client/cli"
	fiatcmd "github.com/comdex-blockchain/x/fiatFactory/client/cli"
	ibccmd "github.com/comdex-blockchain/x/ibc/client/cli"
	
	"github.com/comdex-blockchain/client/lcd"
	"github.com/comdex-blockchain/fiat/app"
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
	
	// Add state commands
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint state querying subcommands",
	}
	tendermintCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(),
	)
	tx.AddCommands(tendermintCmd, cdc)
	
	// Add IBC commands
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
	
	// Add auth and bank commands
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
