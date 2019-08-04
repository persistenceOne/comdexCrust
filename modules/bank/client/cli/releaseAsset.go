package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/client/utils"
	"github.com/commitHub/commitBlockchain/modules/bank/client"
	"github.com/commitHub/commitBlockchain/types"
)

func ReleaseAssetCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "releaseAsset",
		Short: "Releases an locked asset peg.",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toStr := viper.GetString(FlagTo)
			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := types.GetAssetPegHashHex(pegHashStr)
			if err != nil {
				return err
			}

			msg := client.BuildReleaseAssetMsg(cliCtx.GetFromAddress(), to, pegHashHex)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	return cmd
}
