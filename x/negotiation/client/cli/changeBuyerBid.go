package cli

import (
	"os"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/negotiation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ChangeBuyerBidCmd : Change or create a negotiation bid
func ChangeBuyerBidCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changeBuyerBid",
		Short: "Change or create a buyer negotiation bid",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			txCtx := context2.NewTxContextFromCLI().
				WithCodec(cdc)
			
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			
			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}
			
			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			
			toStr := viper.GetString(FlagTo)
			
			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}
			
			bidInt64 := viper.GetInt64(FlagBid)
			timeInt64 := viper.GetInt64(FlagTime)
			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := sdk.GetAssetPegHashHex(pegHashStr)
			negotiationID := sdk.NegotiationID(append(append(from.Bytes(), to.Bytes()...), pegHashHex.Bytes()...))
			
			proposedNegotiation := &sdk.BaseNegotiation{
				NegotiationID: negotiationID,
				BuyerAddress:  from,
				SellerAddress: to,
				PegHash:       pegHashHex,
				Bid:           bidInt64,
				Time:          timeInt64,
			}
			msg := negotiation.BuildMsgChangeBuyerBid(proposedNegotiation)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsBid)
	cmd.Flags().AddFlagSet(fsTime)
	return cmd
}
