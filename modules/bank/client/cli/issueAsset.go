package cli

import (
	"fmt"

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

func IssueAssetCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issueAsset",
		Short: "Initializes asset with the given details and issues to the given address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			var err error
			moderated := viper.GetBool(FlagModerated)
			var to cTypes.AccAddress
			toStr := viper.GetString(FlagTo)
			if moderated && toStr == "" {
				return cTypes.ErrInternal(fmt.Sprintf("must provide toAddress."))
			}

			if toStr == "" {
				to = cliCtx.GetFromAddress()
			} else {
				to, err = cTypes.AccAddressFromBech32(toStr)
				if err != nil {
					return err
				}
				if !moderated {
					if to.String() != cliCtx.GetFromAddress().String() {
						return cTypes.ErrInternal(fmt.Sprintf("Wrong toAddress."))
					}
				}
			}

			assetPeg := &types.BaseAssetPeg{
				AssetQuantity: viper.GetInt64(FlagAssetQuantity),
				AssetType:     viper.GetString(FlagAssetType),
				AssetPrice:    viper.GetInt64(FlagAssetPrice),
				DocumentHash:  viper.GetString(FlagDocumentHash),
				QuantityUnit:  viper.GetString(FlagQuantityUnit),
				Moderated:     moderated,
			}

			msg := client.BuildIssueAssetMsg(cliCtx.GetFromAddress(), to, assetPeg)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsDocumentHash)
	cmd.Flags().AddFlagSet(fsAssetType)
	cmd.Flags().AddFlagSet(fsAssetPrice)
	cmd.Flags().AddFlagSet(fsAssetQuantity)
	cmd.Flags().AddFlagSet(fsQuantityUnit)
	cmd.Flags().AddFlagSet(fsModerated)
	return cmd
}
