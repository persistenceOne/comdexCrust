package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/modules/auth"
	"github.com/persistenceOne/persistenceSDK/modules/auth/client/utils"
	fiatFactoryTypes "github.com/persistenceOne/persistenceSDK/modules/fiatFactory/internal/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

func SendFiatCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send fiat to order with a buyer",
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

			assetPegHashStr := viper.GetString(FlagAssetPegHash)
			assetPegHash, err := types.GetAssetPegHashHex(assetPegHashStr)

			fiatPegHashStr := viper.GetString(FlagFiatPegHash)
			fiatPegHash, err := types.GetFiatPegHashHex(fiatPegHashStr)

			fiatPeg := types.BaseFiatPeg{
				PegHash:           fiatPegHash,
				TransactionAmount: viper.GetInt64(FlagAmount),
			}

			msg := fiatFactoryTypes.BuildSendFiatMsg(cliCtx.GetFromAddress(), ownerAddress, to, assetPegHash,
				types.FiatPegWallet{fiatPeg})

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsAssetPegHash)
	cmd.Flags().AddFlagSet(fsFiatPegHash)
	cmd.Flags().AddFlagSet(fsAmount)
	cmd.Flags().AddFlagSet(fsOwnerAddress)

	return cmd
}