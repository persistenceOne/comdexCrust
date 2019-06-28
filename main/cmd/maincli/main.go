package main

import (
	"os"
	
	"github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/keys"
	"github.com/comdex-blockchain/client/lcd"
	"github.com/comdex-blockchain/client/rpc"
	"github.com/comdex-blockchain/client/tx"
	"github.com/comdex-blockchain/main/app"
	"github.com/comdex-blockchain/version"
	aclcmd "github.com/comdex-blockchain/x/acl/client/cli"
	assetCmd "github.com/comdex-blockchain/x/assetFactory/client/cli"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	bankcmd "github.com/comdex-blockchain/x/bank/client/cli"
	fiatCmd "github.com/comdex-blockchain/x/fiatFactory/client/cli"
	ibccmd "github.com/comdex-blockchain/x/ibc/client/cli"
	negotiationcmd "github.com/comdex-blockchain/x/negotiation/client/cli"
	ordercmd "github.com/comdex-blockchain/x/order/client/cli"
	reputationcmd "github.com/comdex-blockchain/x/reputation/client/cli"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "maincli",
		Short: "Main light-client",
	}
)

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false
	
	// get the codec
	cdc := app.MakeCodec()
	
	// TODO: Setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc.
	
	// add standard rpc, and tx commands
	rpc.AddCommands(rootCmd)
	tx.AddCommands(rootCmd, cdc)
	rootCmd.AddCommand(client.LineBreak)
	
	// Add IBC commands
	ibcCmd := &cobra.Command{
		Use:   "ibc",
		Short: "Inter-Blockchain Communication subcommands",
	}
	ibcCmd.AddCommand(
		client.PostCommands(
			ibccmd.IBCTransferCmd(cdc),
			ibccmd.IBCRelayCmd(cdc),
			ibccmd.IBCIssueAssetCmd(cdc),
			ibccmd.IBCRedeemAssetCmd(cdc),
			ibccmd.IBCIssueFiatCmd(cdc),
			ibccmd.IBCRedeemFiatCmd(cdc),
			ibccmd.IBCSendAssetCmd(cdc),
			ibccmd.IBCSendFiatCmd(cdc),
			ibccmd.IBCBuyerExecuteOrder(cdc),
			ibccmd.IBCSellerExecuteOrder(cdc),
		)...)
	rootCmd.AddCommand(ibcCmd)
	
	// Add auth and bank commands
	rootCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("acc", cdc, authcmd.GetAccountDecoder(cdc)),
			negotiationcmd.GetNegotiationCmd("negotiation", cdc, negotiationcmd.GetNegotiationDecoder(cdc)),
			ordercmd.GetOrderCmd("order", cdc, ordercmd.GetOrderDecoder(cdc)),
			aclcmd.GetACLAccountCmd("acl", cdc, aclcmd.GetACLAccountDecoder(cdc)),
			aclcmd.GetZoneCmd("acl", cdc),
			aclcmd.GetOrganizationCmd("acl", cdc),
			reputationcmd.GetReputationCmd("reputation", cdc, reputationcmd.GetReputationDecoder(cdc)),
			bankcmd.GetAssetCmd("asset", cdc, assetCmd.GetAssetPegDecoder(cdc)),
			bankcmd.GetFiatCmd("fiat", cdc, fiatCmd.GetFiatPegDecoder(cdc)),
		)...)
	
	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
			bankcmd.IssueAssetCmd(cdc),
			bankcmd.RedeemAssetCmd(cdc),
			bankcmd.IssueFiatCmd(cdc),
			bankcmd.SendAssetCmd(cdc),
			bankcmd.SendFiatCmd(cdc),
			bankcmd.BuyerExecuteOrderCmd(cdc),
			bankcmd.SellerExecuteOrderCmd(cdc),
			bankcmd.ReleaseAssetCmd(cdc),
			bankcmd.DefineZoneCmd(cdc),
			bankcmd.DefineOrganizationCmd(cdc),
			bankcmd.DefineACLCmd(cdc),
			negotiationcmd.ChangeBuyerBidCmd(cdc),
			negotiationcmd.ChangeSellerBidCmd(cdc),
			negotiationcmd.ConfirmBuyerBidCmd(cdc),
			negotiationcmd.ConfirmSellerBidCmd(cdc),
			reputationcmd.SubmitBuyerFeedbackCmd(cdc),
			reputationcmd.SubmitSellerFeedbackCmd(cdc),
		)...)
	
	// add proxy, version and key info
	rootCmd.AddCommand(
		lcd.ServeCommand(cdc),
		keys.Commands(),
		version.VersionCmd,
	)
	
	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "CM", os.ExpandEnv("$HOME/.maincli"))
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
