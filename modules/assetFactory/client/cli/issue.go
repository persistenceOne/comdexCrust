package cli

import (
	"fmt"

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

func IssueAssetCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue [from_key_or_address]",
		Short: "Initializes asset with the given details and issues to the given address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			toStr := viper.GetString(FlagTo)
			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			fmt.Println("\n \n \n \n " + toStr)

			documentHashStr := viper.GetString(FlagDocumentHash)
			assetTypeStr := viper.GetString(FlagAssetType)
			assetPriceStr := viper.GetInt64(FlagAssetPrice)
			assetQuantityStr := viper.GetInt64(FlagAssetQuantity)
			quantityUnitStr := viper.GetString(FlagQuantityUnit)
			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := types.GetAssetPegHashHex(pegHashStr)

			assetPeg := &types.BaseAssetPeg{
				AssetQuantity: assetQuantityStr,
				AssetType:     assetTypeStr,
				AssetPrice:    assetPriceStr,
				DocumentHash:  documentHashStr,
				QuantityUnit:  quantityUnitStr,
				PegHash:       pegHashHex,
			}

			msg := assetFactoryTypes.BuildIssueAssetMsg(cliCtx.GetFromAddress(), to, assetPeg)

			fmt.Println("\n \n \n \n ", msg, "\n \n \n \n")

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsDocumentHash)
	cmd.Flags().AddFlagSet(fsAssetType)
	cmd.Flags().AddFlagSet(fsAssetPrice)
	cmd.Flags().AddFlagSet(fsAssetQuantity)
	cmd.Flags().AddFlagSet(fsQuantityUnit)
	cmd.Flags().AddFlagSet(fsPegHash)

	return cmd
}
