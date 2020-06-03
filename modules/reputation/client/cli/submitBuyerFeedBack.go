package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/codec"

	"github.com/commitHub/commitBlockchain/types"

	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/client/utils"
	reputationTypes "github.com/commitHub/commitBlockchain/modules/reputation/internal/types"
)

func SubmitBuyerFeedbackCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buyerFeedback",
		Short: "buyer gives feedback to the seller",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			to := viper.GetString(FlagTo)
			toAddress, err := cTypes.AccAddressFromBech32(to)
			if err != nil {
				return err
			}

			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := types.GetAssetPegHashHex(pegHashStr)
			if err != nil {
				return err
			}

			rating := viper.GetInt64(FlagRating)

			msg := reputationTypes.BuildBuyerFeedbackMsg(cliCtx.GetFromAddress(), toAddress, pegHashHex, rating)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsPeghash)
	cmd.Flags().AddFlagSet(fsRating)
	cmd.Flags().AddFlagSet(fsTo)
	return cmd
}
