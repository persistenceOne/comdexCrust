package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/comdexCrust/codec"
	assetFactoryTypes "github.com/persistenceOne/comdexCrust/modules/assetFactory/internal/types"
	"github.com/persistenceOne/comdexCrust/modules/auth"
	"github.com/persistenceOne/comdexCrust/modules/auth/client/utils"
	"github.com/persistenceOne/comdexCrust/types"
)

func ExecuteAssetCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute",
		Short: "Send a asset to an order with a buyer",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			ownerAddressStr := viper.GetString(FlagOwnerAddress)

			ownerAddress, err := cTypes.AccAddressFromBech32(ownerAddressStr)
			if err != nil {
				return nil
			}

			toStr := viper.GetString(FlagTo)

			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			assetPegHashStr := viper.GetString(FlagPegHash)
			assetPegHash, err := types.GetAssetPegHashHex(assetPegHashStr)

			msg := assetFactoryTypes.BuildExecuteAssetMsg(cliCtx.GetFromAddress(), ownerAddress, to, assetPegHash)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsOwnerAddress)
	cmd.Flags().AddFlagSet(fsAmount)

	return cmd
}
