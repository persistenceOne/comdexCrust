package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/codec"
	assetFactoryTypes "github.com/commitHub/commitBlockchain/modules/assetFactory/internal/types"
	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/client/utils"
	"github.com/commitHub/commitBlockchain/types"
)

func RedeemAssetCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem",
		Short: "Redeem asset with the given details and redeem to the given address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toAddressStr := viper.GetString(FlagTo)
			toAddress, err := cTypes.AccAddressFromBech32(toAddressStr)
			if err != nil {
				return nil
			}

			ownerAddressStr := viper.GetString(FlagOwnerAddress)
			ownerAddress, err := cTypes.AccAddressFromBech32(ownerAddressStr)
			if err != nil {
				return nil
			}

			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := types.GetAssetPegHashHex(pegHashStr)

			msg := assetFactoryTypes.BuildRedeemAssetMsg(cliCtx.GetFromAddress(), ownerAddress, toAddress, pegHashHex)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsOwnerAddress)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsTo)

	return cmd
}
