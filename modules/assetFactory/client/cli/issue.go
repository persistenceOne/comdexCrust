package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/codec"

	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/client/utils"

	assetFactoryTypes "github.com/commitHub/commitBlockchain/modules/assetFactory/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

func IssueAssetCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Initializes asset with the given details and issues to the given address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toStr := viper.GetString(FlagTo)
			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

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
