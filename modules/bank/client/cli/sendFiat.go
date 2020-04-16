package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/modules/auth"
	"github.com/persistenceOne/comdexCrust/modules/auth/client/utils"
	"github.com/persistenceOne/comdexCrust/modules/bank/client"
	"github.com/persistenceOne/comdexCrust/types"
)

func SendFiatCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendFiat",
		Short: "Sends an fiat peg to an order transaction with a given address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toStr := viper.GetString(FlagTo)

			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			amount := viper.GetInt64(FlagAmount)

			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := types.GetAssetPegHashHex(pegHashStr)

			msg := client.BuildSendFiatMsg(cliCtx.GetFromAddress(), to, pegHashHex, amount)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsAmount)
	cmd.Flags().AddFlagSet(fsPegHash)
	return cmd
}
