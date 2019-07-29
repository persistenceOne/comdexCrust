package cli

import (
	"github.com/commitHub/commitBlockchain/client/context"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	context2 "github.com/commitHub/commitBlockchain/x/auth/client/context"
	"github.com/commitHub/commitBlockchain/x/negotiation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"os"

	"github.com/commitHub/commitBlockchain/client/utils"
)

//ChangeSellerBidCmd : Change or create a negotiation bid
func ChangeSellerBidCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changeSellerBid",
		Short: "Change or create a seller negotiation bid",
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
			negotiationID := sdk.NegotiationID(append(append(to.Bytes(), from.Bytes()...), pegHashHex.Bytes()...))

			proposedNegotiation := &sdk.BaseNegotiation{
				NegotiationID: negotiationID,
				BuyerAddress:  to,
				SellerAddress: from,
				PegHash:       pegHashHex,
				Bid:           bidInt64,
				Time:          timeInt64,
			}
			msg := negotiation.BuildMsgChangeSellerBid(proposedNegotiation)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})

		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsBid)
	cmd.Flags().AddFlagSet(fsTime)
	return cmd
}
