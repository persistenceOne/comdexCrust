package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/codec"

	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/client/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	negotiationTypes "github.com/commitHub/commitBlockchain/modules/negotiation/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

func ChangeSellerBidCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-seller-bid",
		Short: "Change or create a seller negotiation bid",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toStr := viper.GetString(FlagTo)

			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			bid := viper.GetInt64(FlagBid)
			time := viper.GetInt64(FlagTime)
			hashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := types.GetAssetPegHashHex(hashStr)
			negotiationID := negotiationTypes.NegotiationID(append(append(to.Bytes(), cliCtx.GetFromAddress().Bytes()...), pegHashHex...))

			proposedNegotiation := &negotiationTypes.BaseNegotiation{
				NegotiationID: negotiationID,
				BuyerAddress:  to,
				SellerAddress: cliCtx.GetFromAddress(),
				PegHash:       pegHashHex,
				Bid:           bid,
				Time:          time,
			}

			msg := negotiationTypes.BuildMsgChangeSellerBid(proposedNegotiation)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsBid)
	cmd.Flags().AddFlagSet(fsTime)
	return cmd
}
