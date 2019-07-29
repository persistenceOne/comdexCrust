package main

import (
	"os"

	"github.com/commitHub/commitBlockchain/client"
	"github.com/commitHub/commitBlockchain/client/keys"
	"github.com/commitHub/commitBlockchain/client/lcd"
	"github.com/commitHub/commitBlockchain/client/rpc"
	"github.com/commitHub/commitBlockchain/client/tx"
	"github.com/commitHub/commitBlockchain/main/app"
	"github.com/commitHub/commitBlockchain/version"
	aclcmd "github.com/commitHub/commitBlockchain/x/acl/client/cli"
	assetCmd "github.com/commitHub/commitBlockchain/x/assetFactory/client/cli"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	bankcmd "github.com/commitHub/commitBlockchain/x/bank/client/cli"
	fiatCmd "github.com/commitHub/commitBlockchain/x/fiatFactory/client/cli"
	gov "github.com/commitHub/commitBlockchain/x/gov/client/cli"
	ibccmd "github.com/commitHub/commitBlockchain/x/ibc/client/cli"
	negotiationcmd "github.com/commitHub/commitBlockchain/x/negotiation/client/cli"
	ordercmd "github.com/commitHub/commitBlockchain/x/order/client/cli"
	reputationcmd "github.com/commitHub/commitBlockchain/x/reputation/client/cli"
	stake "github.com/commitHub/commitBlockchain/x/stake/client/cli"
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

	//Add IBC commands
	ibcCmd := &cobra.Command{
		Use:   "ibc",
		Short: "Inter-Blockchain Communication subcommands",
	}
	stakecmd := &cobra.Command{
		Use:   "stake",
		Short: "stake module subcommands",
	}
	stakecmd.AddCommand(
		client.PostCommands(
			stake.GetCmdCreateValidator(cdc),
			stake.GetCmdDelegate(cdc),
			stake.GetCmdEditValidator(cdc),
			stake.GetCmdRedelegate("stake", cdc),
			stake.GetCmdUnbond("stake", cdc),
		)...)
	stakecmd.AddCommand(
		client.GetCommands(
			stake.GetCmdQueryValidator("stake", cdc),
			stake.GetCmdQueryValidators("stake", cdc),
			stake.GetCmdQueryDelegation("stake", cdc),
			stake.GetCmdQueryDelegations("stake", cdc),
			stake.GetCmdQueryUnbondingDelegation("stake", cdc),
			stake.GetCmdQueryRedelegation("stake", cdc),
			stake.GetCmdQueryRedelegations("stake", cdc),
			stake.GetCmdQueryPool("stake", cdc),
			stake.GetCmdQueryParams("stake", cdc),
		)...)
	rootCmd.AddCommand(stakecmd)

	govcmd := &cobra.Command{
		Use:   "gov",
		Short: "gov module subcommands",
	}

	govcmd.AddCommand(
		client.PostCommands(
			gov.GetCmdSubmitProposal(cdc),
			gov.GetCmdDeposit(cdc),
			gov.GetCmdVote(cdc),
		)...)
	govcmd.AddCommand(
		client.GetCommands(
			gov.GetCmdQueryProposal("gov", cdc),
			gov.GetCmdQueryProposals("gov", cdc),
			gov.GetCmdQueryVote("gov", cdc),
			gov.GetCmdQueryVotes("gov", cdc),
			gov.GetCmdQueryDeposit("gov", cdc),
			gov.GetCmdQueryDeposits("gov", cdc),
			gov.GetCmdQueryTally("gov", cdc),
		)...)
	rootCmd.AddCommand(govcmd)
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

	//Add auth and bank commands
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
