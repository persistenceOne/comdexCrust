package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/modules/auth"
	"github.com/persistenceOne/comdexCrust/modules/auth/client/utils"
	fiatFactoryTypes "github.com/persistenceOne/comdexCrust/modules/fiatFactory/internal/types"
	"github.com/persistenceOne/comdexCrust/types"
)

func ExecuteFiatCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute",
		Short: "Send a fiat to an order with a buyer",
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

			assetPegHashStr := viper.GetString(FlagAssetPegHash)
			assetPegHash, err := types.GetAssetPegHashHex(assetPegHashStr)

			fiatPegHashStr := viper.GetString(FlagFiatPegHash)
			fiatPegHash, err := types.GetFiatPegHashHex(fiatPegHashStr)

			fiatPeg := types.BaseFiatPeg{
				PegHash:           fiatPegHash,
				TransactionAmount: viper.GetInt64(FlagAmount),
			}

			msg := fiatFactoryTypes.BuildExecuteFiatMsg(cliCtx.GetFromAddress(), ownerAddress, to, assetPegHash,
				types.FiatPegWallet{fiatPeg})

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsAssetPegHash)
	cmd.Flags().AddFlagSet(fsFiatPegHash)
	cmd.Flags().AddFlagSet(fsOwnerAddress)
	cmd.Flags().AddFlagSet(fsAmount)

	return cmd
}
