package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/codec"

	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/client/utils"

	"github.com/commitHub/commitBlockchain/types"

	assetFactoryTypes "github.com/commitHub/commitBlockchain/modules/assetFactory/internal/types"
)

func SendAssetCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send asset to order with a buyer",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			ownerAddr := viper.GetString(FlagOwnerAddress)

			ownerAddress, err := cTypes.AccAddressFromBech32(ownerAddr)
			if err != nil {
				return nil
			}

			toAddr := viper.GetString(FlagTo)
			to, err := cTypes.AccAddressFromBech32(toAddr)
			if err != nil {
				return nil
			}

			assetPegHashStr := viper.GetString(FlagPegHash)
			assetPegHash, err := types.GetAssetPegHashHex(assetPegHashStr)

			msg := assetFactoryTypes.BuildSendAssetMsg(cliCtx.GetFromAddress(), ownerAddress, to, assetPegHash)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsOwnerAddress)

	return cmd
}
