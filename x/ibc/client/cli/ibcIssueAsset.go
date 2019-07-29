package cli

import (
	"fmt"
	"os"

	"github.com/commitHub/commitBlockchain/client"
	"github.com/commitHub/commitBlockchain/client/utils"
	context2 "github.com/commitHub/commitBlockchain/x/auth/client/context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/client/context"
	sdk "github.com/commitHub/commitBlockchain/types"

	"github.com/commitHub/commitBlockchain/wire"

	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	"github.com/commitHub/commitBlockchain/x/ibc"
)

const (
	flagDocumentHash  = "documentHash"
	flagAssetType     = "assetType"
	flagAssetPrice    = "assetPrice"
	flagAssetQuantity = "assetQuantity"
	flagQuantityUnit  = "quantityUnit"
	flagModerated     = "moderated"
	// flagChainID       = "chain-id"
)

//IBCIssueAssetCmd : IBC issue asset command
func IBCIssueAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issueAsset",
		Short: "Initializes asset with the given details and issues to the given address",

		RunE: func(cmd *cobra.Command, args []string) error {

			txCtx := context2.NewTxContextFromCLI().
				WithCodec(cdc)

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			sourceChain := viper.GetString(client.FlagChainID)
			destinationChain := viper.GetString(flagChain)
			documentHashStr := viper.GetString(flagDocumentHash)
			assetTypeStr := viper.GetString(flagAssetType)
			assetPriceStr := viper.GetInt64(flagAssetPrice)
			assetQuantityStr := viper.GetInt64(flagAssetQuantity)
			quantityUnitStr := viper.GetString(flagQuantityUnit)
			moderated := viper.GetBool(flagModerated)

			var to sdk.AccAddress
			toStr := viper.GetString(flagTo)
			if moderated && toStr == "" {
				return sdk.ErrInternal(fmt.Sprintf("must provide toAddress."))
			}
			if toStr == "" {
				to = from
			} else {
				to, err = sdk.AccAddressFromBech32(viper.GetString(flagTo))
				if err != nil {
					return err
				}
				if !moderated {
					if to.String() != from.String() {
						return sdk.ErrInternal(fmt.Sprintf("Wrong toAddress."))
					}
				}
			}
			assetPeg := &sdk.BaseAssetPeg{
				AssetQuantity: assetQuantityStr,
				AssetType:     assetTypeStr,
				AssetPrice:    assetPriceStr,
				DocumentHash:  documentHashStr,
				QuantityUnit:  quantityUnitStr,
				Moderated:     moderated,
			}

			msg := ibc.BuildIssueAssetMsg(from, to, assetPeg, sourceChain, destinationChain)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTo, "", "Address to issue asset to")
	cmd.Flags().String(flagDocumentHash, "", "doccument hash")
	cmd.Flags().String(flagAssetType, "", "Asset type")
	cmd.Flags().String(flagAssetPrice, "", "Asset price")
	cmd.Flags().String(flagAssetQuantity, "", "Asset quantity")
	cmd.Flags().String(flagQuantityUnit, "", "Quantity type")
	cmd.Flags().String(flagChain, "", "Destination chain to send coins")
	cmd.Flags().Bool(flagModerated, false, "moderated")
	return cmd
}
