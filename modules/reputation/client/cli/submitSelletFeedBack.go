package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/modules/auth"
	"github.com/persistenceOne/persistenceSDK/modules/auth/client/utils"
	reputationTypes "github.com/persistenceOne/persistenceSDK/modules/reputation/internal/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

func SubmitSellerFeedbackCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sellerFeedback",
		Short: "seller gives feedback to the buyer",
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

			msg := reputationTypes.BuildSellerFeedbackMsg(toAddress, cliCtx.GetFromAddress(), pegHashHex, rating)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsPeghash)
	cmd.Flags().AddFlagSet(fsRating)
	cmd.Flags().AddFlagSet(fsTo)
	return cmd
}