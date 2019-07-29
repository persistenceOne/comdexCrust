package cli

import (
	"os"

	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/assetFactory"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	context2 "github.com/commitHub/commitBlockchain/x/auth/client/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//IssueAssetCmd : create a init asset tx and sign it with the give key
func IssueAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Initializes asset with the given details and issues to the given address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txCtx := context2.NewTxContextFromCLI().WithCodec(cdc)

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			toStr := viper.GetString(FlagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			documentHashStr := viper.GetString(FlagDocumentHash)
			assetTypeStr := viper.GetString(FlagAssetType)
			assetPriceStr := viper.GetInt64(FlagAssetPrice)
			assetQuantityStr := viper.GetInt64(FlagAssetQuantity)
			quantityUnitStr := viper.GetString(FlagQuantityUnit)
			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := sdk.GetAssetPegHashHex(pegHashStr)

			assetPeg := &sdk.BaseAssetPeg{
				AssetQuantity: assetQuantityStr,
				AssetType:     assetTypeStr,
				AssetPrice:    assetPriceStr,
				DocumentHash:  documentHashStr,
				QuantityUnit:  quantityUnitStr,
				PegHash:       pegHashHex,
			}

			msg := assetFactory.BuildIssueAssetMsg(from, to, assetPeg)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
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
