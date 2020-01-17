package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/modules/auth"
	"github.com/persistenceOne/persistenceSDK/modules/auth/client/utils"
	"github.com/persistenceOne/persistenceSDK/modules/bank/client"
	"github.com/persistenceOne/persistenceSDK/types"
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
